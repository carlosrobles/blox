// Automatically generated by MockGen. DO NOT EDIT!
// Source: handler/store/datastore.go

package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
)

// Mock of DataStore interface
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *_MockDataStoreRecorder
}

// Recorder for MockDataStore (not exported)
type _MockDataStoreRecorder struct {
	mock *MockDataStore
}

func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &_MockDataStoreRecorder{mock}
	return mock
}

func (_m *MockDataStore) EXPECT() *_MockDataStoreRecorder {
	return _m.recorder
}

func (_m *MockDataStore) GetWithPrefix(keyPrefix string) (map[string]string, error) {
	ret := _m.ctrl.Call(_m, "GetWithPrefix", keyPrefix)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDataStoreRecorder) GetWithPrefix(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetWithPrefix", arg0)
}

func (_m *MockDataStore) Get(key string) (map[string]string, error) {
	ret := _m.ctrl.Call(_m, "Get", key)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDataStoreRecorder) Get(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0)
}

func (_m *MockDataStore) Add(key string, value string) error {
	ret := _m.ctrl.Call(_m, "Add", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockDataStoreRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Add", arg0, arg1)
}

func (_m *MockDataStore) StreamWithPrefix(ctx context.Context, keyPrefix string) (chan map[string]string, error) {
	ret := _m.ctrl.Call(_m, "StreamWithPrefix", ctx, keyPrefix)
	ret0, _ := ret[0].(chan map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDataStoreRecorder) StreamWithPrefix(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "StreamWithPrefix", arg0, arg1)
}