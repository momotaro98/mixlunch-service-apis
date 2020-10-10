package userservice

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/momotaro98/stew"

	"github.com/momotaro98/mixlunch-service-api/conventions"
	"github.com/momotaro98/mixlunch-service-api/domainerror"
	"github.com/momotaro98/mixlunch-service-api/tagservice"
	"github.com/momotaro98/mixlunch-service-api/utils"
)

type UserPublic struct {
	UserId             string                     `json:"user_id"`
	Name               string                     `json:"name"`
	Email              string                     `json:"email"`
	NickName           string                     `json:"nick_name"`
	PhotoUrl           string                     `json:"photo_url"`
	Position           string                     `json:"position"`
	AcademicBackground string                     `json:"academic_background"`
	Company            string                     `json:"company"`
	SelfIntroduction   string                     `json:"self_introduction"`
	Languages          []Language                 `json:"languages"`
	OccupationIDs      []OccupationID             `json:"occupation_ids"`
	InterestTags       []*tagservice.CategoryTags `json:"interest_tags"`
	SkillTags          []*tagservice.CategoryTags `json:"skill_tags"`
}

type User struct {
	UserId             string                     `json:"user_id"`
	Name               string                     `json:"name"`
	Email              string                     `json:"email"`
	NickName           string                     `json:"nick_name"`
	Sex                string                     `json:"sex"`
	Birthday           time.Time                  `json:"birthday"`
	PhotoUrl           string                     `json:"photo_url"`
	Location           conventions.Location       `json:"location"`
	Position           string                     `json:"position"`
	AcademicBackground string                     `json:"academic_background"`
	Company            string                     `json:"company"`
	SelfIntroduction   string                     `json:"self_introduction"`
	Languages          []Language                 `json:"languages"`
	OccupationIDs      []OccupationID             `json:"occupation_ids"`
	InterestTags       []*tagservice.CategoryTags `json:"interest_tags"`
	SkillTags          []*tagservice.CategoryTags `json:"skill_tags"`
	BlockingUsers      []string                   `json:"blocking_users"`
}

// UserForCommand is a user struct to register/update user info to DB.
// These Business validation specification: https://trello.com/c/zH2wCbgo/264-spec-for-user-input
// Validator 3rd party tag definitions: https://github.com/go-playground/validator/blob/912fcd16a8b47541d7a8a1768a3e3dc5b3dac220/translations/en/en_test.go#L35-L147
type UserForCommand struct {
	UserId             string               `json:"user_id" validate:"required"`
	Name               string               `json:"name" validate:"required,min=1,max=200"`
	Email              string               `json:"email" validate:"required,email"`
	NickName           string               `json:"nick_name" validate:"omitempty,min=1,max=50"`
	Sex                string               `json:"sex" validate:"required,oneof=1 2 9"`
	Birthday           string               `json:"birthday" validate:"required,datetime=2006-01-02"`
	PhotoUrl           string               `json:"photo_url" validate:"omitempty,url"`
	Location           conventions.Location `json:"location" validate:"required"`
	PositionId         uint8                `json:"position_id" validate:"omitempty,min=1"`
	AcademicBackground string               `json:"academic_background" validate:"omitempty,min=0,max=200"`
	Company            string               `json:"company" validate:"omitempty,min=0,max=200"`
	SelfIntroduction   string               `json:"self_introduction" validate:"required,min=30,max=500"`
	Languages          []Language           `json:"languages" validate:"required,min=1,max=100,dive,len=2"`
	OccupationIDs      []uint8              `json:"occupation_ids" validate:"required,min=1,max=100,dive,min=1"`
	InterestTagIds     []uint16             `json:"interest_tag_ids" validate:"omitempty,min=0,max=300,dive,min=1"`
	SkillTagIds        []uint16             `json:"skill_tag_ids" validate:"omitempty,min=0,max=300,dive,min=1"`
}

type (
	Language     string
	OccupationID uint16 // Issue: Somehow uint8 makes weird JSON response like "occupation_ids": "AQI=" not "occupation_ids": [1, 3]
)

type UserServer interface {
	GetUserByUserId(userId string) (*User, error)
	GetUserPublicByUserId(userId string) (*UserPublic, error)
	RegisterUser(newUser *UserForCommand) (*User, error)
	RegisterUserBlock(newUserBlock *UserBlockForCommand) ([]*UserBlockForQuery, error)
}

type realUserServer struct {
	userQueryRepository   IUserQueryRepository
	tagServer             tagservice.TagServer
	userCommandRepository IUserCommandRepository
}

func ProvideUserServer(userQueryRepository IUserQueryRepository,
	tagServer tagservice.TagServer,
	userCommandRepository IUserCommandRepository) UserServer {
	return &realUserServer{
		userQueryRepository:   userQueryRepository,
		tagServer:             tagServer,
		userCommandRepository: userCommandRepository,
	}
}

// GetUserByUserId does query User info by user ID.
// If the user is not in DB, return (nil, nil)
func (s *realUserServer) GetUserByUserId(userId string) (*User, error) {
	var user User
	uDto, err := s.userQueryRepository.QueryUserFullByUsingUserId(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, stew.Wrap(err)
		}
	}

	// Map from DTO to Service Model
	user.UserId = uDto.userId
	user.Name = uDto.name
	user.Email = uDto.email
	user.NickName = uDto.nickName.String
	user.Sex = uDto.sex
	user.Birthday = uDto.birthday
	user.PhotoUrl = uDto.photoUrl.String
	user.AcademicBackground = uDto.academicBackground.String
	user.Company = uDto.company.String
	user.SelfIntroduction = uDto.selfIntroduction.String
	user.Location = conventions.Location{Latitude: uDto.latitude, Longitude: uDto.longitude}
	user.Position = uDto.positionName.String

	// languages
	for _, lDto := range uDto.userlangs {
		user.Languages = append(user.Languages, Language(lDto))
	}

	// useroccupations
	// [Note] For now, occupations table master is not needed for API since
	//        Front side has the occupation master instead.
	//for _, uoDto := range uDto.useroccupations {
	//	user.OccupationIDs = append(user.OccupationIDs, OccupationID(uoDto))
	//}
	for _, uoDto := range uDto.useroccupations {
		user.OccupationIDs = append(user.OccupationIDs, OccupationID(uoDto))
	}

	// Query tags for user tags
	user.InterestTags, err = s.tagServer.GetTagsByTagTypeAndTagIds(tagservice.Interest, uDto.usertags)
	user.SkillTags, err = s.tagServer.GetTagsByTagTypeAndTagIds(tagservice.Skill, uDto.usertags)

	// Blocking Users
	if uDto.blockingUsers == nil {
		user.BlockingUsers = []string{}
	} else {
		user.BlockingUsers = uDto.blockingUsers
	}

	// Is review done

	return &user, nil
}

// GetUserByUserId does query User with simple model info by user ID.
// If the user is not in DB, return (nil, nil)
func (s *realUserServer) GetUserPublicByUserId(userId string) (*UserPublic, error) {
	user, err := s.GetUserByUserId(userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	var userPublic = &UserPublic{
		UserId:             user.UserId,
		Name:               user.Name,
		Email:              user.Email,
		NickName:           user.NickName,
		PhotoUrl:           user.PhotoUrl,
		Position:           user.Position,
		AcademicBackground: user.AcademicBackground,
		Company:            user.Company,
		SelfIntroduction:   user.SelfIntroduction,
		OccupationIDs:      user.OccupationIDs,
		Languages:          user.Languages,
		InterestTags:       user.InterestTags,
		SkillTags:          user.SkillTags,
	}

	return userPublic, nil
}

func (s *realUserServer) RegisterUser(newUser *UserForCommand) (*User, error) {
	// Validation
	if err := Validate(newUser); err != nil {
		return nil, domainerror.NewValidationError(err)
	}

	// Map from user domain model to DTO with Validation
	var uDto UserCommandDto
	uDto.userId = newUser.UserId
	uDto.name = newUser.Name
	uDto.email = newUser.Email
	uDto.nickName = utils.NewNullString(newUser.NickName)
	uDto.sex = newUser.Sex
	uDto.birthday, _ = utils.MakeDateTimeFromStringDate(newUser.Birthday)
	uDto.photoUrl = utils.NewNullString(newUser.PhotoUrl)
	{
		// Location
		uDto.latitude = newUser.Location.Latitude
		uDto.longitude = newUser.Location.Longitude
	}
	uDto.positionId = utils.NewNullInt32(int32(newUser.PositionId))
	uDto.academicBackground = utils.NewNullString(newUser.AcademicBackground)
	uDto.company = utils.NewNullString(newUser.Company)
	uDto.selfIntroduction = utils.NewNullString(newUser.SelfIntroduction)
	{
		userlangs := make([]string, 0, len(newUser.Languages))
		for _, l := range newUser.Languages {
			userlangs = append(userlangs, string(l))
		}
		uDto.userlangs = userlangs
	}
	uDto.occupationIDs = newUser.OccupationIDs
	{
		// InterestTags and SkillTags
		var usertags []uint16
		for _, tagId := range newUser.InterestTagIds {
			usertags = append(usertags, tagId)
		}
		for _, tagId := range newUser.SkillTagIds {
			usertags = append(usertags, tagId)
		}
		uDto.usertags = usertags
	}

	// Add the new user into DB
	err := s.userCommandRepository.InsertUserInfo(&uDto)
	if err != nil {
		var repoErr RepositoryError
		if errors.As(err, &repoErr) {
			switch repoErr.(type) {
			case *DuplicatePrimaryKeyError:
				return nil, NewDuplicateUserRegisterError(newUser.UserId)
			}
		}
		return nil, stew.Wrap(err)
	}

	// Query the new user
	if registeredUser, err := s.GetUserByUserId(uDto.userId); err == nil {
		if registeredUser == nil { // No user is too irregular case since registering succeeded
			return nil, errors.New(fmt.Sprintf("User was created but could not get the information. user id: %s", uDto.userId))
		}
		return registeredUser, err
	}
	return nil, stew.Wrap(err)
}

type UserBlockForCommand struct {
	Blocker string `json:"blocker" validate:"required"`
	Blockee string `json:"blockee" validate:"required"`
}

type UserBlockForQuery struct {
	Blocker   string    `json:"blocker"`
	Blockee   string    `json:"blockee"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *realUserServer) RegisterUserBlock(newUserBlock *UserBlockForCommand) ([]*UserBlockForQuery, error) {
	// Validation
	if err := Validate(newUserBlock); err != nil {
		return nil, domainerror.NewValidationError(err)
	}

	// Map from user domain model to DTO with Validation
	var ubDto = UserBlockCommandDto{
		blocker: newUserBlock.Blocker,
		blockee: newUserBlock.Blockee,
	}

	if err := s.userCommandRepository.InsertUserBlock(&ubDto); err != nil {
		var repoErr RepositoryError
		if errors.As(err, &repoErr) {
			switch repoErr.(type) {
			case *DuplicatePrimaryKeyError:
				return nil, NewDuplicateUserBlockRegisterError(newUserBlock.Blocker, newUserBlock.Blockee)
			case *NoReferenceRowError:
				return nil, NewInconsistencyUserBlockError(newUserBlock.Blocker, newUserBlock.Blockee)
			}
		}
		return nil, stew.Wrap(err)
	}

	ubOfTheBlocker, err := s.userQueryRepository.QueryUserBlockWhereBlocker(newUserBlock.Blocker)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	for _, ubb := range ubOfTheBlocker {
		if ubb.blockee == newUserBlock.Blockee {
			return []*UserBlockForQuery{
				{
					Blocker:   ubb.blocker,
					Blockee:   ubb.blockee,
					CreatedAt: ubb.createdAt,
				},
			}, nil
		}
	}

	return []*UserBlockForQuery{}, nil
}
