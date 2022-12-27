// Code generated by MockGen. DO NOT EDIT.
// Source: ./contracts.go

// Package usecases is a generated GoMock package.
package usecases

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/valerii-smirnov/petli-test-task/internal/domain"
)

// MockDogAdapter is a mock of DogAdapter interface.
type MockDogAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockDogAdapterMockRecorder
}

// MockDogAdapterMockRecorder is the mock recorder for MockDogAdapter.
type MockDogAdapterMockRecorder struct {
	mock *MockDogAdapter
}

// NewMockDogAdapter creates a new mock instance.
func NewMockDogAdapter(ctrl *gomock.Controller) *MockDogAdapter {
	mock := &MockDogAdapter{ctrl: ctrl}
	mock.recorder = &MockDogAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDogAdapter) EXPECT() *MockDogAdapterMockRecorder {
	return m.recorder
}

// AddReaction mocks base method.
func (m *MockDogAdapter) AddReaction(ctx context.Context, reaction domain.Reaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReaction", ctx, reaction)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddReaction indicates an expected call of AddReaction.
func (mr *MockDogAdapterMockRecorder) AddReaction(ctx, reaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReaction", reflect.TypeOf((*MockDogAdapter)(nil).AddReaction), ctx, reaction)
}

// Create mocks base method.
func (m *MockDogAdapter) Create(ctx context.Context, dog domain.Dog) (domain.Dog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, dog)
	ret0, _ := ret[0].(domain.Dog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockDogAdapterMockRecorder) Create(ctx, dog interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDogAdapter)(nil).Create), ctx, dog)
}

// Delete mocks base method.
func (m *MockDogAdapter) Delete(ctx context.Context, dogID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, dogID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDogAdapterMockRecorder) Delete(ctx, dogID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDogAdapter)(nil).Delete), ctx, dogID)
}

// Get mocks base method.
func (m *MockDogAdapter) Get(ctx context.Context, dogID uuid.UUID) (domain.Dog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, dogID)
	ret0, _ := ret[0].(domain.Dog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDogAdapterMockRecorder) Get(ctx, dogID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDogAdapter)(nil).Get), ctx, dogID)
}

// List mocks base method.
func (m *MockDogAdapter) List(ctx context.Context, userID uuid.UUID, pagination domain.Pagination) (domain.DogList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, userID, pagination)
	ret0, _ := ret[0].(domain.DogList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockDogAdapterMockRecorder) List(ctx, userID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDogAdapter)(nil).List), ctx, userID, pagination)
}

// Matches mocks base method.
func (m *MockDogAdapter) Matches(ctx context.Context, dogID uuid.UUID, pagination domain.Pagination) (domain.DogList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Matches", ctx, dogID, pagination)
	ret0, _ := ret[0].(domain.DogList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Matches indicates an expected call of Matches.
func (mr *MockDogAdapterMockRecorder) Matches(ctx, dogID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Matches", reflect.TypeOf((*MockDogAdapter)(nil).Matches), ctx, dogID, pagination)
}

// Update mocks base method.
func (m *MockDogAdapter) Update(ctx context.Context, dogID uuid.UUID, dog domain.Dog) (domain.Dog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, dogID, dog)
	ret0, _ := ret[0].(domain.Dog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockDogAdapterMockRecorder) Update(ctx, dogID, dog interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDogAdapter)(nil).Update), ctx, dogID, dog)
}

// MockUserAdapter is a mock of UserAdapter interface.
type MockUserAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockUserAdapterMockRecorder
}

// MockUserAdapterMockRecorder is the mock recorder for MockUserAdapter.
type MockUserAdapterMockRecorder struct {
	mock *MockUserAdapter
}

// NewMockUserAdapter creates a new mock instance.
func NewMockUserAdapter(ctrl *gomock.Controller) *MockUserAdapter {
	mock := &MockUserAdapter{ctrl: ctrl}
	mock.recorder = &MockUserAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserAdapter) EXPECT() *MockUserAdapterMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserAdapter) Create(ctx context.Context, su domain.SignUp) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, su)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserAdapterMockRecorder) Create(ctx, su interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserAdapter)(nil).Create), ctx, su)
}

// Exists mocks base method.
func (m *MockUserAdapter) Exists(ctx context.Context, email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockUserAdapterMockRecorder) Exists(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockUserAdapter)(nil).Exists), ctx, email)
}

// Get mocks base method.
func (m *MockUserAdapter) Get(ctx context.Context, si domain.SingIn) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, si)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserAdapterMockRecorder) Get(ctx, si interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserAdapter)(nil).Get), ctx, si)
}

// MockPasswordHasher is a mock of PasswordHasher interface.
type MockPasswordHasher struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordHasherMockRecorder
}

// MockPasswordHasherMockRecorder is the mock recorder for MockPasswordHasher.
type MockPasswordHasherMockRecorder struct {
	mock *MockPasswordHasher
}

// NewMockPasswordHasher creates a new mock instance.
func NewMockPasswordHasher(ctrl *gomock.Controller) *MockPasswordHasher {
	mock := &MockPasswordHasher{ctrl: ctrl}
	mock.recorder = &MockPasswordHasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPasswordHasher) EXPECT() *MockPasswordHasherMockRecorder {
	return m.recorder
}

// StringHash mocks base method.
func (m *MockPasswordHasher) StringHash(in string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StringHash", in)
	ret0, _ := ret[0].(string)
	return ret0
}

// StringHash indicates an expected call of StringHash.
func (mr *MockPasswordHasherMockRecorder) StringHash(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StringHash", reflect.TypeOf((*MockPasswordHasher)(nil).StringHash), in)
}

// MockTokenGenerator is a mock of TokenGenerator interface.
type MockTokenGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockTokenGeneratorMockRecorder
}

// MockTokenGeneratorMockRecorder is the mock recorder for MockTokenGenerator.
type MockTokenGeneratorMockRecorder struct {
	mock *MockTokenGenerator
}

// NewMockTokenGenerator creates a new mock instance.
func NewMockTokenGenerator(ctrl *gomock.Controller) *MockTokenGenerator {
	mock := &MockTokenGenerator{ctrl: ctrl}
	mock.recorder = &MockTokenGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenGenerator) EXPECT() *MockTokenGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockTokenGenerator) Generate(userID uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockTokenGeneratorMockRecorder) Generate(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockTokenGenerator)(nil).Generate), userID)
}