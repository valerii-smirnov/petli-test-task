package presenters

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/valerii-smirnov/petli-test-task/internal/domain"
	"github.com/valerii-smirnov/petli-test-task/internal/presenters/messages"
	httpErrors "github.com/valerii-smirnov/petli-test-task/pkg/errors/http"
	"github.com/valerii-smirnov/petli-test-task/pkg/errors/ierr"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuth_SignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	mockAuthUsecase := NewMockAuthUsecase(controller)

	email := "test@email.com"
	password := "testhashedpassword"

	type fields struct {
		auth *Auth
	}
	tests := []struct {
		name              string
		fields            fields
		mocksInitFn       func()
		getRequestFn      func() *http.Request
		resultAssertionFn func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "request body validation error",
			fields: fields{
				NewAuth(mockAuthUsecase),
			},
			mocksInitFn: func() {},
			getRequestFn: func() *http.Request {
				rb := messages.SignUpRequestBody{
					Email:    "worng-email",
					Password: password,
				}

				b, err := json.Marshal(rb)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				auth: NewAuth(mockAuthUsecase),
			},
			mocksInitFn: func() {
				up := domain.SignUp{
					Email:    email,
					Password: password,
				}

				err := ierr.New(ierr.Internal, "test-error")

				mockAuthUsecase.EXPECT().SignUp(gomock.Any(), gomock.Eq(up)).Return(err)
			},
			getRequestFn: func() *http.Request {
				rb := messages.SignUpRequestBody{
					Email:    email,
					Password: password,
				}

				b, err := json.Marshal(rb)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)

				var httpErr messages.InternalServerError
				if err := json.Unmarshal(recorder.Body.Bytes(), &httpErr); err != nil {
					assert.Error(t, err)
				}

				assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
				assert.Equal(t, httpErrors.InternalServerErrorDefaultText, httpErr.Message)
			},
		},
		{
			name: "success",
			fields: fields{
				auth: NewAuth(mockAuthUsecase),
			},
			mocksInitFn: func() {
				up := domain.SignUp{
					Email:    email,
					Password: password,
				}

				mockAuthUsecase.EXPECT().SignUp(gomock.Any(), gomock.Eq(up)).Return(nil)
			},
			getRequestFn: func() *http.Request {
				rb := messages.SignUpRequestBody{
					Email:    email,
					Password: password,
				}

				b, err := json.Marshal(rb)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInitFn()

			recorder := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(recorder)
			engine = InitRoutes(engine, tt.fields.auth)

			req := tt.getRequestFn()
			engine.ServeHTTP(recorder, req)
			tt.resultAssertionFn(recorder)
		})
	}
}

func TestAuth_SignIn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	mockAuthUsecase := NewMockAuthUsecase(controller)

	email := "test@email.com"
	password := "testhashedpassword"

	token := "jwttoken"

	type fields struct {
		auth *Auth
	}
	tests := []struct {
		name              string
		fields            fields
		mocksInitFn       func()
		getRequestFn      func() *http.Request
		resultAssertionFn func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "request body validation error",
			fields: fields{
				auth: NewAuth(mockAuthUsecase),
			},
			mocksInitFn: func() {},
			getRequestFn: func() *http.Request {

				rb := messages.SignInRequestBody{
					Email:    "wrong-email",
					Password: password,
				}

				b, err := json.Marshal(rb)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				auth: NewAuth(mockAuthUsecase),
			},
			mocksInitFn: func() {
				si := domain.SingIn{
					Email:    email,
					Password: password,
				}

				err := ierr.New(ierr.Internal, "test-error")

				mockAuthUsecase.EXPECT().SignIn(gomock.Any(), gomock.Eq(si)).Return(domain.Token(""), err)
			},
			getRequestFn: func() *http.Request {

				rb := messages.SignInRequestBody{
					Email:    email,
					Password: password,
				}

				b, err := json.Marshal(rb)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				auth: NewAuth(mockAuthUsecase),
			},
			mocksInitFn: func() {
				si := domain.SingIn{
					Email:    email,
					Password: password,
				}

				mockAuthUsecase.EXPECT().SignIn(gomock.Any(), gomock.Eq(si)).Return(domain.Token(token), nil)
			},
			getRequestFn: func() *http.Request {

				rb := messages.SignInRequestBody{
					Email:    email,
					Password: password,
				}

				b, err := json.Marshal(rb)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				var si messages.SignInResponseBody
				if err := json.Unmarshal(recorder.Body.Bytes(), &si); err != nil {
					assert.Error(t, err)
				}

				assert.Equal(t, token, si.Token)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInitFn()

			recorder := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(recorder)
			engine = InitRoutes(engine, tt.fields.auth)

			req := tt.getRequestFn()
			engine.ServeHTTP(recorder, req)
			tt.resultAssertionFn(recorder)
		})
	}
}
