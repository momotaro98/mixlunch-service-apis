package main

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/momotaro98/mixlunch-service-api/cmd/grpc/testmock"
	"github.com/momotaro98/mixlunch-service-api/logger"
	"github.com/momotaro98/mixlunch-service-api/partyservice"
	"github.com/momotaro98/mixlunch-service-api/userservice"
)

//go:generate mockgen -source=../../userscheduleservice/domain.go -destination=testmock/userscheduleservice.go -package=testmock
//go:generate mockgen -source=../../partyservice/domain.go -destination=testmock/partyservice.go -package=testmock
//go:generate mockgen -source=../../userservice/domain.go -destination=testmock/userservice.go -package=testmock

func TestAssembleBlacklist(t *testing.T) {
	const (
		userID = "user-id-test"
		mateA  = "user-id-A"
		mateB  = "user-id-B"
		mateC  = "user-id-C"
		mateD  = "user-id-D"
		mateE  = "user-id-E"
	)
	// mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Input
	partyMock := testmock.NewMockPartyServer(mockCtrl)
	partyMock.EXPECT().
		GetLastNPartiesOfAUser(userID, 3).
		Return(&partyservice.Parties{
			Parties: []*partyservice.Party{
				{
					Members: []*userservice.UserPublic{
						{
							UserId: userID,
						},
						{
							UserId: mateA,
						},
					},
				},
				{
					Members: []*userservice.UserPublic{
						{
							UserId: userID,
						},
						{
							UserId: mateB,
						},
					},
				},
				{
					Members: []*userservice.UserPublic{
						{
							UserId: userID,
						},
						{
							UserId: mateC,
						},
					},
				},
			},
		}, nil)

	var (
		begin = time.Now().AddDate(0, 0, -14)
	)
	partyMock.EXPECT().
		GetPartyByUserIdAndTimeRange(userID, begin.Format(time.RFC3339), "").
		Return(&partyservice.Parties{
			Parties: []*partyservice.Party{
				{
					Members: []*userservice.UserPublic{
						{
							UserId: userID,
						},
						{
							UserId: mateD,
						},
					},
				},
			},
		}, nil)

	grpcServer := provideGRPCMixLunchServer(
		logger.ProvideLogger(&logger.Config{ErrorLevel: "debug"}),
		testmock.NewMockUserScheduleServer(mockCtrl),
		partyMock,
		testmock.NewMockUserServer(mockCtrl),
	)

	// Input
	user := &userservice.User{
		UserId:        userID,
		BlockingUsers: []string{mateE},
	}

	// Act
	blackList, _ := grpcServer.assembleBlacklist(user)

	// Assert
	if len(blackList) != 5 {
		t.Errorf("expected: 5, got: %d", len(blackList))
	}
	sort.Strings(blackList)
	if !reflect.DeepEqual([]string{mateA, mateB, mateC, mateD, mateE}, blackList) {
		t.Errorf("expected: 5 mates, got: %+v", blackList)
	}
}
