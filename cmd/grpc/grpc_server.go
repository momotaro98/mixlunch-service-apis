package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/momotaro98/stew"

	"github.com/momotaro98/mixlunch-service-api/conventions"
	"github.com/momotaro98/mixlunch-service-api/logger"
	"github.com/momotaro98/mixlunch-service-api/partyservice"
	"github.com/momotaro98/mixlunch-service-api/pb"
	"github.com/momotaro98/mixlunch-service-api/tagservice"
	usService "github.com/momotaro98/mixlunch-service-api/userscheduleservice"
	"github.com/momotaro98/mixlunch-service-api/userservice"
)

func makeRFC3399FormatDateTimeStr(year int, month int, day int, hour int, min int) string {
	// RFC3339 format is accepted.
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:00Z", year, month, day, hour, min)
}

func generateBeginEndOfTheDay(targetDateStr string) (beginDTStr string, endDTStr string) {
	ymdSlice := strings.Split(targetDateStr, "-")
	year, _ := strconv.Atoi(ymdSlice[0])
	month, _ := strconv.Atoi(ymdSlice[1])
	day, _ := strconv.Atoi(ymdSlice[2])
	beginDTStr = makeRFC3399FormatDateTimeStr(year, month, day, 0, 0)
	endDTStr = makeRFC3399FormatDateTimeStr(year, month, day, 23, 59)
	return beginDTStr, endDTStr
}

type gRPCMixLunchServer struct {
	logger      logger.Logger
	usServer    usService.UserScheduleServer
	partyServer partyservice.PartyServer
	userServer  userservice.UserServer
}

func provideGRPCMixLunchServer(logger logger.Logger,
	usServer usService.UserScheduleServer,
	partyServer partyservice.PartyServer,
	userServer userservice.UserServer,
) *gRPCMixLunchServer {
	return &gRPCMixLunchServer{
		logger:      logger,
		usServer:    usServer,
		partyServer: partyServer,
		userServer:  userServer,
	}
}

func (s *gRPCMixLunchServer) GetUsersForMatching(targetDate *pb.TargetDate, stream pb.MixLunch_GetUsersForMatchingServer) error {
	s.logger.Log(logger.Info, "", fmt.Sprintf("Start GetUsersForMatching process with TargetDate, %v", *targetDate))
	// Retrieve users from DB
	beginDateTimeStr, endDateTimeStr := generateBeginEndOfTheDay(targetDate.Date)
	eachUserSchedulesOfTheDate, err := s.usServer.GetEachUserSchedules(beginDateTimeStr, endDateTimeStr)
	if err != nil {
		s.logger.Log(logger.Error, "", err.Error())
		return stew.Wrap(err)
	}

	// Assign the data into pb.UserModelForMatching and send it to client with gRPC stream
	for _, aUserSchedule := range eachUserSchedulesOfTheDate {
		// Request to user service
		user, err := s.userServer.GetUserByUserId(aUserSchedule.UserId)
		if err != nil {
			s.logger.Log(logger.Error, "", err.Error())
			return stew.Wrap(err)
		}

		// [Business Logic] Assemble Blacklist User
		blacklistOfTheUser, err := s.assembleBlacklist(user)
		if err != nil {
			return stew.Wrap(err)
		}

		// Populate to gRPC proto buffer model
		var userModelForMatching = pb.UserModelForMatching{
			UserId:       aUserSchedule.UserId,
			FreeFrom:     aUserSchedule.UserSchedules[0].FromDateTime.Format(conventions.TimeFormat),
			FreeTo:       aUserSchedule.UserSchedules[0].ToDateTime.Format(conventions.TimeFormat),
			UserName:     user.Name,
			Email:        user.Email,
			HaveTags:     convertSliceTagsToIDs(user.InterestTags, user.SkillTags),
			WantTags:     convertSliceTagsToIDs(aUserSchedule.UserSchedules[0].Tags),
			Latitude:     aUserSchedule.UserSchedules[0].Location.Latitude,
			Longitude:    aUserSchedule.UserSchedules[0].Location.Longitude,
			LocationType: int32(aUserSchedule.UserSchedules[0].Location.LocationTypeID),
			Blacklist:    blacklistOfTheUser,
			Languages:    convertSliceLangToString(user.Languages),
		}

		// Send to Party service via gRPC
		if err := stream.Send(&userModelForMatching); err != nil {
			s.logger.Log(logger.Error, "", err.Error())
			return stew.Wrap(err)
		}
	}

	s.logger.Log(logger.Info, "", "Process succeeded")
	return nil
}

func (s *gRPCMixLunchServer) assembleBlacklist(user *userservice.User) (blacklistUsers []string, err error) {
	// Pure black list
	blacklistUsers = append(blacklistUsers, user.BlockingUsers...)

	// Add black list by Avoiding Business Logic
	const (
		ignoreTimes = 3
		daysAgo     = -14
	)

	// [Business Logic] 最後のランチから ignoreTimes 回分のランチメイトは無視する。
	pLastN, err := s.partyServer.GetLastNPartiesOfAUser(user.UserId, ignoreTimes)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	// [Business Logic] 現在から直前の daysAgo 日間でランチしたランチメイトは無視する。
	var (
		begin = time.Now().AddDate(0, 0, daysAgo)
	)
	pTimeRange, err := s.partyServer.GetPartyByUserIdAndTimeRange(
		user.UserId, begin.Format(time.RFC3339), "")
	if err != nil {
		return nil, stew.Wrap(err)
	}

	targetParties := append(pLastN.Parties, pTimeRange.Parties...)

	for _, party := range targetParties {
		for _, member := range party.Members {
			if member.UserId != user.UserId {
				blacklistUsers = append(blacklistUsers, member.UserId)
			}
		}
	}

	return blacklistUsers, nil
}

func convertSliceTagsToIDs(sliceCTags ...[]*tagservice.CategoryTags) []int32 {
	var sTags []*tagservice.SmallTag
	for _, cTags := range sliceCTags {
		for _, cTag := range cTags {
			sTags = append(sTags, cTag.Tags...)
		}
	}
	ret := make([]int32, 0, len(sTags))
	for _, t := range sTags {
		ret = append(ret, int32(t.TagId))
	}
	return ret
}

func convertSliceLangToString(langs []userservice.Language) []string {
	ret := make([]string, 0, len(langs))
	for _, l := range langs {
		ret = append(ret, string(l))
	}
	return ret
}

func (s *gRPCMixLunchServer) CreateParties(stream pb.MixLunch_CreatePartiesServer) error {
	s.logger.Log(logger.Info, "", "Start CreateParties process")
	partyChan := make(chan *partyservice.PartyForCommand)

	// Receive parties via gRPC
	recErr := make(chan error)
	go receivePartiesFromMatchingModule(stream, partyChan, recErr)

	// Insert parties into DB
	go func() {
		parties := make([]*partyservice.PartyForCommand, 0)
		for party := range partyChan {
			// Generate Chat room by using passed Chat Room ID
			go func(p *partyservice.PartyForCommand) {
				if err := s.partyServer.GenerateChatRoom(p.ChatRoomId); err != nil {
					s.logger.Log(logger.Error, "", fmt.Sprintf("Failed to generate chatroom ID: %v, err: %v", party, err))
				}
			}(party)
			// Add the party to DB
			parties = append(parties, party)
		}

		if err := s.partyServer.UpsertParties(parties); err != nil {
			s.logger.Log(logger.Error, "", fmt.Sprintf("Upserting the parties failed. err: %+v", err))
			panic(err)
		} else {
			s.logger.Log(logger.Info, "", "Upserting the parties succeeded")
		}
	}()

	if err, open := <-recErr; open {
		s.logger.Log(logger.Error, "", err.Error())
		return err
	}

	return nil
}

func receivePartiesFromMatchingModule(stream pb.MixLunch_CreatePartiesServer, partyChan chan<- *partyservice.PartyForCommand, errChan chan<- error) {
	defer func() {
		_ = stream.SendAndClose(&pb.Empty{})
		close(partyChan)
		close(errChan)
	}()
	for {
		party, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			errChan <- err
			break
		}
		// Translate protocol buffer model to our domain model
		// startFrom
		startFrom, err := time.Parse(time.RFC3339, party.StartFrom)
		if err != nil {
			errChan <- err
			break
		}
		// endTo
		endTo, err := time.Parse(time.RFC3339, party.EndTo)
		if err != nil {
			errChan <- err
			break
		}
		// chatRoomId
		chatRoomId := party.RoomId
		// []*User
		var membersModel []*userservice.UserPublic
		for _, member := range party.Members {
			var m = &userservice.UserPublic{UserId: member.UserId}
			membersModel = append(membersModel, m)
		}
		// Send a party domain model
		partyChan <- partyservice.NewPartyForCommand(startFrom, endTo, chatRoomId, membersModel)
	}
}

func (s *gRPCMixLunchServer) GetParties(targetDate *pb.TargetDate, stream pb.MixLunch_GetPartiesServer) error {
	s.logger.Log(logger.Info, "", fmt.Sprintf("Start GetParties process with TargetDate, %v", *targetDate))
	// Retrieve parties from DB
	beginDateTimeStr, endDateTimeStr := generateBeginEndOfTheDay(targetDate.Date)
	partiesOfTheDate, err := s.partyServer.GetParties(beginDateTimeStr, endDateTimeStr)
	if err != nil {
		s.logger.Log(logger.Error, "", err.Error())
		return err
	}
	// Assign the data into pb.Party and send it to client with gRPC stream
	for _, party := range partiesOfTheDate.Parties {
		var partyToSend pb.Party
		// StartFrom
		fDT := party.StartFrom
		partyToSend.StartFrom = makeRFC3399FormatDateTimeStr(fDT.Year(), int(fDT.Month()), fDT.Day(), fDT.Hour(), fDT.Minute())
		// EndTo
		tDT := party.EndTo
		partyToSend.EndTo = makeRFC3399FormatDateTimeStr(tDT.Year(), int(tDT.Month()), tDT.Day(), tDT.Hour(), tDT.Minute())
		// Members
		var members []*pb.UserModelForMatching
		for _, member := range party.Members {
			var memberToSend pb.UserModelForMatching
			memberToSend.UserId = member.UserId
			memberToSend.FreeFrom = makeRFC3399FormatDateTimeStr(fDT.Year(), int(fDT.Month()), fDT.Day(), fDT.Hour(), fDT.Minute())
			memberToSend.FreeTo = makeRFC3399FormatDateTimeStr(tDT.Year(), int(tDT.Month()), tDT.Day(), tDT.Hour(), tDT.Minute())
			memberToSend.UserName = member.Name
			memberToSend.Email = member.Email
			members = append(members, &memberToSend)
		}
		partyToSend.Members = members
		// ChatRoomId
		partyToSend.ChatRoomId = party.ChatRoomId

		// Send to client
		if err := stream.Send(&partyToSend); err != nil {
			s.logger.Log(logger.Error, "", err.Error())
			return err
		}
	}
	s.logger.Log(logger.Info, "", "Process succeeded")
	return nil
}
