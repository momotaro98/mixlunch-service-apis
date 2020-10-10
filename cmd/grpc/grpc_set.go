package main

import (
	"github.com/google/wire"

	"github.com/momotaro98/mixlunch-service-api/partyservice"
	"github.com/momotaro98/mixlunch-service-api/tagservice"
	usService "github.com/momotaro98/mixlunch-service-api/userscheduleservice"
	"github.com/momotaro98/mixlunch-service-api/userservice"
)

var gRPCSuperSet = wire.NewSet(
	// Tag service
	tagservice.ProvideDB,
	tagservice.ProvideTagQueryRepository,
	tagservice.ProvideTagServer,
	// User Schedule service, which uses Tag service
	usService.ProvideDB,
	usService.ProvideUserScheduleRepository,
	usService.ProvideRealUserScheduleUpdateRepository,
	usService.ProvideUserScheduleServer,
	// User service, which uses Tag service
	userservice.ProvideDB,
	userservice.ProvideUserQueryRepository,
	userservice.ProvideUserCommandRepository,
	userservice.ProvideUserServer,
	// Party service, which uses User service
	partyservice.ProvideDB,
	partyservice.ProvideApp,
	partyservice.ProvidePartyQueryRepository,
	partyservice.ProvidePartyCommandRepository,
	partyservice.ProvideChatRoomRepository,
	partyservice.ProvidePartyServer,
)
