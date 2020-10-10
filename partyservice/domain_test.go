package partyservice

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/momotaro98/mixlunch-service-api/partyservice/testmock"
	"github.com/momotaro98/mixlunch-service-api/userservice"
	"github.com/momotaro98/mixlunch-service-api/utils"
)

//go:generate mockgen -source=repositories.go -destination=repositories_mock.go -package=partyservice -self_package=github.com/momotaro98/mixlunch-service-api/partyservice
//go:generate mockgen -source=../userservice/domain.go -destination=userservice_mock.go -package=partyservice -self_package=github.com/momotaro98/mixlunch-service-api/partyservice
//go:generate mockgen -source=../tagservice/domain.go -destination=testmock/tagservice.go -package=testmock

const (
	uid       = "userId0123456789"
	baseYear  = 2018
	baseMonth = 11
	baseDay   = 15
	baseHour  = 12
	anyInt64  = int64(100001)
)

func TestGetPartiesByTimeRange_InvalidDatetimeString_ThrowSpecificError(t *testing.T) {
	// mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	partyServer := ProvidePartyServer(
		NewMockIPartyQueryRepository(mockCtrl),
		NewMockIPartyCommandRepository(mockCtrl),
		NewMockUserServer(mockCtrl),
		testmock.NewMockTagServer(mockCtrl),
		NewMockIChatRoomRepository(mockCtrl))
	// How to act and assert
	test := func(t *testing.T, begin, end string) {
		// Act
		parties, err := partyServer.GetParties(begin, end)
		// Assert
		var e *InvalidDateTimeFormatError
		if !errors.As(err, &e) {
			t.Errorf("Test failed. Expected: *InvalidDateTimeFormatError', Actual: %v", err)
		}
		if parties != nil {
			t.Errorf("Test failed. Expected: nil', Actual: %v", parties)
		}
	}

	t.Run("Begin DateTime String is Invalid", func(t *testing.T) {
		// Arrange
		beginDateTimeStr := "2018-10-10 10:10:10" // Invalid
		endDateTimeStr := utils.MakeCorrectFormatDateTimeStr(baseYear, baseMonth, baseDay, baseHour, 0)
		// Act and Assert
		test(t, beginDateTimeStr, endDateTimeStr)
	})
	t.Run("End DateTime String is Invalid", func(t *testing.T) {
		// Arrange
		beginDateTimeStr := utils.MakeCorrectFormatDateTimeStr(baseYear, baseMonth, baseDay, baseHour, 0)
		endDateTimeStr := "2018-10-30 10:10:10" // Invalid
		// Act and Assert
		test(t, beginDateTimeStr, endDateTimeStr)
	})
}

func TestGetPartyByUserIdAndTimeRange_InvalidBeginDatetimeString_ThrowSpecificError(t *testing.T) {
	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	partyServer := ProvidePartyServer(
		NewMockIPartyQueryRepository(mockCtrl),
		NewMockIPartyCommandRepository(mockCtrl),
		NewMockUserServer(mockCtrl),
		testmock.NewMockTagServer(mockCtrl),
		NewMockIChatRoomRepository(mockCtrl))
	/// How to act and assert
	test := func(t *testing.T, userId, begin, end string) {
		// Act
		parties, err := partyServer.GetPartyByUserIdAndTimeRange(userId, begin, end)
		// Assert
		var e *InvalidDateTimeFormatError
		if !errors.As(err, &e) {
			t.Errorf("Test failed. Expected: *InvalidDateTimeFormatError', Actual: %v", err)
		}
		if parties != nil {
			t.Errorf("Test failed. Expected: nil', Actual: %v", parties)
		}
	}

	t.Run("Begin DateTime string is invalid", func(t *testing.T) {
		/// Arrange
		beginDateTimeStr := "2018-10-10 10:10:10" // Invalid
		endDateTimeStr := utils.MakeCorrectFormatDateTimeStr(baseYear, baseMonth, baseDay, baseHour, 0)
		// Act and Assert
		test(t, uid, beginDateTimeStr, endDateTimeStr)
	})
	t.Run("End DateTime string is invalid", func(t *testing.T) {
		/// Arrange
		beginDateTimeStr := utils.MakeCorrectFormatDateTimeStr(baseYear, baseMonth, baseDay, baseHour, 0)
		endDateTimeStr := "2018-10-30 10:10:10" // Invalid
		// Act and Assert
		test(t, uid, beginDateTimeStr, endDateTimeStr)
	})
}

func TestGetLastNPartiesOfAUser(t *testing.T) {
	const (
		userID = "user-id"
	)

	var partyDtos = []*PartyDto{
		{
			id:        1,
			startFrom: time.Now(),
			endTo:     time.Now(),
		},
		{
			id:        2,
			startFrom: time.Now(),
			endTo:     time.Now(),
		},
		{
			id:        3,
			startFrom: time.Now(),
			endTo:     time.Now(),
		},
	}

	// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	partyQueryRepoMock := NewMockIPartyQueryRepository(mockCtrl)

	partyQueryRepoMock.EXPECT().
		QueryPartyMembersWherePartyId(gomock.Any()).
		Return([]*PartyMemberDto{
			{
				partyId: 1,
				userId:  userID,
			},
			{
				partyId: 1,
				userId:  "lunch-mate",
			},
		}, nil).AnyTimes()

	userServerMock := NewMockUserServer(mockCtrl)
	userServerMock.EXPECT().GetUserPublicByUserId(gomock.Any()).
		Return(&userservice.UserPublic{
			UserId: "lunch-mate",
		}, nil).AnyTimes()

	partyQueryRepoMock.EXPECT().QueryPartyTagsWherePartyId(gomock.Any()).
		Return(&PartyTagsDto{
			tagIds: []uint16{24, 56, 60}, // random
		}, nil).AnyTimes()

	tagServerMock := testmock.NewMockTagServer(mockCtrl)
	tagServerMock.EXPECT().GetTagsByTagTypeAndTagIds(gomock.Any(), gomock.Any()).
		Return(nil, nil).AnyTimes()

	assert := func(t *testing.T, partyServer PartyServer, n, expected int) {
		// Act
		parties, err := partyServer.GetLastNPartiesOfAUser(userID, n)
		if err != nil {
			t.Errorf("expected: nil, got: %+v", err)
		}
		if len(parties.Parties) != expected {
			t.Errorf("expected: %d, got: %d", expected, len(parties.Parties))
		}
	}

	t.Run("n is equal to parties", func(t *testing.T) {
		n := 3
		partyQueryRepoMock.EXPECT().
			QueryPartiesWhereUserIdLastN(userID, n).
			Return(partyDtos, nil)
		partyServer := ProvidePartyServer(
			partyQueryRepoMock,
			NewMockIPartyCommandRepository(mockCtrl),
			userServerMock,
			tagServerMock,
			NewMockIChatRoomRepository(mockCtrl),
		)
		assert(t, partyServer, n, 3)
	})

	t.Run("n is more than parties", func(t *testing.T) {
		n := 4
		partyQueryRepoMock.EXPECT().
			QueryPartiesWhereUserIdLastN(userID, n).
			Return(partyDtos, nil)
		partyServer := ProvidePartyServer(
			partyQueryRepoMock,
			NewMockIPartyCommandRepository(mockCtrl),
			userServerMock,
			tagServerMock,
			NewMockIChatRoomRepository(mockCtrl),
		)
		assert(t, partyServer, n, 3)
	})

	t.Run("parties is empty then return empty", func(t *testing.T) {
		n := 3
		partyDtos = []*PartyDto{}
		partyQueryRepoMock.EXPECT().
			QueryPartiesWhereUserIdLastN(userID, n).
			Return(partyDtos, nil)
		partyServer := ProvidePartyServer(
			partyQueryRepoMock,
			NewMockIPartyCommandRepository(mockCtrl),
			userServerMock,
			tagServerMock,
			NewMockIChatRoomRepository(mockCtrl),
		)
		assert(t, partyServer, n, 0)
	})
}

func TestUpsertParties_Successfully_Nil(t *testing.T) {
	// Arrange
	/// Business
	userId1 := "user-id-1"
	userId2 := "user-id-2"
	startFrom, endTo := time.Now(), time.Now()
	chatRoomId := ""
	members1 := []*userservice.UserPublic{
		{UserId: userId1},
		{UserId: userId2},
	}
	partyModel1 := NewPartyForCommand(startFrom, endTo, chatRoomId, members1)
	userId3 := "user-id-3"
	userId4 := "user-id-4"
	startFrom2, endTo2 := time.Now(), time.Now()
	chatRoomId2 := ""
	members2 := []*userservice.UserPublic{
		{UserId: userId3},
		{UserId: userId4},
	}
	partyModel2 := NewPartyForCommand(startFrom2, endTo2, chatRoomId2, members2)
	parties := []*PartyForCommand{partyModel1, partyModel2}

	/// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// Query
	partyQueryRepository := NewMockIPartyQueryRepository(mockCtrl)
	// Command
	partyCommandRepository := NewMockIPartyCommandRepository(mockCtrl)
	partyCommandRepository.EXPECT().Tran().Return(&sql.Tx{}, nil)
	partyCommandRepository.EXPECT().Commit(gomock.Any()).Return(nil)
	//partyCommandRepository.EXPECT().Rollback(gomock.Any()).Return(nil)
	partyCommandRepository.EXPECT().InsertParty(gomock.Any(), gomock.Any()).Return(anyInt64, nil).AnyTimes()
	partyCommandRepository.EXPECT().DeletePartiesWithADay(gomock.Any(), gomock.Any()).Return(nil)
	partyServer := ProvidePartyServer(
		partyQueryRepository,
		partyCommandRepository,
		NewMockUserServer(mockCtrl),
		testmock.NewMockTagServer(mockCtrl),
		NewMockIChatRoomRepository(mockCtrl))

	// Act
	err := partyServer.UpsertParties(parties)
	// Assert
	if err != nil {
		t.Errorf("Test failed. Expected: nil', Actual: %v", err)
	}
}
