package userscheduleservice

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/momotaro98/mixlunch-service-api/conventions"
	"github.com/momotaro98/mixlunch-service-api/tagservice"
	"github.com/momotaro98/mixlunch-service-api/userscheduleservice/testmock"
	"github.com/momotaro98/mixlunch-service-api/utils"
)

//go:generate mockgen -source=repositories.go -destination=repositories_mock.go -package=userscheduleservice -self_package=github.com/momotaro98/mixlunch-service-api/userscheduleservice
//go:generate mockgen -source=../tagservice/domain.go -destination=testmock/tagservice.go -package=testmock

const (
	uid       = "userId0123456789"
	userId2   = "john1234"
	userId3   = "denny5678"
	baseYear  = 2018
	baseMonth = 11
	baseDay   = 15
	baseHour  = 12
	anyInt64  = int64(10000001)
)

var (
	tagIDs   = []uint16{1, 3, 5}
	location = conventions.Location{
		Latitude:  35.0,
		Longitude: 135.0,
	}
)

// Helpers for tests

func makeOnlyOneUserScheduleDtos(userId string, id int64) []*UserScheduleDto {
	fromDateTime1 := time.Date(baseYear, baseMonth, baseDay, baseHour-1, 0, 0, 0, time.Local)
	toDateTime1 := time.Date(baseYear, baseMonth, baseDay, baseHour+1, 0, 0, 0, time.Local)
	var dto1 = UserScheduleDto{userScheduleId: id, userId: userId, fromDateTime: fromDateTime1, toDateTime: toDateTime1}
	return []*UserScheduleDto{&dto1}
}

func makeSomeUserScheduleDtos(userId string) []*UserScheduleDto {
	fromDateTime1 := time.Date(baseYear, baseMonth, baseDay, baseHour-1, 0, 0, 0, time.Local)
	toDateTime1 := time.Date(baseYear, baseMonth, baseDay, baseHour+1, 0, 0, 0, time.Local)
	var dto1 = UserScheduleDto{userScheduleId: 100, userId: userId, fromDateTime: fromDateTime1, toDateTime: toDateTime1}
	fromDateTime2 := time.Date(baseYear, baseMonth, baseDay+3, baseHour-1, 0, 0, 0, time.Local)
	toDateTime2 := time.Date(baseYear, baseMonth, baseDay+3, baseHour, 30, 0, 0, time.Local)
	var dto2 = UserScheduleDto{userScheduleId: 101, userId: userId, fromDateTime: fromDateTime2, toDateTime: toDateTime2}
	return []*UserScheduleDto{&dto1, &dto2}
}

func makeOneUserScheduleDto() *UserScheduleDto {
	fromDateTime1 := time.Date(baseYear, baseMonth, baseDay, baseHour-1, 0, 0, 0, time.Local)
	toDateTime1 := time.Date(baseYear, baseMonth, baseDay, baseHour+1, 0, 0, 0, time.Local)
	var dto = UserScheduleDto{userScheduleId: 999999999, userId: uid, fromDateTime: fromDateTime1, toDateTime: toDateTime1}
	return &dto
}

func makeMultipleUsersSchedulesDtos() []*UserScheduleDto {
	fromDateTime1 := time.Date(baseYear, baseMonth, baseDay, baseHour-1, 0, 0, 0, time.Local)
	toDateTime1 := time.Date(baseYear, baseMonth, baseDay, baseHour+1, 0, 0, 0, time.Local)
	fromDateTime2 := time.Date(baseYear, baseMonth, baseDay+3, baseHour-1, 0, 0, 0, time.Local)
	toDateTime2 := time.Date(baseYear, baseMonth, baseDay+3, baseHour, 30, 0, 0, time.Local)
	var user1Dto1 = UserScheduleDto{userScheduleId: 101, userId: uid, fromDateTime: fromDateTime1, toDateTime: toDateTime1}
	var user1Dto2 = UserScheduleDto{userScheduleId: 102, userId: uid, fromDateTime: fromDateTime2, toDateTime: toDateTime2}
	var user2Dto1 = UserScheduleDto{userScheduleId: 103, userId: userId2, fromDateTime: fromDateTime1, toDateTime: toDateTime1}
	var user3Dto1 = UserScheduleDto{userScheduleId: 104, userId: userId3, fromDateTime: fromDateTime1, toDateTime: toDateTime1}
	var user3Dto2 = UserScheduleDto{userScheduleId: 105, userId: userId3, fromDateTime: fromDateTime2, toDateTime: toDateTime2}
	return []*UserScheduleDto{&user1Dto1, &user1Dto2, &user2Dto1, &user3Dto1, &user3Dto2}
}

func makeRegularTags() (ret []*tagservice.CategoryTags) {
	cateTags1 := &tagservice.CategoryTags{
		Category: &tagservice.Category{CategoryId: 1, Name: "cate1"},
		Tags: []*tagservice.SmallTag{
			{TagId: 1, Name: "tag1"},
			{TagId: 3, Name: "tag3"},
		},
	}
	cateTags2 := &tagservice.CategoryTags{
		Category: &tagservice.Category{CategoryId: 2, Name: "cate2"},
		Tags: []*tagservice.SmallTag{
			{TagId: 2, Name: "tag2"},
			{TagId: 4, Name: "tag4"},
		},
	}
	return append(ret, cateTags1, cateTags2)
}

func TestGetUserSchedulesByTimeRange_SomeDTOs_ReturnModels(t *testing.T) {
	// Arrange
	beginDateTime := utils.MakeCorrectFormatDateTimeStr(baseYear-1, baseMonth, 1, 0, 0)
	endDateTime := utils.MakeCorrectFormatDateTimeStr(baseYear+1, baseMonth, 30, 0, 0)

	// Mock preparation
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// query mock
	userScheduleQueryRepositoryMock := NewMockIUserScheduleQueryRepository(mockCtrl)
	userScheduleQueryRepositoryMock.EXPECT().
		QueryUserSchedulesWhereTimeRange(gomock.Any(), gomock.Any(), uid).
		Return(makeSomeUserScheduleDtos(uid), nil)
	// tag service mock
	tagServiceMock := testmock.NewMockTagServer(mockCtrl)
	tagServiceMock.EXPECT().
		GetTagsByTagTypeAndTagIds(tagservice.All, gomock.Any()).
		Return(makeRegularTags(), nil).
		AnyTimes()
	// Initialize mock
	userScheduleServer := ProvideUserScheduleServer(
		userScheduleQueryRepositoryMock,
		NewMockIUserScheduleCommandRepository(mockCtrl),
		tagServiceMock,
	)

	// Act
	uSchedules, _ := userScheduleServer.GetUserSchedulesByTimeRange(uid, beginDateTime, endDateTime)
	// Assert
	if uSchedules.UserId != uid {
		t.Errorf("Test failed. Expected: %s', Actual: %s", uid, uSchedules.UserId)
	}
	if len(uSchedules.UserSchedules) != 2 {
		t.Errorf("Test failed. Expected: 2', Actual: %d", len(uSchedules.UserSchedules))
	}
}

func TestGetUserSchedulesByTimeRange_DateTimeFormatIsIncorrect_ReturnNil(t *testing.T) {
	// Arrange : begin is correct. end is Incorrect.
	/// Business
	beginDateTime := utils.MakeCorrectFormatDateTimeStr(baseYear-1, baseMonth, 1, 0, 0)
	endDateTime := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00", baseYear+1, baseMonth, 2, 0, 0) // Incorrect format
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userScheduleServer := ProvideUserScheduleServer(
		NewMockIUserScheduleQueryRepository(mockCtrl),
		NewMockIUserScheduleCommandRepository(mockCtrl),
		testmock.NewMockTagServer(mockCtrl),
	)
	// Act
	uSchedules, err := userScheduleServer.GetUserSchedulesByTimeRange(uid, beginDateTime, endDateTime)
	// Assert
	if uSchedules != nil {
		t.Errorf("Test failed. Expected: nil', Actual: %v", uSchedules)
	}
	if err == nil {
		t.Errorf("Test failed. Expected: error has something', Actual: nil")
	}
}

func TestGetUserSchedulesByTimeRange_EmptyDots_ZeroValueUserSchedules(t *testing.T) {
	// Arrange
	/// Business
	beginDateTime := utils.MakeCorrectFormatDateTimeStr(baseYear-1, baseMonth, 1, 0, 0)
	endDateTime := utils.MakeCorrectFormatDateTimeStr(baseYear+1, baseMonth, 30, 0, 0)
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userScheduleQueryRepositoryMock := NewMockIUserScheduleQueryRepository(mockCtrl)
	userScheduleQueryRepositoryMock.EXPECT().
		QueryUserSchedulesWhereTimeRange(gomock.Any(), gomock.Any(), uid).
		Return([]*UserScheduleDto{}, nil) // Empty user is expected
	userScheduleServer := ProvideUserScheduleServer(
		userScheduleQueryRepositoryMock,
		NewMockIUserScheduleCommandRepository(mockCtrl),
		testmock.NewMockTagServer(mockCtrl),
	)
	// Act
	uSchedules, _ := userScheduleServer.GetUserSchedulesByTimeRange(uid, beginDateTime, endDateTime)
	// Assert
	if uSchedules.UserId != uid {
		t.Errorf("Test failed. Expected: %s', Actual: %s", uid, uSchedules.UserId)
	}
	if len(uSchedules.UserSchedules) != 0 {
		t.Errorf("Test failed. Expected: 0', Actual: %d", len(uSchedules.UserSchedules))
	}
}

func TestGetEachUserSchedules_RegularMultiUsersDtos_Expected(t *testing.T) {
	// Arrange
	beginDateTime := utils.MakeCorrectFormatDateTimeStr(baseYear-1, baseMonth, 1, 0, 0)
	endDateTime := utils.MakeCorrectFormatDateTimeStr(baseYear+1, baseMonth, 30, 0, 0)
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// user-schedule service mock
	userScheduleQueryRepositoryMock := NewMockIUserScheduleQueryRepository(mockCtrl)
	userScheduleQueryRepositoryMock.EXPECT().
		QueryUserSchedulesWhereTimeRange(gomock.Any(), gomock.Any(), ""). // Empty userId is expected
		Return(makeMultipleUsersSchedulesDtos(), nil)                     // some users DTO
	// tag service mock
	tagServiceMock := testmock.NewMockTagServer(mockCtrl)
	tagServiceMock.EXPECT().
		GetTagsByTagTypeAndTagIds(tagservice.All, gomock.Any()).
		Return(makeRegularTags(), nil).AnyTimes()
	userScheduleServer := ProvideUserScheduleServer(
		userScheduleQueryRepositoryMock,
		NewMockIUserScheduleCommandRepository(mockCtrl),
		tagServiceMock,
	)
	// Act
	eachUserSchedules, _ := userScheduleServer.GetEachUserSchedules(beginDateTime, endDateTime)
	// Assert
	if len(eachUserSchedules) != 3 { // since makeMultipleUsersSchedulesDtos makes 3 users schedules
		t.Errorf("Test failed. Expected: %d', Actual: %d", 3, len(eachUserSchedules))
	}
	if eachUserSchedules[0].UserId != uid {
		t.Errorf("Test failed. Expected: %s', Actual: %s", uid, eachUserSchedules[0].UserId)
	}
	if eachUserSchedules[1].UserId != userId2 {
		t.Errorf("Test failed. Expected: %s', Actual: %s", userId2, eachUserSchedules[1].UserId)
	}
	if eachUserSchedules[2].UserId != userId3 {
		t.Errorf("Test failed. Expected: %s', Actual: %s", userId3, eachUserSchedules[2].UserId)
	}
	if len(eachUserSchedules[0].UserSchedules) != 2 { // based on makeMultipleUsersSchedulesDtos implementation
		t.Errorf("Test failed. Expected: %d', Actual: %d", 2, len(eachUserSchedules[0].UserSchedules))
	}
	if len(eachUserSchedules[1].UserSchedules) != 1 {
		t.Errorf("Test failed. Expected: %d', Actual: %d", 1, len(eachUserSchedules[1].UserSchedules))
	}
	if len(eachUserSchedules[2].UserSchedules) != 2 {
		t.Errorf("Test failed. Expected: %d', Actual: %d", 2, len(eachUserSchedules[2].UserSchedules))
	}
}

func TestAddUserSchedule_EverythingIsOk_NoError(t *testing.T) {
	// Arrange
	/// Business
	fromDateTime := time.Date(baseYear, baseMonth, baseDay+5, baseHour, 30, 0, 0, time.Local)
	toDateTime := time.Date(baseYear, baseMonth, baseDay+5, baseHour+1, 30, 0, 0, time.Local)
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	//// Mock of Query repository
	lastInsertedIdOfUserSchedule := anyInt64
	userScheduleQueryRepositoryMock := NewMockIUserScheduleQueryRepository(mockCtrl)
	userScheduleQueryRepositoryMock.EXPECT().
		QueryUserSchedulesWhereTimeRange(gomock.Any(), gomock.Any(), uid).
		Return([]*UserScheduleDto{}, nil) // Empty user is expected to avoid duplicate registering error
	userScheduleQueryRepositoryMock.EXPECT().
		QueryUserScheduleWhereId(lastInsertedIdOfUserSchedule).
		Return(makeOneUserScheduleDto(), nil) // Registered schedule
	//// Mock of Command repository
	userScheduleCommandRepositoryMock := NewMockIUserScheduleCommandRepository(mockCtrl)
	userScheduleCommandRepositoryMock.EXPECT().
		InsertUserSchedule(gomock.Any()).
		Return(lastInsertedIdOfUserSchedule, nil)
	// tag service mock
	tagServiceMock := testmock.NewMockTagServer(mockCtrl)
	tagServiceMock.EXPECT().
		GetTagsByTagTypeAndTagIds(tagservice.All, gomock.Any()).
		Return(makeRegularTags(), nil)
	//// Initialize server with mocks
	userScheduleServer := ProvideUserScheduleServer(
		userScheduleQueryRepositoryMock,
		userScheduleCommandRepositoryMock,
		tagServiceMock,
	)
	// Act
	uSchedules, err := userScheduleServer.AddUserSchedule(uid,
		&UserScheduleForCommand{
			FromDateTime: fromDateTime,
			ToDateTime:   toDateTime,
			TagIds:       tagIDs,
			Location:     location,
		},
	)
	// Assert
	if err != nil {
		t.Errorf("Test failed. Expected: no error', Actual: %s", err)
	}
	if uSchedules.UserId != uid {
		t.Errorf("Test failed. Expected: %s', Actual: %s", uid, uSchedules.UserId)
	}
	if len(uSchedules.UserSchedules) != 1 {
		t.Errorf("Test failed. Expected: 1', Actual: %d", len(uSchedules.UserSchedules))
	}
}

func TestAddUserSchedule_TimeRangeInvalidCases_ReturnError(t *testing.T) {
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userScheduleServer := ProvideUserScheduleServer(
		NewMockIUserScheduleQueryRepository(mockCtrl),
		NewMockIUserScheduleCommandRepository(mockCtrl),
		testmock.NewMockTagServer(mockCtrl),
	)

	t.Run("From datetime is after than To datetime", func(t *testing.T) {
		// Arrange
		fromDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour+1, 30, 0, 0, time.Local) // After than toDateTime
		toDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour, 0, 0, 0, time.Local)
		// Act
		uSchedules, err := userScheduleServer.AddUserSchedule(uid,
			&UserScheduleForCommand{
				FromDateTime: fromDateTime,
				ToDateTime:   toDateTime,
				TagIds:       tagIDs,
				Location:     location,
			},
		)
		// Assert
		if _, ok := err.(*FromIsAfterToError); !ok {
			t.Errorf("Test failed. Expected: *FromIsAfterToError', Actual: %v", err)
		}
		if uSchedules != nil {
			t.Errorf("Test failed. Expected: nil', Actual: %v", uSchedules)
		}
	})
	t.Run("From datetime and To datetime are not in a same day", func(t *testing.T) {
		// Arrange
		fromDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour, 30, 0, 0, time.Local)
		toDateTime := time.Date(baseYear, baseMonth, baseDay+1, baseHour+1, 30, 0, 0, time.Local) // Different date from fromDateTime
		// Act
		uSchedules, err := userScheduleServer.AddUserSchedule(uid,
			&UserScheduleForCommand{
				FromDateTime: fromDateTime,
				ToDateTime:   toDateTime,
				TagIds:       tagIDs,
				Location:     location,
			},
		)
		// Assert
		if _, ok := err.(*DifferentDayFromAndToError); !ok {
			t.Errorf("Test failed. Expected: *FromIsAfterToError', Actual: %v", err)
		}
		if uSchedules != nil {
			t.Errorf("Test failed. Expected: nil', Actual: %v", uSchedules)
		}
	})
	t.Run("Time range From and To is less than business required speficication one", func(t *testing.T) {
		// Arrange
		fromDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour, 0, 0, 0, time.Local)
		toDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour, 59, 0, 0, time.Local) // Within 60 minutes
		// Act
		uSchedules, err := userScheduleServer.AddUserSchedule(uid,
			&UserScheduleForCommand{
				FromDateTime: fromDateTime,
				ToDateTime:   toDateTime,
				TagIds:       tagIDs,
				Location:     location,
			},
		)
		// Assert
		if _, ok := err.(*TimeRangeIsLessThanSpecifiedError); !ok {
			t.Errorf("Test failed. Expected: *TimeRangeIsLessThanSpecifiedError', Actual: %v", err)
		}
		if uSchedules != nil {
			t.Errorf("Test failed. Expected: nil', Actual: %v", uSchedules)
		}
	})
}

func TestAddUserSchedule_DuplicateDayInASameDay_DuplicateInOneDayError(t *testing.T) {
	// Arrange
	/// Business
	fromDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour, 30, 0, 0, time.Local)
	toDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour+1, 30, 0, 0, time.Local)
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userScheduleQueryRepository := NewMockIUserScheduleQueryRepository(mockCtrl)
	userScheduleQueryRepository.EXPECT().
		QueryUserSchedulesWhereTimeRange(gomock.Any(), gomock.Any(), uid).
		Return(makeOnlyOneUserScheduleDtos(uid, anyInt64), nil) // Duplicate in a day
	userScheduleServer := ProvideUserScheduleServer(
		userScheduleQueryRepository,
		NewMockIUserScheduleCommandRepository(mockCtrl),
		testmock.NewMockTagServer(mockCtrl),
	)
	// Act
	uSchedules, err := userScheduleServer.AddUserSchedule(uid,
		&UserScheduleForCommand{
			FromDateTime: fromDateTime,
			ToDateTime:   toDateTime,
			TagIds:       tagIDs,
			Location:     location,
		},
	)
	// Assert
	var e *DuplicateInOneDayError
	if !errors.As(err, &e) {
		t.Errorf("Test failed. Expected: wrappted *DuplicateInOneDayError', Actual: %v", err)
	}
	if uSchedules != nil {
		t.Errorf("Test failed. Expected: nil', Actual: %v", uSchedules)
	}
}

func TestUpdateUserSchedule_EverythingIsOk_NoError(t *testing.T) {
	// Arrange
	/// Business
	fromDateTime := time.Date(baseYear, baseMonth, baseDay+5, baseHour, 30, 0, 0, time.Local)
	toDateTime := time.Date(baseYear, baseMonth, baseDay+5, baseHour+1, 30, 0, 0, time.Local)
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	queriedIdOfUserSchedule := anyInt64
	lastUpdatedIdOfUserSchedule := anyInt64

	userScheduleQueryRepositoryMock := NewMockIUserScheduleQueryRepository(mockCtrl)
	userScheduleQueryRepositoryMock.EXPECT().
		QueryUserSchedulesWhereTimeRange(gomock.Any(), gomock.Any(), uid).
		Return(makeOnlyOneUserScheduleDtos(uid, queriedIdOfUserSchedule), nil) // Only one user is expected before update

	userScheduleCommandRepositoryMock := NewMockIUserScheduleCommandRepository(mockCtrl)
	userScheduleCommandRepositoryMock.EXPECT().
		UpdateUserSchedule(
			queriedIdOfUserSchedule,
			gomock.Any(),
		).
		Return(lastUpdatedIdOfUserSchedule, nil)

	userScheduleQueryRepositoryMock.EXPECT().
		QueryUserScheduleWhereId(lastUpdatedIdOfUserSchedule).
		Return(makeOneUserScheduleDto(), nil) // Get the user after update

	// tag service mock
	tagServiceMock := testmock.NewMockTagServer(mockCtrl)
	tagServiceMock.EXPECT().
		GetTagsByTagTypeAndTagIds(tagservice.All, gomock.Any()).
		Return(makeRegularTags(), nil)

	//// Initialize server with mocks
	userScheduleServer := ProvideUserScheduleServer(
		userScheduleQueryRepositoryMock,
		userScheduleCommandRepositoryMock,
		tagServiceMock,
	)
	// Act
	uSchedules, err := userScheduleServer.UpdateUserSchedule(uid,
		&UserScheduleForCommand{
			FromDateTime: fromDateTime,
			ToDateTime:   toDateTime,
			TagIds:       tagIDs,
			Location:     location,
		},
	)
	// Assert
	if err != nil {
		t.Errorf("Test failed. Expected: no error', Actual: %s", err)
	}
	if uSchedules.UserId != uid {
		t.Errorf("Test failed. Expected: %s', Actual: %s", uid, uSchedules.UserId)
	}
	if len(uSchedules.UserSchedules) != 1 {
		t.Errorf("Test failed. Expected: 1', Actual: %d", len(uSchedules.UserSchedules))
	}
}

func TestUpdateUserSchedule_FromDateTimeIsAfterThanToDateTime_FromIsAfterToError(t *testing.T) {
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userScheduleServer := ProvideUserScheduleServer(
		NewMockIUserScheduleQueryRepository(mockCtrl),
		NewMockIUserScheduleCommandRepository(mockCtrl),
		testmock.NewMockTagServer(mockCtrl),
	)

	t.Run("From datetime is after than To datetime", func(t *testing.T) {
		// Arrange
		fromDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour+1, 30, 0, 0, time.Local) // After than toDateTime
		toDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour, 0, 0, 0, time.Local)
		// Act
		uSchedules, err := userScheduleServer.AddUserSchedule(uid,
			&UserScheduleForCommand{
				FromDateTime: fromDateTime,
				ToDateTime:   toDateTime,
				TagIds:       tagIDs,
				Location:     location,
			},
		)
		// Assert
		if _, ok := err.(*FromIsAfterToError); !ok {
			t.Errorf("Test failed. Expected: *FromIsAfterToError', Actual: %v", err)
		}
		if uSchedules != nil {
			t.Errorf("Test failed. Expected: nil', Actual: %v", uSchedules)
		}
	})
	t.Run("From datetime and To datetime are not in a same day", func(t *testing.T) {
		// Arrange
		fromDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour, 30, 0, 0, time.Local)
		toDateTime := time.Date(baseYear, baseMonth, baseDay+1, baseHour+1, 30, 0, 0, time.Local) // Different date from fromDateTime
		// Act
		uSchedules, err := userScheduleServer.AddUserSchedule(uid,
			&UserScheduleForCommand{
				FromDateTime: fromDateTime,
				ToDateTime:   toDateTime,
				TagIds:       tagIDs,
				Location:     location,
			},
		)
		// Assert
		if _, ok := err.(*DifferentDayFromAndToError); !ok {
			t.Errorf("Test failed. Expected: *FromIsAfterToError', Actual: %v", err)
		}
		if uSchedules != nil {
			t.Errorf("Test failed. Expected: nil', Actual: %v", uSchedules)
		}
	})
	t.Run("Time range From and To is less than business required speficication one", func(t *testing.T) {
		// Arrange
		fromDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour, 0, 0, 0, time.Local)
		toDateTime := time.Date(baseYear, baseMonth, baseDay, baseHour, 59, 0, 0, time.Local) // Within 60 minutes
		// Act
		uSchedules, err := userScheduleServer.AddUserSchedule(uid,
			&UserScheduleForCommand{
				FromDateTime: fromDateTime,
				ToDateTime:   toDateTime,
				TagIds:       tagIDs,
				Location:     location,
			},
		)
		// Assert
		if _, ok := err.(*TimeRangeIsLessThanSpecifiedError); !ok {
			t.Errorf("Test failed. Expected: *TimeRangeIsLessThanSpecifiedError', Actual: %v", err)
		}
		if uSchedules != nil {
			t.Errorf("Test failed. Expected: nil', Actual: %v", uSchedules)
		}
	})
}

func TestUpdateUserSchedule_ThereAreMoreThanOneSchedulesInTheDay_Error(t *testing.T) {
	// Arrange
	fromDateTime := time.Date(baseYear, baseMonth, baseDay+5, baseHour, 30, 0, 0, time.Local)
	toDateTime := time.Date(baseYear, baseMonth, baseDay+5, baseHour+1, 30, 0, 0, time.Local)
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userScheduleQueryRepositoryMock := NewMockIUserScheduleQueryRepository(mockCtrl)
	userScheduleQueryRepositoryMock.EXPECT().
		QueryUserSchedulesWhereTimeRange(gomock.Any(), gomock.Any(), uid).
		Return(makeSomeUserScheduleDtos(uid), nil) // [Error] There are more than one user schedules before updating
	//// Initialize server with mocks
	userScheduleServer := ProvideUserScheduleServer(
		userScheduleQueryRepositoryMock,
		NewMockIUserScheduleCommandRepository(mockCtrl),
		testmock.NewMockTagServer(mockCtrl),
	)
	// Act
	uSchedules, err := userScheduleServer.UpdateUserSchedule(uid,
		&UserScheduleForCommand{
			FromDateTime: fromDateTime,
			ToDateTime:   toDateTime,
			TagIds:       tagIDs,
			Location:     location,
		},
	)
	// Assert
	if err == nil {
		t.Errorf("Test failed. Expected: There is an error, Actual: nil")
	}
	if uSchedules != nil {
		t.Errorf("Test failed. Expected: nil', Actual: %v", uSchedules)
	}
}

func TestUpdateUserSchedule_NoScheduleInTheDay_TheScheduleNotFoundError(t *testing.T) {
	// Arrange
	fromDateTime := time.Date(baseYear, baseMonth, baseDay+5, baseHour, 30, 0, 0, time.Local)
	toDateTime := time.Date(baseYear, baseMonth, baseDay+5, baseHour+1, 30, 0, 0, time.Local)
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userScheduleQueryRepositoryMock := NewMockIUserScheduleQueryRepository(mockCtrl)
	userScheduleQueryRepositoryMock.EXPECT().
		QueryUserSchedulesWhereTimeRange(gomock.Any(), gomock.Any(), uid).
		Return([]*UserScheduleDto{}, nil) // [Error] There's no user schedule before updating
	//// Initialize server with mocks
	userScheduleServer := ProvideUserScheduleServer(
		userScheduleQueryRepositoryMock,
		NewMockIUserScheduleCommandRepository(mockCtrl),
		testmock.NewMockTagServer(mockCtrl),
	)
	// Act
	uSchedules, err := userScheduleServer.UpdateUserSchedule(uid,
		&UserScheduleForCommand{
			FromDateTime: fromDateTime,
			ToDateTime:   toDateTime,
			TagIds:       tagIDs,
			Location:     location,
		},
	)
	// Assert
	var e *TheScheduleNotFoundError
	if !errors.As(err, &e) {
		t.Errorf("Test failed. Expected: *TheScheduleNotFoundError', Actual: %v", err)
	}
	if uSchedules != nil {
		t.Errorf("Test failed. Expected: nil', Actual: %v", uSchedules)
	}
}

func TestDeleteUserSchedule_EverythingIsOk_NoError(t *testing.T) {
	// Arrange
	targetDateTime := time.Date(baseYear, baseMonth, baseDay+5, 0, 0, 0, 0, time.Local)
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	queriedIdOfUserSchedule := anyInt64

	userScheduleQueryRepositoryMock := NewMockIUserScheduleQueryRepository(mockCtrl)
	userScheduleQueryRepositoryMock.EXPECT().
		QueryUserSchedulesWhereTimeRange(gomock.Any(), gomock.Any(), uid).
		Return(makeOnlyOneUserScheduleDtos(uid, queriedIdOfUserSchedule), nil) // Only one user is expected before update

	userScheduleCommandRepositoryMock := NewMockIUserScheduleCommandRepository(mockCtrl)
	userScheduleCommandRepositoryMock.EXPECT().
		DeleteUserSchedule(queriedIdOfUserSchedule).
		Return(nil)

	// tag service mock
	tagServiceMock := testmock.NewMockTagServer(mockCtrl)
	tagServiceMock.EXPECT().
		GetTagsByTagTypeAndTagIds(tagservice.All, gomock.Any()).
		Return(makeRegularTags(), nil)

	//// Initialize server with mocks
	userScheduleServer := ProvideUserScheduleServer(
		userScheduleQueryRepositoryMock,
		userScheduleCommandRepositoryMock,
		tagServiceMock,
	)
	// Act
	uSchedules, err := userScheduleServer.DeleteUserSchedule(uid, targetDateTime)
	// Assert
	if err != nil {
		t.Errorf("Test failed. Expected: no error', Actual: %s", err)
	}
	if uSchedules.UserId != uid {
		t.Errorf("Test failed. Expected: %s', Actual: %s", uid, uSchedules.UserId)
	}
	if len(uSchedules.UserSchedules) != 1 {
		t.Errorf("Test failed. Expected: 1', Actual: %d", len(uSchedules.UserSchedules))
	}
}

func TestDeleteUserSchedule_NoScheduleInTheDay_TheScheduleNotFoundError(t *testing.T) {
	// Arrange
	targetDateTime := time.Date(baseYear, baseMonth, baseDay+5, 0, 0, 0, 0, time.Local)
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userScheduleQueryRepositoryMock := NewMockIUserScheduleQueryRepository(mockCtrl)
	userScheduleQueryRepositoryMock.EXPECT().
		QueryUserSchedulesWhereTimeRange(gomock.Any(), gomock.Any(), uid).
		Return([]*UserScheduleDto{}, nil) // [Error] There's no user schedule before deleting
	//// Initialize server with mocks
	userScheduleServer := ProvideUserScheduleServer(
		userScheduleQueryRepositoryMock,
		NewMockIUserScheduleCommandRepository(mockCtrl),
		testmock.NewMockTagServer(mockCtrl),
	)
	// Act
	uSchedules, err := userScheduleServer.DeleteUserSchedule(uid, targetDateTime)
	// Assert
	var e *TheScheduleNotFoundError
	if !errors.As(err, &e) {
		t.Errorf("Test failed. Expected: *TheScheduleNotFoundError', Actual: %v", err)
	}
	if uSchedules != nil {
		t.Errorf("Test failed. Expected: nil', Actual: %v", uSchedules)
	}
}
