// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/project-radius/radius/pkg/queue (interfaces: Client)

// Package queue is a generated GoMock package.
package queue

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Dequeue mocks base method.
func (m *MockClient) Dequeue(arg0 context.Context, arg1 ...DequeueOptions) (<-chan *Message, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Dequeue", varargs...)
	ret0, _ := ret[0].(<-chan *Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Dequeue indicates an expected call of Dequeue.
func (mr *MockClientMockRecorder) Dequeue(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dequeue", reflect.TypeOf((*MockClient)(nil).Dequeue), varargs...)
}

// Enqueue mocks base method.
func (m *MockClient) Enqueue(arg0 context.Context, arg1 *Message, arg2 ...EnqueueOptions) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Enqueue", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Enqueue indicates an expected call of Enqueue.
func (mr *MockClientMockRecorder) Enqueue(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enqueue", reflect.TypeOf((*MockClient)(nil).Enqueue), varargs...)
}