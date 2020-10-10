// Code generated by MockGen. DO NOT EDIT.
// Source: repositories.go

// Package userscheduleservice is a generated GoMock package.
package userscheduleservice

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockIUserScheduleQueryRepository is a mock of IUserScheduleQueryRepository interface
type MockIUserScheduleQueryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIUserScheduleQueryRepositoryMockRecorder
}

// MockIUserScheduleQueryRepositoryMockRecorder is the mock recorder for MockIUserScheduleQueryRepository
type MockIUserScheduleQueryRepositoryMockRecorder struct {
	mock *MockIUserScheduleQueryRepository
}

// NewMockIUserScheduleQueryRepository creates a new mock instance
func NewMockIUserScheduleQueryRepository(ctrl *gomock.Controller) *MockIUserScheduleQueryRepository {
	mock := &MockIUserScheduleQueryRepository{ctrl: ctrl}
	mock.recorder = &MockIUserScheduleQueryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIUserScheduleQueryRepository) EXPECT() *MockIUserScheduleQueryRepositoryMockRecorder {
	return m.recorder
}

// QueryUserSchedulesWhereTimeRange mocks base method
func (m *MockIUserScheduleQueryRepository) QueryUserSchedulesWhereTimeRange(beginDateTime, endDateTime time.Time, userId string) ([]*UserScheduleDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryUserSchedulesWhereTimeRange", beginDateTime, endDateTime, userId)
	ret0, _ := ret[0].([]*UserScheduleDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryUserSchedulesWhereTimeRange indicates an expected call of QueryUserSchedulesWhereTimeRange
func (mr *MockIUserScheduleQueryRepositoryMockRecorder) QueryUserSchedulesWhereTimeRange(beginDateTime, endDateTime, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryUserSchedulesWhereTimeRange", reflect.TypeOf((*MockIUserScheduleQueryRepository)(nil).QueryUserSchedulesWhereTimeRange), beginDateTime, endDateTime, userId)
}

// QueryUserScheduleWhereId mocks base method
func (m *MockIUserScheduleQueryRepository) QueryUserScheduleWhereId(userScheduleId int64) (*UserScheduleDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryUserScheduleWhereId", userScheduleId)
	ret0, _ := ret[0].(*UserScheduleDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryUserScheduleWhereId indicates an expected call of QueryUserScheduleWhereId
func (mr *MockIUserScheduleQueryRepositoryMockRecorder) QueryUserScheduleWhereId(userScheduleId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryUserScheduleWhereId", reflect.TypeOf((*MockIUserScheduleQueryRepository)(nil).QueryUserScheduleWhereId), userScheduleId)
}

// MockIUserScheduleCommandRepository is a mock of IUserScheduleCommandRepository interface
type MockIUserScheduleCommandRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIUserScheduleCommandRepositoryMockRecorder
}

// MockIUserScheduleCommandRepositoryMockRecorder is the mock recorder for MockIUserScheduleCommandRepository
type MockIUserScheduleCommandRepositoryMockRecorder struct {
	mock *MockIUserScheduleCommandRepository
}

// NewMockIUserScheduleCommandRepository creates a new mock instance
func NewMockIUserScheduleCommandRepository(ctrl *gomock.Controller) *MockIUserScheduleCommandRepository {
	mock := &MockIUserScheduleCommandRepository{ctrl: ctrl}
	mock.recorder = &MockIUserScheduleCommandRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIUserScheduleCommandRepository) EXPECT() *MockIUserScheduleCommandRepositoryMockRecorder {
	return m.recorder
}

// InsertUserSchedule mocks base method
func (m *MockIUserScheduleCommandRepository) InsertUserSchedule(dto *UserScheduleDto) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUserSchedule", dto)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertUserSchedule indicates an expected call of InsertUserSchedule
func (mr *MockIUserScheduleCommandRepositoryMockRecorder) InsertUserSchedule(dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUserSchedule", reflect.TypeOf((*MockIUserScheduleCommandRepository)(nil).InsertUserSchedule), dto)
}

// UpdateUserSchedule mocks base method
func (m *MockIUserScheduleCommandRepository) UpdateUserSchedule(userScheduleId int64, dto *UserScheduleDto) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserSchedule", userScheduleId, dto)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserSchedule indicates an expected call of UpdateUserSchedule
func (mr *MockIUserScheduleCommandRepositoryMockRecorder) UpdateUserSchedule(userScheduleId, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserSchedule", reflect.TypeOf((*MockIUserScheduleCommandRepository)(nil).UpdateUserSchedule), userScheduleId, dto)
}

// DeleteUserSchedule mocks base method
func (m *MockIUserScheduleCommandRepository) DeleteUserSchedule(userScheduleId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserSchedule", userScheduleId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserSchedule indicates an expected call of DeleteUserSchedule
func (mr *MockIUserScheduleCommandRepositoryMockRecorder) DeleteUserSchedule(userScheduleId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserSchedule", reflect.TypeOf((*MockIUserScheduleCommandRepository)(nil).DeleteUserSchedule), userScheduleId)
}