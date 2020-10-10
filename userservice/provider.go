package userservice

import (
	"github.com/google/wire"

	"github.com/momotaro98/mixlunch-service-api/tagservice"
)

var SuperSet = wire.NewSet(
	// Tag service
	tagservice.ProvideDB,
	tagservice.ProvideTagQueryRepository,
	tagservice.ProvideTagServer,
	// User service, which uses Tag service
	ProvideDB,
	ProvideUserQueryRepository,
	ProvideUserCommandRepository,
	ProvideUserServer,
)
