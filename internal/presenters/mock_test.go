// Code generated by MockGen. DO NOT EDIT.
// Source: ./contracts.go

// Package presenters is a generated GoMock package.
package presenters

import (
	context "context"
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/valerii-smirnov/petli-test-task/internal/domain"
)

// MockAuthUsecase is a mock of AuthUsecase interface.
type MockAuthUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUsecaseMockRecorder
}

// MockAuthUsecaseMockRecorder is the mock recorder for MockAuthUsecase.
type MockAuthUsecaseMockRecorder struct {
	mock *MockAuthUsecase
}

// NewMockAuthUsecase creates a new mock instance.
func NewMockAuthUsecase(ctrl *gomock.Controller) *MockAuthUsecase {
	mock := &MockAuthUsecase{ctrl: ctrl}
	mock.recorder = &MockAuthUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUsecase) EXPECT() *MockAuthUsecaseMockRecorder {
	return m.recorder
}

// SignIn mocks base method.
func (m *MockAuthUsecase) SignIn(ctx context.Context, si domain.SingIn) (domain.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", ctx, si)
	ret0, _ := ret[0].(domain.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAuthUsecaseMockRecorder) SignIn(ctx, si interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuthUsecase)(nil).SignIn), ctx, si)
}

// SignUp mocks base method.
func (m *MockAuthUsecase) SignUp(ctx context.Context, su domain.SignUp) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", ctx, su)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthUsecaseMockRecorder) SignUp(ctx, su interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuthUsecase)(nil).SignUp), ctx, su)
}

// MockDogUsecase is a mock of DogUsecase interface.
type MockDogUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockDogUsecaseMockRecorder
}

// MockDogUsecaseMockRecorder is the mock recorder for MockDogUsecase.
type MockDogUsecaseMockRecorder struct {
	mock *MockDogUsecase
}

// NewMockDogUsecase creates a new mock instance.
func NewMockDogUsecase(ctrl *gomock.Controller) *MockDogUsecase {
	mock := &MockDogUsecase{ctrl: ctrl}
	mock.recorder = &MockDogUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDogUsecase) EXPECT() *MockDogUsecaseMockRecorder {
	return m.recorder
}

// AddReaction mocks base method.
func (m *MockDogUsecase) AddReaction(ctx context.Context, userID uuid.UUID, reaction domain.Reaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReaction", ctx, userID, reaction)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddReaction indicates an expected call of AddReaction.
func (mr *MockDogUsecaseMockRecorder) AddReaction(ctx, userID, reaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReaction", reflect.TypeOf((*MockDogUsecase)(nil).AddReaction), ctx, userID, reaction)
}

// Create mocks base method.
func (m *MockDogUsecase) Create(ctx context.Context, dog domain.Dog) (domain.Dog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, dog)
	ret0, _ := ret[0].(domain.Dog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockDogUsecaseMockRecorder) Create(ctx, dog interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDogUsecase)(nil).Create), ctx, dog)
}

// Delete mocks base method.
func (m *MockDogUsecase) Delete(ctx context.Context, dogID, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, dogID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDogUsecaseMockRecorder) Delete(ctx, dogID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDogUsecase)(nil).Delete), ctx, dogID, userID)
}

// Get mocks base method.
func (m *MockDogUsecase) Get(ctx context.Context, dogID uuid.UUID) (domain.Dog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, dogID)
	ret0, _ := ret[0].(domain.Dog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDogUsecaseMockRecorder) Get(ctx, dogID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDogUsecase)(nil).Get), ctx, dogID)
}

// List mocks base method.
func (m *MockDogUsecase) List(ctx context.Context, userID uuid.UUID, pagination domain.Pagination) (domain.DogList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, userID, pagination)
	ret0, _ := ret[0].(domain.DogList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockDogUsecaseMockRecorder) List(ctx, userID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDogUsecase)(nil).List), ctx, userID, pagination)
}

// Matches mocks base method.
func (m *MockDogUsecase) Matches(ctx context.Context, userID, dogID uuid.UUID, pagination domain.Pagination) (domain.DogList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Matches", ctx, userID, dogID, pagination)
	ret0, _ := ret[0].(domain.DogList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Matches indicates an expected call of Matches.
func (mr *MockDogUsecaseMockRecorder) Matches(ctx, userID, dogID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Matches", reflect.TypeOf((*MockDogUsecase)(nil).Matches), ctx, userID, dogID, pagination)
}

// Update mocks base method.
func (m *MockDogUsecase) Update(ctx context.Context, dogID uuid.UUID, dog domain.Dog) (domain.Dog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, dogID, dog)
	ret0, _ := ret[0].(domain.Dog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockDogUsecaseMockRecorder) Update(ctx, dogID, dog interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDogUsecase)(nil).Update), ctx, dogID, dog)
}

// MockTokenParser is a mock of TokenParser interface.
type MockTokenParser struct {
	ctrl     *gomock.Controller
	recorder *MockTokenParserMockRecorder
}

// MockTokenParserMockRecorder is the mock recorder for MockTokenParser.
type MockTokenParserMockRecorder struct {
	mock *MockTokenParser
}

// NewMockTokenParser creates a new mock instance.
func NewMockTokenParser(ctrl *gomock.Controller) *MockTokenParser {
	mock := &MockTokenParser{ctrl: ctrl}
	mock.recorder = &MockTokenParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenParser) EXPECT() *MockTokenParserMockRecorder {
	return m.recorder
}

// Parse mocks base method.
func (m *MockTokenParser) Parse(token string) (*jwt.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", token)
	ret0, _ := ret[0].(*jwt.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parse indicates an expected call of Parse.
func (mr *MockTokenParserMockRecorder) Parse(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockTokenParser)(nil).Parse), token)
}

// MockIdentityExtractor is a mock of IdentityExtractor interface.
type MockIdentityExtractor struct {
	ctrl     *gomock.Controller
	recorder *MockIdentityExtractorMockRecorder
}

// MockIdentityExtractorMockRecorder is the mock recorder for MockIdentityExtractor.
type MockIdentityExtractorMockRecorder struct {
	mock *MockIdentityExtractor
}

// NewMockIdentityExtractor creates a new mock instance.
func NewMockIdentityExtractor(ctrl *gomock.Controller) *MockIdentityExtractor {
	mock := &MockIdentityExtractor{ctrl: ctrl}
	mock.recorder = &MockIdentityExtractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIdentityExtractor) EXPECT() *MockIdentityExtractorMockRecorder {
	return m.recorder
}

// ExtractFromContext mocks base method.
func (m *MockIdentityExtractor) ExtractFromContext(c *gin.Context) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExtractFromContext", c)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExtractFromContext indicates an expected call of ExtractFromContext.
func (mr *MockIdentityExtractorMockRecorder) ExtractFromContext(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExtractFromContext", reflect.TypeOf((*MockIdentityExtractor)(nil).ExtractFromContext), c)
}

// MockPaginator is a mock of Paginator interface.
type MockPaginator struct {
	ctrl     *gomock.Controller
	recorder *MockPaginatorMockRecorder
}

// MockPaginatorMockRecorder is the mock recorder for MockPaginator.
type MockPaginatorMockRecorder struct {
	mock *MockPaginator
}

// NewMockPaginator creates a new mock instance.
func NewMockPaginator(ctrl *gomock.Controller) *MockPaginator {
	mock := &MockPaginator{ctrl: ctrl}
	mock.recorder = &MockPaginatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPaginator) EXPECT() *MockPaginatorMockRecorder {
	return m.recorder
}

// GetPagination mocks base method.
func (m *MockPaginator) GetPagination(c *gin.Context) (domain.Pagination, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPagination", c)
	ret0, _ := ret[0].(domain.Pagination)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPagination indicates an expected call of GetPagination.
func (mr *MockPaginatorMockRecorder) GetPagination(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPagination", reflect.TypeOf((*MockPaginator)(nil).GetPagination), c)
}
