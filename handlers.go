package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/momotaro98/stew"

	"github.com/momotaro98/mixlunch-service-api/domainerror"
	"github.com/momotaro98/mixlunch-service-api/logger"
	"github.com/momotaro98/mixlunch-service-api/partyservice"
	"github.com/momotaro98/mixlunch-service-api/tagservice"
	usService "github.com/momotaro98/mixlunch-service-api/userscheduleservice"
	"github.com/momotaro98/mixlunch-service-api/userservice"
)

var (
	XRequestId = "x-request-id"
)

type ErrorResponse struct {
	Message string                `json:"message"`
	Code    domainerror.ErrorCode `json:"code"`
}

func NewErrorResponse(message string, code domainerror.ErrorCode) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Code:    code,
	}
}

func assembleErrorResponse(domainError domainerror.DomainError) []byte {
	message := domainError.Error()
	code := domainError.Code()
	errorResponse := NewErrorResponse(message, code)
	res, err := json.Marshal(errorResponse)
	if err != nil {
		panic(err)
	}
	return res
}

func responseWithSuccess(l logger.Logger, reqId string, modelToMarshalToJSON interface{}, w http.ResponseWriter) {
	// Success case : JSON Marshal and return response to client
	res, err := json.Marshal(modelToMarshalToJSON)
	if err != nil {
		l.Log(logger.Error, reqId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		l.Log(logger.Error, reqId, err.Error())
		panic(err)
	}
}

func handleError(w http.ResponseWriter, r *http.Request, l logger.Logger, err error) {
	var (
		reqId = r.Header.Get(XRequestId)
	)

	var domainErr domainerror.DomainError

	if errors.As(err, &domainErr) {
		l.Log(logger.Warn, reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write(assembleErrorResponse(domainErr)); err != nil {
			panic(err)
		}
		return
	} else if err != nil {
		l.Log(logger.Error, reqId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func httpGetWrap(w http.ResponseWriter, r *http.Request, l logger.Logger, f func() (interface{}, error)) {
	var (
		reqId = r.Header.Get(XRequestId)
	)
	l.Log(logger.Info, reqId, fmt.Sprintf("Got GET request. URL: %s", r.URL.Path))

	retFromService, err := f()
	if err != nil {
		handleError(w, r, l, err)
		return
	}

	responseWithSuccess(l, reqId, retFromService, w)
}

func httpPostWrap(w http.ResponseWriter, r *http.Request, l logger.Logger, decoding interface{}, f func(decoded interface{}) (interface{}, error)) {
	var (
		reqId = r.Header.Get(XRequestId)
	)
	l.Log(logger.Info, reqId, fmt.Sprintf("Got POST request. URL: %s", r.URL.Path))

	// Parse the request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(decoding)
	if err != nil {
		l.Log(logger.Error, reqId, err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		jsonErr := domainerror.NewJSONParseError(r.URL.Path, err)
		if _, err = w.Write(assembleErrorResponse(jsonErr)); err != nil {
			panic(err)
		}
		return
	}

	retFromService, err := f(decoding)
	if err != nil {
		handleError(w, r, l, err)
		return
	}

	responseWithSuccess(l, reqId, retFromService, w)
}

type UserScheduleHandler struct {
	logger logger.Logger
	server usService.UserScheduleServer
}

func provideUserScheduleHandler(logger logger.Logger, server usService.UserScheduleServer) *UserScheduleHandler {
	return &UserScheduleHandler{
		logger: logger,
		server: server,
	}
}

func (h *UserScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		params        = mux.Vars(r)
		uid           = params["uid"]
		beginDateTime = params["beginDateTime"]
		endDateTime   = params["endDateTime"]
	)
	httpGetWrap(w, r, h.logger, func() (interface{}, error) {
		ret, err := h.server.GetUserSchedulesByTimeRange(uid, beginDateTime, endDateTime)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return ret, nil
	})
}

type AddUserScheduleHandler struct {
	logger logger.Logger
	server usService.UserScheduleServer
}

func provideAddUserScheduleHandler(logger logger.Logger, server usService.UserScheduleServer) *AddUserScheduleHandler {
	return &AddUserScheduleHandler{
		logger: logger,
		server: server,
	}
}

func (h *AddUserScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		params = mux.Vars(r)
		uid    = params["uid"]
	)
	var addingUserSchedule usService.UserScheduleForCommand
	httpPostWrap(w, r, h.logger, &addingUserSchedule, func(decoded interface{}) (interface{}, error) {
		us, _ := decoded.(*usService.UserScheduleForCommand)
		ret, err := h.server.AddUserSchedule(uid, us)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return ret, nil
	})
}

type UpdateUserScheduleHandler struct {
	logger logger.Logger
	server usService.UserScheduleServer
}

func provideUpdateUserScheduleHandler(logger logger.Logger, server usService.UserScheduleServer) *UpdateUserScheduleHandler {
	return &UpdateUserScheduleHandler{
		logger: logger,
		server: server,
	}
}

func (h *UpdateUserScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		params = mux.Vars(r)
		uid    = params["uid"]
	)
	var updatingUserSchedule usService.UserScheduleForCommand
	httpPostWrap(w, r, h.logger, &updatingUserSchedule, func(decoded interface{}) (interface{}, error) {
		us, _ := decoded.(*usService.UserScheduleForCommand)
		ret, err := h.server.UpdateUserSchedule(uid, us)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return ret, nil
	})
}

type DeleteUserScheduleHandler struct {
	logger logger.Logger
	server usService.UserScheduleServer
}

func provideDeleteUserScheduleHandler(logger logger.Logger, server usService.UserScheduleServer) *DeleteUserScheduleHandler {
	return &DeleteUserScheduleHandler{
		logger: logger,
		server: server,
	}
}

func (h *DeleteUserScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		params = mux.Vars(r)
		uid    = params["uid"]
	)
	var deletingDate usService.SpecifiedDate
	httpPostWrap(w, r, h.logger, &deletingDate, func(decoded interface{}) (interface{}, error) {
		sd, _ := decoded.(*usService.SpecifiedDate)
		ret, err := h.server.DeleteUserSchedule(uid, sd.Date)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return ret, nil
	})
}

type PartyHandler struct {
	logger logger.Logger
	server partyservice.PartyServer
}

func providePartyHandler(logger logger.Logger, server partyservice.PartyServer) *PartyHandler {
	return &PartyHandler{
		logger: logger,
		server: server,
	}
}

func (h *PartyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		params        = mux.Vars(r)
		uid           = params["uid"]
		beginDateTime = params["beginDateTime"]
		endDateTime   = params["endDateTime"]
	)
	httpGetWrap(w, r, h.logger, func() (interface{}, error) {
		ret, err := h.server.GetPartyByUserIdAndTimeRange(uid, beginDateTime, endDateTime)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return ret, nil
	})
}

type PartyReviewMemberHandler struct {
	logger logger.Logger
	server partyservice.PartyServer
}

func providePartyReviewMemberHandler(logger logger.Logger, server partyservice.PartyServer) *PartyReviewMemberHandler {
	return &PartyReviewMemberHandler{
		logger: logger,
		server: server,
	}
}

func (h *PartyReviewMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var newReviewMember partyservice.PartyReviewMember
	httpPostWrap(w, r, h.logger, &newReviewMember, func(decoded interface{}) (interface{}, error) {
		reviewMember, _ := decoded.(*partyservice.PartyReviewMember)
		err := h.server.PostPartyReviewMember(reviewMember)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return nil, nil
	})
}

type PartyReviewMemberDoneHandler struct {
	logger logger.Logger
	server partyservice.PartyServer
}

func providePartyReviewMemberDoneHandler(logger logger.Logger, server partyservice.PartyServer) *PartyReviewMemberDoneHandler {
	return &PartyReviewMemberDoneHandler{
		logger: logger,
		server: server,
	}
}

func (h *PartyReviewMemberDoneHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		params   = mux.Vars(r)
		reviewer = params["reviewer"]
	)
	httpGetWrap(w, r, h.logger, func() (interface{}, error) {
		ret, err := h.server.GetIsLatestPartyReviewDone(reviewer)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return ret, nil
	})
}

type TagsHandler struct {
	logger logger.Logger
	server tagservice.TagServer
}

func provideTagsHandler(logger logger.Logger, server tagservice.TagServer) *TagsHandler {
	return &TagsHandler{
		logger: logger,
		server: server,
	}
}

func (h *TagsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		params = mux.Vars(r)
		ttid   = params["ttid"]
	)
	var tagType tagservice.TagType
	if ttid == "" {
		tagType = tagservice.All
	} else {
		tagTypeId, _ := strconv.Atoi(ttid)
		tagType = tagservice.TagType(tagTypeId)
	}

	httpGetWrap(w, r, h.logger, func() (interface{}, error) {
		ret, err := h.server.GetTagsByTagType(tagType)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return ret, nil
	})
}

type UserHandler struct {
	logger logger.Logger
	server userservice.UserServer
}

func provideUserHandler(logger logger.Logger, server userservice.UserServer) *UserHandler {
	return &UserHandler{
		logger: logger,
		server: server,
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		params = mux.Vars(r)
		uid    = params["uid"]
	)
	httpGetWrap(w, r, h.logger, func() (interface{}, error) {
		ret, err := h.server.GetUserByUserId(uid)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return ret, nil
	})
}

type UserPublicHandler struct {
	logger logger.Logger
	server userservice.UserServer
}

func provideUserPublicHandler(logger logger.Logger, server userservice.UserServer) *UserPublicHandler {
	return &UserPublicHandler{
		logger: logger,
		server: server,
	}
}

func (h *UserPublicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		params = mux.Vars(r)
		uid    = params["uid"]
	)
	httpGetWrap(w, r, h.logger, func() (interface{}, error) {
		ret, err := h.server.GetUserPublicByUserId(uid)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return ret, nil
	})
}

type UserRegisterHandler struct {
	logger logger.Logger
	server userservice.UserServer
}

func provideUserRegisterHandler(logger logger.Logger, server userservice.UserServer) *UserRegisterHandler {
	return &UserRegisterHandler{
		logger: logger,
		server: server,
	}
}

func (h *UserRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var newUser userservice.UserForCommand
	httpPostWrap(w, r, h.logger, &newUser, func(decoded interface{}) (interface{}, error) {
		user, _ := decoded.(*userservice.UserForCommand)
		ret, err := h.server.RegisterUser(user)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return ret, nil
	})
}

type UserBlockRegisterHandler struct {
	logger logger.Logger
	server userservice.UserServer
}

func provideUserBlockRegisterHandler(logger logger.Logger, server userservice.UserServer) *UserBlockRegisterHandler {
	return &UserBlockRegisterHandler{
		logger: logger,
		server: server,
	}
}

func (h *UserBlockRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var newUserBlock userservice.UserBlockForCommand
	httpPostWrap(w, r, h.logger, &newUserBlock, func(decoded interface{}) (interface{}, error) {
		userBlock, _ := decoded.(*userservice.UserBlockForCommand)
		ret, err := h.server.RegisterUserBlock(userBlock)
		if err != nil {
			return nil, stew.Wrap(err)
		}
		return ret, nil
	})
}
