package partyservice

import (
	"database/sql"
	"errors"
	"time"

	"github.com/momotaro98/stew"

	"github.com/momotaro98/mixlunch-service-api/domainerror"
	"github.com/momotaro98/mixlunch-service-api/tagservice"
	"github.com/momotaro98/mixlunch-service-api/userservice"
	"github.com/momotaro98/mixlunch-service-api/utils"
)

type Parties struct {
	Parties []*Party `json:"parties"`
}

type Party struct {
	PartyID    int                        `json:"party_id"`
	StartFrom  time.Time                  `json:"start_from"`
	EndTo      time.Time                  `json:"end_to"`
	ChatRoomId string                     `json:"chat_room_id"`
	Members    []*userservice.UserPublic  `json:"members"`
	Tags       []*tagservice.CategoryTags `json:"tags"`
}

func NewParty(partyID int, startFrom, endTo time.Time, chatRoomId string,
	members []*userservice.UserPublic, tags []*tagservice.CategoryTags) *Party {
	return &Party{
		PartyID:    partyID,
		StartFrom:  startFrom,
		EndTo:      endTo,
		ChatRoomId: chatRoomId,
		Members:    members,
		Tags:       tags,
	}
}

type PartyForCommand struct {
	StartFrom  time.Time                 `json:"start_from"`
	EndTo      time.Time                 `json:"end_to"`
	ChatRoomId string                    `json:"chat_room_id"`
	Members    []*userservice.UserPublic `json:"members"`
}

func NewPartyForCommand(startFrom, endTo time.Time, chatRoomId string,
	members []*userservice.UserPublic) *PartyForCommand {
	return &PartyForCommand{
		StartFrom:  startFrom,
		EndTo:      endTo,
		ChatRoomId: chatRoomId,
		Members:    members,
	}
}

type IsLatestReviewDone struct {
	IsReviewDone bool `json:"is_review_done"`
}

// [Note] Trying to separate "Command" and "Query" model in User type constructor

type PartyServer interface {
	GetParties(beginDateTimeStr, endDateTimeStr string) (*Parties, error)
	GetPartyByUserIdAndTimeRange(userId, beginDateTimeStr, endDateTimeStr string) (*Parties, error)
	GetIsLatestPartyReviewDone(userId string) (*IsLatestReviewDone, error)
	GetLastNPartiesOfAUser(userId string, n int) (*Parties, error)
	PostPartyReviewMember(reviewMember *PartyReviewMember) error
	UpsertParties(partyModel []*PartyForCommand) error
	GenerateChatRoom(chatRoomId string) error
}

func ProvidePartyServer(
	queryRepository IPartyQueryRepository,
	updateRepository IPartyCommandRepository,
	userServer userservice.UserServer,
	tagServer tagservice.TagServer,
	chatRoomRepository IChatRoomRepository) PartyServer {
	return &realPartyServer{
		partyQueryRepository:   queryRepository,
		partyCommandRepository: updateRepository,
		userServer:             userServer,
		tagServer:              tagServer,
		chatRoomRepository:     chatRoomRepository,
	}
}

type realPartyServer struct {
	partyQueryRepository   IPartyQueryRepository
	partyCommandRepository IPartyCommandRepository
	userServer             userservice.UserServer
	tagServer              tagservice.TagServer
	chatRoomRepository     IChatRoomRepository
}

func (s *realPartyServer) GetParties(beginDateTimeStr, endDateTimeStr string) (*Parties, error) {
	// Parse begin and end DateTime string to RFC3339 spec
	beginDateTime, endDateTime, err := parseBeginEndDateTime(beginDateTimeStr, endDateTimeStr)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	// Query parties
	partyDtos, err := s.partyQueryRepository.QueryPartiesWhereTimeRange(&PartyQueryDto{
		beginDateTime: beginDateTime,
		endDateTime:   endDateTime,
	})
	if err != nil {
		return nil, stew.Wrap(err)
	}
	if len(partyDtos) < 1 {
		return &Parties{}, nil
	}
	// Assign Parties
	parties, err := s.populateIntoParties(partyDtos) // Parties of GetParties doesn't require user ID
	if err != nil {
		return nil, stew.Wrap(err)
	}
	return parties, nil
}

func (s *realPartyServer) GetPartyByUserIdAndTimeRange(userId, beginDateTimeStr, endDateTimeStr string) (*Parties, error) {
	// Parse begin and end DateTime string to RFC3339 spec
	beginDateTime, endDateTime, err := parseBeginEndDateTime(beginDateTimeStr, endDateTimeStr)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	// Query parties
	partyDtos, err := s.partyQueryRepository.QueryPartiesWhereUserIdAndTimeRange(
		userId,
		&PartyQueryDto{
			beginDateTime: beginDateTime,
			endDateTime:   endDateTime,
		},
	)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	if len(partyDtos) < 1 {
		return &Parties{
			Parties: make([]*Party, 0),
		}, nil
	}
	// Assign Parties
	parties, err := s.populateIntoParties(partyDtos)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	return parties, nil
}

func (s *realPartyServer) populateIntoParties(partyDtos []*PartyDto) (*Parties, error) {
	// Issue: This is 1+N query issue. Fix by involving GetParties and GetPartyByUserIdAndTimeRange when we have time.
	var parties = Parties{
		Parties: make([]*Party, 0),
	}
	for _, pDto := range partyDtos {
		var members []*userservice.UserPublic
		// Get Party members in each party
		memberDtos, err := s.partyQueryRepository.QueryPartyMembersWherePartyId(pDto.id)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		for _, memberDto := range memberDtos {
			// Assign user model and add members list
			// UserId
			userId := memberDto.userId
			// UserName and Email
			userPublic, err := s.userServer.GetUserPublicByUserId(userId)
			if err != nil {
				return nil, stew.Wrap(err)
			}
			members = append(members, userPublic)
		}

		// Get Tags of the party
		tagIdsDto, err := s.partyQueryRepository.QueryPartyTagsWherePartyId(pDto.id)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		tags, err := s.tagServer.GetTagsByTagTypeAndTagIds(tagservice.All, tagIdsDto.tagIds)
		if err != nil {
			return nil, stew.Wrap(err)
		}

		// Assign Party domain model
		party := NewParty(int(pDto.id), pDto.startFrom, pDto.endTo, pDto.chatRoomId.String, members, tags)
		// Add a party to party list
		parties.Parties = append(parties.Parties, party)
	}

	return &parties, nil
}

func (s *realPartyServer) GetIsLatestPartyReviewDone(userId string) (*IsLatestReviewDone, error) {
	parties, err := s.GetLastNPartiesOfAUser(userId, 1)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	if parties == nil || parties.Parties == nil || len(parties.Parties) < 1 {
		return &IsLatestReviewDone{IsReviewDone: true}, nil
	}

	latestParty := parties.Parties[0]

	query := &ReviewMemberQuery{
		PartyID:  latestParty.PartyID,
		Reviewer: userId,
	}

	reviews, err := s.SearchPartyReviewMember(query)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	numOfPartyMember := len(latestParty.Members)
	numOfMembersInTheReview := len(reviews)

	return &IsLatestReviewDone{
		IsReviewDone: numOfPartyMember-1 == numOfMembersInTheReview,
	}, nil
}

func (s *realPartyServer) GetLastNPartiesOfAUser(userId string, n int) (*Parties, error) {
	partiesDto, err := s.partyQueryRepository.QueryPartiesWhereUserIdLastN(userId, n)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	if len(partiesDto) < 1 {
		return &Parties{
			Parties: []*Party{},
		}, nil
	}

	parties, err := s.populateIntoParties(partiesDto)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	return parties, nil
}

type ReviewMemberQuery struct {
	PartyID  int
	Reviewer string
	Reviewee string
}

func (s *realPartyServer) SearchPartyReviewMember(reviewMemberQuery *ReviewMemberQuery) ([]*PartyReviewMember, error) {
	var queryDto = &ReviewMemberQueryDto{
		partyID:  int64(reviewMemberQuery.PartyID),
		reviewer: reviewMemberQuery.Reviewer,
		reviewee: reviewMemberQuery.Reviewee,
	}
	retDtos, err := s.partyQueryRepository.QueryPartyReviewMembers(queryDto)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	var ret []*PartyReviewMember
	for _, dto := range retDtos {
		var prm = &PartyReviewMember{
			PartyID:  int(dto.partyID),
			Reviewer: dto.reviewer,
			Reviewee: dto.reviewee,
			Score:    dto.score,
			Comment:  dto.comments,
		}
		ret = append(ret, prm)
	}
	return ret, nil
}

func (s *realPartyServer) tran(txFunc func(*sql.Tx) (interface{}, error)) (data interface{}, err error) {
	// [Note] Make Tran, Rollback, Commit as Interface method for Dependency Injection
	tx, err := s.partyCommandRepository.Tran()
	if err != nil {
		return nil, stew.Wrap(err)
	}
	defer func() {
		if p := recover(); p != nil {
			s.partyCommandRepository.Rollback(tx)
			panic(p)
		} else if err != nil {
			s.partyCommandRepository.Rollback(tx)
		} else {
			err = s.partyCommandRepository.Commit(tx)
		}
	}()
	data, err = txFunc(tx)
	return
}

type PartyReviewMember struct {
	PartyID  int     `json:"party_id" validate:"required"`
	Reviewer string  `json:"reviewer" validate:"required"`
	Reviewee string  `json:"reviewee" validate:"required"`
	Score    float64 `json:"score" validate:"required"`
	Comment  string  `json:"comment" validate:"omitempty,min=0,max=300"`
}

func (s *realPartyServer) PostPartyReviewMember(reviewMember *PartyReviewMember) error {
	// Validation
	if err := Validate(reviewMember); err != nil {
		return domainerror.NewValidationError(err)
	}

	_, err := s.tran(func(tx *sql.Tx) (interface{}, error) {
		reviewMemberDto := &PartyMemberReviewDto{
			partyID:  int64(reviewMember.PartyID),
			reviewer: reviewMember.Reviewer,
			reviewee: reviewMember.Reviewee,
			score:    reviewMember.Score,
			comments: reviewMember.Comment,
		}

		err := s.partyCommandRepository.InsertPartyMemberReview(tx, reviewMemberDto)
		if err != nil {
			var repoErr RepositoryError
			if errors.As(err, &repoErr) {
				switch repoErr.(type) {
				case *DuplicatePrimaryKeyError:
					return nil, NewDuplicateReviewError(reviewMember.PartyID, reviewMember.Reviewer, reviewMember.Reviewee)
				case *NoReferenceRowError:
					return nil, NewInconsistencyReviewError(reviewMember.PartyID, reviewMember.Reviewer, reviewMember.Reviewee)
				}
			}
			return nil, stew.Wrap(err)
		}
		return nil, nil
	})
	if err != nil {
		return stew.Wrap(err)
	}

	return nil
}

func (s *realPartyServer) UpsertParties(partyModels []*PartyForCommand) error {
	_, err := s.tran(func(tx *sql.Tx) (interface{}, error) {
		if len(partyModels) < 1 {
			return nil, nil
		}
		// Delete the existing date's parties
		err := s.partyCommandRepository.DeletePartiesWithADay(tx, partyModels[0].StartFrom)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		// Register the new parties
		for _, party := range partyModels {
			memUserIDs := make([]string, 0, len(party.Members))
			for _, member := range party.Members {
				memUserIDs = append(memUserIDs, member.UserId)
			}
			partyDto := PartyCommandDto{
				startFrom:     party.StartFrom,
				endTo:         party.EndTo,
				chatRoomId:    utils.NewNullString(party.ChatRoomId),
				memberUserIDs: memUserIDs,
			}
			_, err := s.partyCommandRepository.InsertParty(tx, &partyDto)
			if err != nil {
				return nil, stew.Wrap(err)
			}
		}
		return nil, nil
	})
	if err != nil {
		return stew.Wrap(err)
	}

	return nil
}

// GenerateChatRoom generates chat room of a party in storage service for app users
func (s *realPartyServer) GenerateChatRoom(chatRoomId string) error {
	return s.chatRoomRepository.CreateChatRoom(chatRoomId)
}

func parseBeginEndDateTime(beginDateTimeStr, endDateTimeStr string) (begin *time.Time, end *time.Time, err error) {
	if beginDateTimeStr != "" {
		dt, err := time.Parse(time.RFC3339, beginDateTimeStr)
		if err != nil {
			return nil, nil, NewInvalidDateTimeFormatError(beginDateTimeStr)
		}
		begin = &dt
	}
	if endDateTimeStr != "" {
		dt, err := time.Parse(time.RFC3339, endDateTimeStr)
		if err != nil {
			return nil, nil, NewInvalidDateTimeFormatError(endDateTimeStr)
		}
		end = &dt
	}
	return begin, end, nil
}
