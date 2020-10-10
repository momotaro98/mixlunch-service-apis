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

func initializeUserScheduleHandler(loggerConfig *logger.Config, userScheduleServiceConfig *usService.Config, tagServiceConfig *tagservice.Config) *UserScheduleHandler {
	wire.Build(logger.SuperSet, usService.SuperSet, provideUserScheduleHandler)
	return nil
}

func initializeUpdateUserScheduleHandler(loggerConfig *logger.Config, userScheduleServiceConfig *usService.Config, tagServiceConfig *tagservice.Config) *UpdateUserScheduleHandler {
	wire.Build(logger.SuperSet, usService.SuperSet, provideUpdateUserScheduleHandler)
	return nil
}

func initializeDeleteUserScheduleHandler(loggerConfig *logger.Config, userScheduleServiceConfig *usService.Config, tagServiceConfig *tagservice.Config) *DeleteUserScheduleHandler {
	wire.Build(logger.SuperSet, usService.SuperSet, provideDeleteUserScheduleHandler)
	return nil
}

func initializeAddUserScheduleHandler(loggerConfig *logger.Config, userScheduleServiceConfig *usService.Config, tagServiceConfig *tagservice.Config) *AddUserScheduleHandler {
	wire.Build(logger.SuperSet, usService.SuperSet, provideAddUserScheduleHandler)
	return nil
}

func initializePartyHandler(loggerConfig *logger.Config, partyServiceConfig *partyservice.Config, userServiceConfig *userservice.Config, tagServiceConfig *tagservice.Config) *PartyHandler {
	wire.Build(logger.SuperSet, partyservice.SuperSet, providePartyHandler)
	return nil
}

func initializePartyReviewMemberHandler(loggerConfig *logger.Config, partyServiceConfig *partyservice.Config, userServiceConfig *userservice.Config, tagServiceConfig *tagservice.Config) *PartyReviewMemberHandler {
	wire.Build(logger.SuperSet, partyservice.SuperSet, providePartyReviewMemberHandler)
	return nil
}

func initializePartyReviewMemberDoneHandler(loggerConfig *logger.Config, partyServiceConfig *partyservice.Config, userServiceConfig *userservice.Config, tagServiceConfig *tagservice.Config) *PartyReviewMemberDoneHandler {
	wire.Build(logger.SuperSet, partyservice.SuperSet, providePartyReviewMemberDoneHandler)
	return nil
}

func initializeTagsHandler(loggerConfig *logger.Config, tagServiceConfig *tagservice.Config) *TagsHandler {
	wire.Build(logger.SuperSet, tagservice.SuperSet, provideTagsHandler)
	return nil
}

func initializeUserHandler(loggerConfig *logger.Config, userServiceConfig *userservice.Config, tagServiceConfig *tagservice.Config) *UserHandler {
	wire.Build(logger.SuperSet, userservice.SuperSet, provideUserHandler)
	return nil
}

func initializeUserPublicHandler(loggerConfig *logger.Config, userServiceConfig *userservice.Config, tagServiceConfig *tagservice.Config) *UserPublicHandler {
	wire.Build(logger.SuperSet, userservice.SuperSet, provideUserPublicHandler)
	return nil
}

func initializeUserRegisterHandler(loggerConfig *logger.Config, userServiceConfig *userservice.Config, tagServiceConfig *tagservice.Config) *UserRegisterHandler {
	wire.Build(logger.SuperSet, userservice.SuperSet, provideUserRegisterHandler)
	return nil
}

func initializeUserBlockRegisterHandler(loggerConfig *logger.Config, userServiceConfig *userservice.Config, tagServiceConfig *tagservice.Config) *UserBlockRegisterHandler {
	wire.Build(logger.SuperSet, userservice.SuperSet, provideUserBlockRegisterHandler)
	return nil
}
