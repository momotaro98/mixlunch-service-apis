package partyservice

import (
	"github.com/google/wire"

	"github.com/momotaro98/mixlunch-service-api/tagservice"
	"github.com/momotaro98/mixlunch-service-api/userservice"
)

var SuperSet = wire.NewSet(
	// Tag service
	tagservice.ProvideDB,
	tagservice.ProvideTagQueryRepository,
	tagservice.ProvideTagServer,
	// User service, which uses Tag service
	userservice.ProvideDB,
	userservice.ProvideUserQueryRepository,
	userservice.ProvideUserCommandRepository,
	userservice.ProvideUserServer,
	// Party service, which uses User service
	ProvideDB,
	ProvideApp,
	ProvidePartyQueryRepository,
	ProvidePartyCommandRepository,
	ProvideChatRoomRepository,
	ProvidePartyServer,
)
