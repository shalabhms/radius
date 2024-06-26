// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/radius-project/radius/pkg/cli/credential (interfaces: AzureCredentialManagementClientInterface)
//
// Generated by this command:
//
//	mockgen -typed -destination=./mock_azure_credential_management.go -package=credential -self_package github.com/radius-project/radius/pkg/cli/credential github.com/radius-project/radius/pkg/cli/credential AzureCredentialManagementClientInterface
//

// Package credential is a generated GoMock package.
package credential

import (
	context "context"
	reflect "reflect"

	v20231001preview "github.com/radius-project/radius/pkg/ucp/api/v20231001preview"
	gomock "go.uber.org/mock/gomock"
)

// MockAzureCredentialManagementClientInterface is a mock of AzureCredentialManagementClientInterface interface.
type MockAzureCredentialManagementClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAzureCredentialManagementClientInterfaceMockRecorder
}

// MockAzureCredentialManagementClientInterfaceMockRecorder is the mock recorder for MockAzureCredentialManagementClientInterface.
type MockAzureCredentialManagementClientInterfaceMockRecorder struct {
	mock *MockAzureCredentialManagementClientInterface
}

// NewMockAzureCredentialManagementClientInterface creates a new mock instance.
func NewMockAzureCredentialManagementClientInterface(ctrl *gomock.Controller) *MockAzureCredentialManagementClientInterface {
	mock := &MockAzureCredentialManagementClientInterface{ctrl: ctrl}
	mock.recorder = &MockAzureCredentialManagementClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAzureCredentialManagementClientInterface) EXPECT() *MockAzureCredentialManagementClientInterfaceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockAzureCredentialManagementClientInterface) Delete(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockAzureCredentialManagementClientInterfaceMockRecorder) Delete(arg0, arg1 any) *MockAzureCredentialManagementClientInterfaceDeleteCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAzureCredentialManagementClientInterface)(nil).Delete), arg0, arg1)
	return &MockAzureCredentialManagementClientInterfaceDeleteCall{Call: call}
}

// MockAzureCredentialManagementClientInterfaceDeleteCall wrap *gomock.Call
type MockAzureCredentialManagementClientInterfaceDeleteCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAzureCredentialManagementClientInterfaceDeleteCall) Return(arg0 bool, arg1 error) *MockAzureCredentialManagementClientInterfaceDeleteCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAzureCredentialManagementClientInterfaceDeleteCall) Do(f func(context.Context, string) (bool, error)) *MockAzureCredentialManagementClientInterfaceDeleteCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAzureCredentialManagementClientInterfaceDeleteCall) DoAndReturn(f func(context.Context, string) (bool, error)) *MockAzureCredentialManagementClientInterfaceDeleteCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Get mocks base method.
func (m *MockAzureCredentialManagementClientInterface) Get(arg0 context.Context, arg1 string) (ProviderCredentialConfiguration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(ProviderCredentialConfiguration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAzureCredentialManagementClientInterfaceMockRecorder) Get(arg0, arg1 any) *MockAzureCredentialManagementClientInterfaceGetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAzureCredentialManagementClientInterface)(nil).Get), arg0, arg1)
	return &MockAzureCredentialManagementClientInterfaceGetCall{Call: call}
}

// MockAzureCredentialManagementClientInterfaceGetCall wrap *gomock.Call
type MockAzureCredentialManagementClientInterfaceGetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAzureCredentialManagementClientInterfaceGetCall) Return(arg0 ProviderCredentialConfiguration, arg1 error) *MockAzureCredentialManagementClientInterfaceGetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAzureCredentialManagementClientInterfaceGetCall) Do(f func(context.Context, string) (ProviderCredentialConfiguration, error)) *MockAzureCredentialManagementClientInterfaceGetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAzureCredentialManagementClientInterfaceGetCall) DoAndReturn(f func(context.Context, string) (ProviderCredentialConfiguration, error)) *MockAzureCredentialManagementClientInterfaceGetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockAzureCredentialManagementClientInterface) List(arg0 context.Context) ([]CloudProviderStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]CloudProviderStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockAzureCredentialManagementClientInterfaceMockRecorder) List(arg0 any) *MockAzureCredentialManagementClientInterfaceListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockAzureCredentialManagementClientInterface)(nil).List), arg0)
	return &MockAzureCredentialManagementClientInterfaceListCall{Call: call}
}

// MockAzureCredentialManagementClientInterfaceListCall wrap *gomock.Call
type MockAzureCredentialManagementClientInterfaceListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAzureCredentialManagementClientInterfaceListCall) Return(arg0 []CloudProviderStatus, arg1 error) *MockAzureCredentialManagementClientInterfaceListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAzureCredentialManagementClientInterfaceListCall) Do(f func(context.Context) ([]CloudProviderStatus, error)) *MockAzureCredentialManagementClientInterfaceListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAzureCredentialManagementClientInterfaceListCall) DoAndReturn(f func(context.Context) ([]CloudProviderStatus, error)) *MockAzureCredentialManagementClientInterfaceListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Put mocks base method.
func (m *MockAzureCredentialManagementClientInterface) Put(arg0 context.Context, arg1 v20231001preview.AzureCredentialResource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockAzureCredentialManagementClientInterfaceMockRecorder) Put(arg0, arg1 any) *MockAzureCredentialManagementClientInterfacePutCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockAzureCredentialManagementClientInterface)(nil).Put), arg0, arg1)
	return &MockAzureCredentialManagementClientInterfacePutCall{Call: call}
}

// MockAzureCredentialManagementClientInterfacePutCall wrap *gomock.Call
type MockAzureCredentialManagementClientInterfacePutCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAzureCredentialManagementClientInterfacePutCall) Return(arg0 error) *MockAzureCredentialManagementClientInterfacePutCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAzureCredentialManagementClientInterfacePutCall) Do(f func(context.Context, v20231001preview.AzureCredentialResource) error) *MockAzureCredentialManagementClientInterfacePutCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAzureCredentialManagementClientInterfacePutCall) DoAndReturn(f func(context.Context, v20231001preview.AzureCredentialResource) error) *MockAzureCredentialManagementClientInterfacePutCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
