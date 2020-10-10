// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"

	"github.com/momotaro98/mixlunch-service-api/logger"
	"github.com/momotaro98/mixlunch-service-api/partyservice"
	"github.com/momotaro98/mixlunch-service-api/tagservice"
	usService "github.com/momotaro98/mixlunch-service-api/userscheduleservice"
	"github.com/momotaro98/mixlunch-service-api/userservice"
)

func initializeGRPCServer(loggerConfig *logger.Config, usServiceDbConfig *usService.Config, partyServiceConfig *partyservice.Config, userServiceConfig *userservice.Config, tagServiceConfig *tagservice.Config) *gRPCMixLunchServer {
	wire.Build(logger.SuperSet, gRPCSuperSet, provideGRPCMixLunchServer)
	return nil
}
