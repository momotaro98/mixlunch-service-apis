package userservice

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/momotaro98/mixlunch-service-api/conventions"
	"github.com/momotaro98/mixlunch-service-api/domainerror"
	"github.com/momotaro98/mixlunch-service-api/tagservice"
	mock "github.com/momotaro98/mixlunch-service-api/userservice/testmock"
)

// [Note] I had to include repositories_mock in userservice package with mockgen's `self_package` flag to avoid `import cycle` issue.
//go:generate mockgen -source=repositories.go -destination=repositories_mock.go -package=userservice -self_package=github.com/momotaro98/mixlunch-service-api/userservice
//go:generate mockgen -source=../tagservice/domain.go -destination=testmock/tagservice.go -package=testmock

const (
	uid      = "USER_ID"
	English  = "en"
	Japanese = "ja"
)

func genRegularUserForCommand() *UserForCommand {
	return &UserForCommand{
		UserId:   uid,
		Name:     "David John",
		Email:    "a@a.com",
		NickName: "Josh",
		Sex:      "1",
		Birthday: "1992-04-04",
		PhotoUrl: "https://s3.aws.com/john199",
		Location: conventions.Location{
			Latitude:  35.681236,
			Longitude: 139.767125,
		},
		PositionId:         1,
		AcademicBackground: "Tokyo University",
		Company:            "Microsoft, Inc.",
		SelfIntroduction:   "Hello, I'm John. Nice to meet you! I look forward to seeing you guys!",
		Languages:          []Language{"en", "fr"},
		OccupationIDs:      []uint8{1, 2},
		InterestTagIds:     []uint16{1, 2, 3},
		SkillTagIds:        []uint16{8, 9, 10},
	}
}

var regularUserFullQueryDto = &UserFullQueryDto{
	userId:    uid,
	userlangs: []string{English, Japanese},
}

func TestGetUserByUserId_RegularAllCase_Success(t *testing.T) {
	// Arrange
	// mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// Mock for query repository
	userQueryRepositoryMock := NewMockIUserQueryRepository(mockCtrl)
	userQueryRepositoryMock.EXPECT().
		QueryUserFullByUsingUserId(uid).
		Return(regularUserFullQueryDto, nil)
	// Mock for tag service
	tagServerMock := mock.NewMockTagServer(mockCtrl)
	// GetTagsByTagTypeAndTagIds method is called twice in GeUserByUserId method
	gomock.InOrder(
		tagServerMock.EXPECT().GetTagsByTagTypeAndTagIds(tagservice.Interest, nil).
			Return([]*tagservice.CategoryTags{}, nil),
		tagServerMock.EXPECT().GetTagsByTagTypeAndTagIds(tagservice.Skill, nil).
			Return([]*tagservice.CategoryTags{}, nil),
	)
	// Mock for command repository
	userCommandRepositoryMock := NewMockIUserCommandRepository(mockCtrl)
	// No method of CommandRepository is called in GetUserByUserId method

	userServer := ProvideUserServer(userQueryRepositoryMock, tagServerMock, userCommandRepositoryMock)

	// Act
	actUser, err := userServer.GetUserByUserId(uid)
	// Assert
	if err != nil {
		t.Failed()
	}
	if act := actUser.UserId; act != uid {
		t.Errorf("Test failed. Expected: %s, Actual: %s", uid, act)
	}
	if act := actUser.Languages; !reflect.DeepEqual([]Language{English, Japanese}, act) {
		t.Errorf("Test failed. Expected: %s, Actual: %s", uid, act)
	}
}

func TestRegisterUser_ValidateErrors(t *testing.T) {
	// Prepare mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// No mock methods are called because this test case is for validation before calling mock methods
	userQueryRepositoryMock := NewMockIUserQueryRepository(mockCtrl)
	tagServerMock := mock.NewMockTagServer(mockCtrl)
	userCommandRepositoryMock := NewMockIUserCommandRepository(mockCtrl)
	userServer := ProvideUserServer(userQueryRepositoryMock, tagServerMock, userCommandRepositoryMock)

	testValidate := func(t *testing.T, userServer UserServer, user *UserForCommand) {
		// Act
		_, err := userServer.RegisterUser(user)
		// Assert
		if err == nil {
			t.Errorf("Test failed. Expected: There's an error.', Actual: No error")
		}
		if _, ok := err.(*domainerror.ValidationError); !ok {
			t.Errorf("Test failed. Expected: ValidationError', Actual: %v", err)
		}
	}

	t.Run("'Name' Longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Name = strings.Repeat("A", 201) // Spec
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'Email' Invalid email format", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Email = "test@@example.com"
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'NickName' Longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.NickName = strings.Repeat("A", 51) // Spec
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'Sex' Out of scope", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Sex = "3"
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'Birthday' Invalid format", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Birthday = "1991/04/30"
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'PhotoUrl' Invalid format", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.PhotoUrl = "Invalid"
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'Location' latitude is longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Location = conventions.Location{
			Latitude:  90.1,       // invalid
			Longitude: 139.767125, // valid
		}
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'Location' latitude is shorter than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Location = conventions.Location{
			Latitude:  -90.1,      // invalid
			Longitude: 139.767125, // valid
		}
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'Location' longitude is longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Location = conventions.Location{
			Latitude:  34.1234,    // valid
			Longitude: 180.767125, // invalid
		}
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'Location' longitude is shorter than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Location = conventions.Location{
			Latitude:  34.1234,     // valid
			Longitude: -180.767125, // invalid
		}
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'AcademicBackground' is longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.AcademicBackground = strings.Repeat("A", 201) // Spec
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'Company' is longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Company = strings.Repeat("A", 201) // Spec
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'SelfIntroduction' is longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.SelfIntroduction = strings.Repeat("A", 501) // Spec
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'Languages' is longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		langs := make([]Language, 101)
		for i := 0; i < len(langs); i++ {
			langs[i] = Language(strings.Repeat("A", 2))
		}
		user.Languages = langs
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'Languages' content is shorter than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Languages = []Language{"e", "fr"}
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'Languages' content longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Languages = []Language{"en", "frn"}
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'OccupationIDs' is empty", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.OccupationIDs = []uint8{}
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'OccupationIDs' is longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		ids := make([]uint8, 101)
		for i := 0; i < len(ids); i++ {
			ids[i] = uint8(1)
		}
		user.OccupationIDs = ids
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'OccupationIDs' content is smaller than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.OccupationIDs = []uint8{0, 100}
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'InterestTagIds' is longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		ids := make([]uint16, 301)
		for i := 0; i < len(ids); i++ {
			ids[i] = uint16(1)
		}
		user.InterestTagIds = ids
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'InterestTagIds' content is smaller than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.InterestTagIds = []uint16{0, 100}
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'SkillTagIds' is longer than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		ids := make([]uint16, 301)
		for i := 0; i < len(ids); i++ {
			ids[i] = uint16(1)
		}
		user.SkillTagIds = ids
		// Act and Assert
		testValidate(t, userServer, user)
	})
	t.Run("'SkillTagIds' content is smaller than expected", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.SkillTagIds = []uint16{0, 100}
		// Act and Assert
		testValidate(t, userServer, user)
	})
}

func TestRegisterUser_ValidateAsRequired(t *testing.T) {
	// Prepare mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// No mock methods are called because this test case is for validation before calling mock methods
	userQueryRepositoryMock := NewMockIUserQueryRepository(mockCtrl)
	tagServerMock := mock.NewMockTagServer(mockCtrl)
	userCommandRepositoryMock := NewMockIUserCommandRepository(mockCtrl)
	userServer := ProvideUserServer(userQueryRepositoryMock, tagServerMock, userCommandRepositoryMock)

	testValidateAsRequired := func(t *testing.T, userServer UserServer, user *UserForCommand) {
		// Act
		_, err := userServer.RegisterUser(user)
		// Assert
		if err == nil {
			t.Errorf("Test failed. Expected: There's an error.', Actual: No error")
		}
		if _, ok := err.(*domainerror.ValidationError); !ok {
			t.Errorf("Test failed. Expected: ValidationError', Actual: %v", err)
		}
	}

	t.Run("No UserId", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.UserId = ""
		// Act and Assert
		testValidateAsRequired(t, userServer, user)
	})
	t.Run("No Name", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Name = ""
		// Act and Assert
		testValidateAsRequired(t, userServer, user)
	})
	t.Run("No Email", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Email = ""
		// Act and Assert
		testValidateAsRequired(t, userServer, user)
	})
	t.Run("No Birthday", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Birthday = ""
		// Act and Assert
		testValidateAsRequired(t, userServer, user)
	})
	t.Run("No Sex", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Sex = ""
		// Act and Assert
		testValidateAsRequired(t, userServer, user)
	})
	t.Run("No Location", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Location = conventions.Location{}
		// Act and Assert
		testValidateAsRequired(t, userServer, user)
	})
	t.Run("No Self Introduction", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.SelfIntroduction = ""
		// Act and Assert
		testValidateAsRequired(t, userServer, user)
	})
	t.Run("No Languages", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.Languages = []Language{}
		// Act and Assert
		testValidateAsRequired(t, userServer, user)
	})
	t.Run("No OccupationIDs", func(t *testing.T) {
		// Arrange
		user := genRegularUserForCommand()
		user.OccupationIDs = []uint8{}
		// Act and Assert
		testValidateAsRequired(t, userServer, user)
	})
}

func TestRegisterUser_MapToDto(t *testing.T) {
	// Prepare mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// No mock methods are called because this test case is for validation before calling mock methods
	userQueryRepositoryMock := NewMockIUserQueryRepository(mockCtrl)
	userQueryRepositoryMock.EXPECT().
		QueryUserFullByUsingUserId(gomock.Any()).
		Return(&UserFullQueryDto{}, nil)
	tagServerMock := mock.NewMockTagServer(mockCtrl)
	tagServerMock.EXPECT().
		GetTagsByTagTypeAndTagIds(gomock.Any(), gomock.Any()).
		Return([]*tagservice.CategoryTags{
			{}, {},
		}, nil).
		AnyTimes()
	userCommandRepositoryMock := NewMockIUserCommandRepository(mockCtrl)
	var actUserCommand *UserCommandDto
	userCommandRepositoryMock.EXPECT().
		InsertUserInfo(gomock.Any()).
		Do(func(uDto *UserCommandDto) {
			actUserCommand = uDto
		}).
		Return(nil)
	userServer := ProvideUserServer(userQueryRepositoryMock, tagServerMock, userCommandRepositoryMock)
	// Act
	_, err := userServer.RegisterUser(genRegularUserForCommand())
	// Assert
	if err != nil {
		t.Errorf("expected: %+v, got: %+v", nil, err)
	}
	expectedUser := genRegularUserForCommand()
	if actUserCommand.userId != expectedUser.UserId ||
		actUserCommand.name != expectedUser.Name ||
		actUserCommand.email != expectedUser.Email ||
		actUserCommand.latitude != expectedUser.Location.Latitude ||
		actUserCommand.longitude != expectedUser.Location.Longitude ||
		actUserCommand.sex != expectedUser.Sex {
		// @@@Need@@@: Add other fields with check logic
		t.Errorf("expected: %+v, got: %+v", expectedUser, actUserCommand)
	}
}

func TestRegisterUserBlock(t *testing.T) {
	var (
		aBlockee                 = "blockee-user"
		aCreatedAt               = time.Date(2019, 3, 11, 1, 2, 3, 4, time.Local)
		regularUserBlockForQuery = []*UserBlockQueryDto{
			{
				blocker:   uid,
				blockee:   "user-blocked-1",
				createdAt: time.Date(2020, 7, 17, 1, 2, 3, 4, time.Local),
			},
			{
				blocker:   uid,
				blockee:   aBlockee,
				createdAt: aCreatedAt,
			},
			{
				blocker:   uid,
				blockee:   "user-blocked-2",
				createdAt: time.Date(2020, 6, 11, 1, 2, 3, 4, time.Local),
			},
		}
	)

	// Prepare mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("success", func(t *testing.T) {
		// Arrange
		// Mock
		userQueryRepositoryMock := NewMockIUserQueryRepository(mockCtrl)
		userQueryRepositoryMock.EXPECT().
			QueryUserBlockWhereBlocker(uid).
			Return(regularUserBlockForQuery, nil)
		userCommandRepositoryMock := NewMockIUserCommandRepository(mockCtrl)
		userCommandRepositoryMock.EXPECT().InsertUserBlock(&UserBlockCommandDto{
			blocker: uid,
			blockee: aBlockee,
		}).Return(nil)
		userServer := ProvideUserServer(userQueryRepositoryMock, mock.NewMockTagServer(mockCtrl), userCommandRepositoryMock)
		// Input
		var inputUserBlock = &UserBlockForCommand{
			Blocker: uid,
			Blockee: aBlockee,
		}
		// Act
		ret, err := userServer.RegisterUserBlock(inputUserBlock)
		// Assert
		if err != nil {
			t.Errorf("expected: nil, actual: %+v", err)
		}
		if len(ret) != 1 {
			t.Errorf("expected: 1, actual: %d", len(ret))
		}
		if ret[0].Blocker != uid || ret[0].Blockee != aBlockee || ret[0].CreatedAt != aCreatedAt {
			t.Errorf("expected: %+v, actual: %+v,",
				&UserBlockForQuery{
					Blocker:   uid,
					Blockee:   aBlockee,
					CreatedAt: aCreatedAt,
				}, ret)
		}
	})

	testValidateAsRequired := func(t *testing.T, userServer UserServer, input *UserBlockForCommand) {
		// Act
		_, err := userServer.RegisterUserBlock(input)
		// Assert
		if err == nil {
			t.Errorf("Test failed. Expected: There's an error.', Actual: No error")
		}
		if _, ok := err.(*domainerror.ValidationError); !ok {
			t.Errorf("Test failed. Expected: ValidationError', Actual: %v", err)
		}
	}

	t.Run("blocker is not set. required validation error", func(t *testing.T) {
		// Arrange
		// Mock
		userServer := ProvideUserServer(
			NewMockIUserQueryRepository(mockCtrl),
			mock.NewMockTagServer(mockCtrl),
			NewMockIUserCommandRepository(mockCtrl),
		)
		// Input
		var input = &UserBlockForCommand{
			Blockee: aBlockee,
		}
		// Act and Assert
		testValidateAsRequired(t, userServer, input)
	})

	t.Run("blockee is not set. required validation error", func(t *testing.T) {
		// Arrange
		// Mock
		userServer := ProvideUserServer(
			NewMockIUserQueryRepository(mockCtrl),
			mock.NewMockTagServer(mockCtrl),
			NewMockIUserCommandRepository(mockCtrl),
		)
		// Input
		var input = &UserBlockForCommand{
			Blocker: uid,
		}
		// Act and Assert
		testValidateAsRequired(t, userServer, input)
	})
}
