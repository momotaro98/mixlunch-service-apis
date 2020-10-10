package userscheduleservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/momotaro98/stew"

	"github.com/momotaro98/mixlunch-service-api/conventions"
	"github.com/momotaro98/mixlunch-service-api/domainerror"
	"github.com/momotaro98/mixlunch-service-api/tagservice"
)

type UserSchedules struct {
	UserId        string          `json:"user_id"`
	UserSchedules []*UserSchedule `json:"user_schedules"`
}

type UserSchedule struct {
	FromDateTime time.Time                  `json:"from_date_time"`
	ToDateTime   time.Time                  `json:"to_date_time"`
	Tags         []*tagservice.CategoryTags `json:"tags"`
	Location     conventions.Location       `json:"location"`
}

func NewUserSchedule(
	fromDateTime, toDateTime time.Time,
	tags []*tagservice.CategoryTags,
	location conventions.Location,
) *UserSchedule {
	return &UserSchedule{
		FromDateTime: fromDateTime,
		ToDateTime:   toDateTime,
		Tags:         tags,
		Location:     location,
	}
}

// UserScheduleForCommand is a user schedule struct to register/update.
// Validator 3rd party tag definitions: https://github.com/go-playground/validator/blob/912fcd16a8b47541d7a8a1768a3e3dc5b3dac220/translations/en/en_test.go#L35-L147
type UserScheduleForCommand struct {
	FromDateTime time.Time            `json:"from_date_time" validate:"required"`
	ToDateTime   time.Time            `json:"to_date_time" validate:"required"`
	TagIds       []uint16             `json:"tag_ids" validate:"required,min=0,dive,min=1"`
	Location     conventions.Location `json:"location" validate:"required"`
}

type SpecifiedDate struct {
	Date time.Time `json:"date"`
}

func validateFromDateTimeAndToDateTime(fromDateTime, toDateTime time.Time) error {
	// Validate if fromDateTime is before than toDateTime
	if !fromDateTime.Before(toDateTime) {
		return NewFromIsAfterToError(fromDateTime, toDateTime)
	}
	// Validate if fromDateTime and toDateTime are not in a same day
	if fromDateTime.Year() != toDateTime.Year() ||
		fromDateTime.Month() != toDateTime.Month() ||
		fromDateTime.Day() != toDateTime.Day() {
		return NewDifferentDayFromAndToError(fromDateTime, toDateTime)
	}
	// [Business requirement] Validate time range between fromDateTime and toDateTime should be within specified range
	if toDateTime.Sub(fromDateTime).Minutes() < regulatedTimeDurationMinutes {
		return NewTimeRangeIsLessThanSpecifiedError(fromDateTime, toDateTime)
	}
	return nil
}

type UserScheduleServer interface {
	GetUserSchedulesByTimeRange(userId, beginDateTimeStr, endDateTimeStr string) (*UserSchedules, error)
	GetEachUserSchedules(beginDateTimeStr, endDateTimeStr string) ([]*UserSchedules, error)
	AddUserSchedule(userId string, usComm *UserScheduleForCommand) (*UserSchedules, error)
	UpdateUserSchedule(userId string, usComm *UserScheduleForCommand) (*UserSchedules, error)
	DeleteUserSchedule(userId string, targetDate time.Time) (*UserSchedules, error)
}

type realUserScheduleServer struct {
	userScheduleQueryRepository   IUserScheduleQueryRepository
	userScheduleCommandRepository IUserScheduleCommandRepository
	tagServer                     tagservice.TagServer
}

func ProvideUserScheduleServer(queryRepository IUserScheduleQueryRepository, updateRepository IUserScheduleCommandRepository, tagServer tagservice.TagServer) UserScheduleServer {
	return &realUserScheduleServer{
		userScheduleQueryRepository:   queryRepository,
		userScheduleCommandRepository: updateRepository,
		tagServer:                     tagServer,
	}
}

// extractAScheduleDtoWithValidation returns DTO of one user schedule
// with validation to check if there is only one user schedule in the specified day.
func (s *realUserScheduleServer) extractAScheduleDtoWithValidation(userId string, dateTimeOfTheTargetDate time.Time) (*UserScheduleDto, error) {
	begin := time.Date(dateTimeOfTheTargetDate.Year(), dateTimeOfTheTargetDate.Month(), dateTimeOfTheTargetDate.Day(),
		0, 0, 0, 0, time.Local)
	end := time.Date(dateTimeOfTheTargetDate.Year(), dateTimeOfTheTargetDate.Month(), dateTimeOfTheTargetDate.Day(),
		23, 59, 59, 999, time.Local)
	// Try to retrieve the target user schedule and Validate
	queriedDtosToValidateAndSpecify, err := s.userScheduleQueryRepository.QueryUserSchedulesWhereTimeRange(begin, end, userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	if len(queriedDtosToValidateAndSpecify) > 1 {
		// Business logic error. A day must not have more than one user schedule.
		return nil, fmt.Errorf("[Business Error] There are more than one user schedules in the day. The date: %v", dateTimeOfTheTargetDate)
	}
	if len(queriedDtosToValidateAndSpecify) == 0 {
		return nil, NewTheScheduleNotFoundError(dateTimeOfTheTargetDate)
	}
	return queriedDtosToValidateAndSpecify[0], nil
}

func (s *realUserScheduleServer) GetUserSchedulesByTimeRange(userId, beginDateTimeStr, endDateTimeStr string) (*UserSchedules, error) {
	// Parse DateTime string to RFC3339 spec
	beginDateTime, err := time.Parse(time.RFC3339, beginDateTimeStr)
	if err != nil {
		return nil, NewInvalidDateTimeFormatError(beginDateTimeStr)
	}
	endDateTime, err := time.Parse(time.RFC3339, endDateTimeStr)
	if err != nil {
		return nil, NewInvalidDateTimeFormatError(endDateTimeStr)
	}

	// Query Dto from DB through repository layer
	dtos, err := s.userScheduleQueryRepository.QueryUserSchedulesWhereTimeRange(beginDateTime, endDateTime, userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	var uSchedules = UserSchedules{
		UserId:        userId,
		UserSchedules: make([]*UserSchedule, 0),
	}
	if len(dtos) < 1 {
		return &uSchedules, nil
	}

	// Assign to domain
	for _, dto := range dtos {
		// Query tags from tagservice
		tags, err := s.tagServer.GetTagsByTagTypeAndTagIds(tagservice.All, dto.tagIds)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		// New user schedule model
		us := NewUserSchedule(
			dto.fromDateTime, dto.toDateTime,
			tags,
			conventions.NewLocation(dto.latitude, dto.longitude, dto.locationTypeID),
		)
		uSchedules.UserSchedules = append(uSchedules.UserSchedules, us)
	}
	return &uSchedules, nil
}

func (s *realUserScheduleServer) GetEachUserSchedules(beginDateTimeStr, endDateTimeStr string) ([]*UserSchedules, error) {
	// Parse DateTime string to RFC3339 spec
	beginDateTime, err := time.Parse(time.RFC3339, beginDateTimeStr)
	if err != nil {
		return nil, NewInvalidDateTimeFormatError(beginDateTimeStr)
	}
	endDateTime, err := time.Parse(time.RFC3339, endDateTimeStr)
	if err != nil {
		return nil, NewInvalidDateTimeFormatError(endDateTimeStr)
	}
	// Query Dto from DB through repository layer
	dtos, err := s.userScheduleQueryRepository.QueryUserSchedulesWhereTimeRange(beginDateTime, endDateTime, "")
	if err != nil {
		return nil, stew.Wrap(err)
	}
	/* Expected flow

	input(dtos) =>
	{userID:kent, fromdt:a, todt:b},
	{userID:kent, fromdt:c, todt:d},
	{userID:nick, fromdt:a, todt:b},
	{userID:john, fromdt:a, todt:b},
	{userID:john, fromdt:c, todt:d},

	output =>
	{[
		{userID:kent, UserSchedules: [{fdt:c, tdt:d}, {fdt:a, tdt:b}]},
		{userID:nick, UserSchedules: [{fdt:c, tdt:d}]},
		{userID:john, UserSchedules: [{fdt:c, tdt:d}, {fdt:a, tdt:b}]},
	]}

	*/
	var ret []*UserSchedules
	if len(dtos) < 1 {
		return ret, nil
	}
	var currentUserSchedule UserSchedules
	for _, dto := range dtos {
		// Query tags from tagservice
		tags, err := s.tagServer.GetTagsByTagTypeAndTagIds(tagservice.All, dto.tagIds)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		// New user schedule model
		uSchedule := NewUserSchedule(
			dto.fromDateTime, dto.toDateTime,
			tags,
			conventions.NewLocation(dto.latitude, dto.longitude, dto.locationTypeID),
		)

		if currentUserSchedule.UserId == "" {
			// initialize
			currentUserSchedule = UserSchedules{
				UserId:        dto.userId,
				UserSchedules: []*UserSchedule{uSchedule},
			}
		} else if dto.userId == currentUserSchedule.UserId {
			// Append the user schedule to current slice
			currentUserSchedule.UserSchedules = append(currentUserSchedule.UserSchedules, uSchedule)
		} else {
			// flush current stacks
			copiedSchedule := currentUserSchedule
			ret = append(ret, &copiedSchedule)
			// refresh
			currentUserSchedule = UserSchedules{
				UserId:        dto.userId,
				UserSchedules: []*UserSchedule{uSchedule},
			}
		}
	}
	// flush last stacks
	copiedSchedule := currentUserSchedule
	ret = append(ret, &copiedSchedule)

	return ret, nil
}

func (s *realUserScheduleServer) AddUserSchedule(userId string, usComm *UserScheduleForCommand) (*UserSchedules, error) {
	// Validation
	if err := ValidateUserSchedule(usComm); err != nil {
		return nil, domainerror.NewValidationError(err)
	}

	// Datetime validation
	if err := validateFromDateTimeAndToDateTime(usComm.FromDateTime, usComm.ToDateTime); err != nil {
		return nil, err
	}

	// Validate if non duplicate user schedule in one day by querying the table
	beginDateTime := time.Date(
		usComm.FromDateTime.Year(), usComm.FromDateTime.Month(), usComm.FromDateTime.Day(),
		0, 0, 0, 0, time.Local)
	endDateTime := time.Date(
		usComm.FromDateTime.Year(), usComm.FromDateTime.Month(), usComm.FromDateTime.Day(),
		23, 59, 59, 999, time.Local)
	queriedDtosToValidate, err := s.userScheduleQueryRepository.QueryUserSchedulesWhereTimeRange(beginDateTime, endDateTime, userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	if len(queriedDtosToValidate) > 0 {
		return nil, NewDuplicateInOneDayError(beginDateTime)
	}

	// Add a user schedule
	lastInsertedId, err := s.userScheduleCommandRepository.InsertUserSchedule(
		NewUserScheduleDtoForCommand(userId,
			usComm.FromDateTime, usComm.ToDateTime,
			usComm.TagIds,
			usComm.Location.Latitude, usComm.Location.Longitude,
		),
	)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	// Get the newly added user schedule to return
	lastInsertedDto, err := s.userScheduleQueryRepository.QueryUserScheduleWhereId(lastInsertedId)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	if lastInsertedDto == nil {
		return nil, errors.New("can't get the inserted user schedule in AddUserSchedule")
	}

	var uSchedules UserSchedules // variable to return
	// UserId
	uSchedules.UserId = lastInsertedDto.userId
	// Tags of the schedule
	tags, err := s.tagServer.GetTagsByTagTypeAndTagIds(tagservice.All, lastInsertedDto.tagIds)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	// This service method should make sure that user schedules has only one model.
	oneUserSchedule := NewUserSchedule(
		lastInsertedDto.fromDateTime, lastInsertedDto.toDateTime,
		tags,
		conventions.NewLocation(lastInsertedDto.latitude, lastInsertedDto.longitude, lastInsertedDto.locationTypeID),
	)
	uSchedules.UserSchedules = append(uSchedules.UserSchedules, oneUserSchedule)
	return &uSchedules, nil
}

func (s *realUserScheduleServer) UpdateUserSchedule(userId string, usComm *UserScheduleForCommand) (*UserSchedules, error) {
	// Validation
	if err := ValidateUserSchedule(usComm); err != nil {
		return nil, domainerror.NewValidationError(err)
	}

	// Datetime validation
	if err := validateFromDateTimeAndToDateTime(usComm.FromDateTime, usComm.ToDateTime); err != nil {
		return nil, err
	}

	// Check if there is a user schedule in the day
	targetDtoToUpdate, err := s.extractAScheduleDtoWithValidation(userId, usComm.FromDateTime)
	if err != nil {
		return nil, err
	}

	// Update the target user schedule
	lastUpdatedId, err := s.userScheduleCommandRepository.UpdateUserSchedule(
		targetDtoToUpdate.userScheduleId,
		NewUserScheduleDtoForCommand(userId,
			usComm.FromDateTime, usComm.ToDateTime,
			usComm.TagIds,
			usComm.Location.Latitude, usComm.Location.Longitude,
		),
	)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	// Get the updated user schedule to return
	lastUpdatedDto, err := s.userScheduleQueryRepository.QueryUserScheduleWhereId(lastUpdatedId)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	if lastUpdatedDto == nil {
		return nil, errors.New("error in UpdateUserSchedule: can't get the inserted user schedule")
	}

	var uSchedules UserSchedules // variable to return
	uSchedules.UserId = lastUpdatedDto.userId
	// Tags of the schedule
	tags, err := s.tagServer.GetTagsByTagTypeAndTagIds(tagservice.All, lastUpdatedDto.tagIds)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	// This service method should make sure that user schedules has only one model.
	oneUserSchedule := NewUserSchedule(
		lastUpdatedDto.fromDateTime, lastUpdatedDto.toDateTime,
		tags,
		conventions.NewLocation(lastUpdatedDto.latitude, lastUpdatedDto.longitude, lastUpdatedDto.locationTypeID),
	)
	uSchedules.UserSchedules = append(uSchedules.UserSchedules, oneUserSchedule)
	return &uSchedules, nil
}

func (s *realUserScheduleServer) DeleteUserSchedule(userId string, targetDate time.Time) (*UserSchedules, error) {
	// Check if there is a user schedule in the day
	targetDtoToDelete, err := s.extractAScheduleDtoWithValidation(userId, targetDate)
	if err != nil {
		return nil, err
	}

	// Query Tags information before deleting
	tags, err := s.tagServer.GetTagsByTagTypeAndTagIds(tagservice.All, targetDtoToDelete.tagIds)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	// Delete the user schedule
	err = s.userScheduleCommandRepository.DeleteUserSchedule(targetDtoToDelete.userScheduleId)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	var uSchedules UserSchedules // variable to return
	uSchedules.UserId = targetDtoToDelete.userId
	// This service method should make sure that user schedules has only one model.
	oneUserSchedule := NewUserSchedule(
		targetDtoToDelete.fromDateTime, targetDtoToDelete.toDateTime,
		tags,
		conventions.NewLocation(targetDtoToDelete.latitude, targetDtoToDelete.longitude, targetDtoToDelete.locationTypeID),
	)
	uSchedules.UserSchedules = append(uSchedules.UserSchedules, oneUserSchedule)
	return &uSchedules, nil
}
