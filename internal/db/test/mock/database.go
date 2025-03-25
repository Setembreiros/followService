// Code generated by MockGen. DO NOT EDIT.
// Source: database.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	model "followservice/internal/model/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDatabaseClient is a mock of DatabaseClient interface.
type MockDatabaseClient struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseClientMockRecorder
}

// MockDatabaseClientMockRecorder is the mock recorder for MockDatabaseClient.
type MockDatabaseClientMockRecorder struct {
	mock *MockDatabaseClient
}

// NewMockDatabaseClient creates a new mock instance.
func NewMockDatabaseClient(ctrl *gomock.Controller) *MockDatabaseClient {
	mock := &MockDatabaseClient{ctrl: ctrl}
	mock.recorder = &MockDatabaseClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabaseClient) EXPECT() *MockDatabaseClientMockRecorder {
	return m.recorder
}

// Clean mocks base method.
func (m *MockDatabaseClient) Clean() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Clean")
}

// Clean indicates an expected call of Clean.
func (mr *MockDatabaseClientMockRecorder) Clean() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clean", reflect.TypeOf((*MockDatabaseClient)(nil).Clean))
}

// CreateRelationship mocks base method.
func (m *MockDatabaseClient) CreateRelationship(relationship *model.UserPairRelationship) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRelationship", relationship)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRelationship indicates an expected call of CreateRelationship.
func (mr *MockDatabaseClientMockRecorder) CreateRelationship(relationship interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRelationship", reflect.TypeOf((*MockDatabaseClient)(nil).CreateRelationship), relationship)
}

// DeleteRelationship mocks base method.
func (m *MockDatabaseClient) DeleteRelationship(relationship *model.UserPairRelationship) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRelationship", relationship)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRelationship indicates an expected call of DeleteRelationship.
func (mr *MockDatabaseClientMockRecorder) DeleteRelationship(relationship interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRelationship", reflect.TypeOf((*MockDatabaseClient)(nil).DeleteRelationship), relationship)
}

// GetUserFollowers mocks base method.
func (m *MockDatabaseClient) GetUserFollowers(username, lastPostId string, limit int) ([]string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFollowers", username, lastPostId, limit)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUserFollowers indicates an expected call of GetUserFollowers.
func (mr *MockDatabaseClientMockRecorder) GetUserFollowers(username, lastPostId, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFollowers", reflect.TypeOf((*MockDatabaseClient)(nil).GetUserFollowers), username, lastPostId, limit)
}

// RelationshipExists mocks base method.
func (m *MockDatabaseClient) RelationshipExists(userPair *model.UserPairRelationship) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RelationshipExists", userPair)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RelationshipExists indicates an expected call of RelationshipExists.
func (mr *MockDatabaseClientMockRecorder) RelationshipExists(userPair interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RelationshipExists", reflect.TypeOf((*MockDatabaseClient)(nil).RelationshipExists), userPair)
}
