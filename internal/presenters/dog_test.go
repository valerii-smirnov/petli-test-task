package presenters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/valerii-smirnov/petli-test-task/internal/domain"
	"github.com/valerii-smirnov/petli-test-task/internal/presenters/messages"
	"github.com/valerii-smirnov/petli-test-task/pkg/errors/ierr"
	"github.com/valerii-smirnov/petli-test-task/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDog_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	tokenProcessor := token.NewJWT("test-secret", time.Minute*5)

	mockDogUsecase := NewMockDogUsecase(controller)
	mockIdentityExtractor := NewMockIdentityExtractor(controller)
	mockPaginator := NewMockPaginator(controller)
	authMiddleware := NewAuthMiddleware(tokenProcessor)

	userID := uuid.New()

	dList := domain.DogList{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Name:      "dog1",
			Sex:       "male",
			Age:       2,
			Breed:     "test",
			Image:     "http://test.com/dog1.jpeg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			UserID:    userID,
			Name:      "dog2",
			Sex:       "feamle",
			Age:       3,
			Breed:     "test",
			Image:     "http://test.com/dog2.jpeg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mList := messages.DogListResponseBody{
		{
			ID:    dList[0].ID.String(),
			Name:  dList[0].Name,
			Sex:   dList[0].Sex.String(),
			Age:   dList[0].Age,
			Breed: dList[0].Breed,
			Image: dList[0].Image,
		},
		{
			ID:    dList[1].ID.String(),
			Name:  dList[1].Name,
			Sex:   dList[1].Sex.String(),
			Age:   dList[1].Age,
			Breed: dList[1].Breed,
			Image: dList[1].Image,
		},
	}

	type fields struct {
		dog *Dog
	}
	tests := []struct {
		name              string
		fields            fields
		mocksInitFn       func()
		getRequestFn      func() *http.Request
		resultAssertionFn func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "getting user ID from context error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Unauthenticated, "testing-error")
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(uuid.Nil, err)
			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/dog", nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "getting pagination from url error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.InvalidArgument, "testing-error")
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockPaginator.EXPECT().GetPagination(gomock.Any()).Return(domain.Pagination{}, err)
			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/dog", nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Internal, "testing-error")
				pagination := domain.Pagination{Page: 1, PerPage: 10}

				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockPaginator.EXPECT().GetPagination(gomock.Any()).Return(pagination, nil)
				mockDogUsecase.EXPECT().List(gomock.Any(), userID, pagination).Return(nil, err)
			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/dog", nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "success",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				pagination := domain.Pagination{Page: 1, PerPage: 2}

				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockPaginator.EXPECT().GetPagination(gomock.Any()).Return(pagination, nil)
				mockDogUsecase.EXPECT().List(gomock.Any(), userID, pagination).Return(dList, nil)
			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/dog", nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				var resp messages.DogListResponseBody
				if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
					assert.Error(t, err)
				}

				assert.Equal(t, mList, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInitFn()

			recorder := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(recorder)
			engine = InitRoutes(engine, tt.fields.dog)

			req := tt.getRequestFn()
			engine.ServeHTTP(recorder, req)
			tt.resultAssertionFn(recorder)
		})
	}
}

func TestDog_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	tokenProcessor := token.NewJWT("test-secret", time.Minute*5)

	mockDogUsecase := NewMockDogUsecase(controller)
	mockIdentityExtractor := NewMockIdentityExtractor(controller)
	mockPaginator := NewMockPaginator(controller)
	authMiddleware := NewAuthMiddleware(tokenProcessor)

	userID := uuid.New()
	dogID := uuid.New()

	dDog := domain.Dog{
		ID:        dogID,
		UserID:    userID,
		Name:      "dog1",
		Sex:       "male",
		Age:       3,
		Breed:     "test",
		Image:     "http://test.com/image1.jpeg",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	mDog := messages.DogResponseBody{
		ID:    dDog.ID.String(),
		Name:  dDog.Name,
		Sex:   dDog.Sex.String(),
		Age:   dDog.Age,
		Breed: dDog.Breed,
		Image: dDog.Image,
	}

	type fields struct {
		dog *Dog
	}
	tests := []struct {
		name              string
		fields            fields
		mocksInitFn       func()
		getRequestFn      func() *http.Request
		resultAssertionFn func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "getting dog id param error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {

			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/dog/wrong-dog-id-param", nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Internal, "testing-error")
				mockDogUsecase.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(domain.Dog{}, err)
			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/dog/%s", dogID.String()), nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "success",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				mockDogUsecase.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(dDog, nil)
			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/dog/%s", dogID.String()), nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				var resp messages.DogResponseBody
				if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
					assert.Error(t, err)
				}

				assert.Equal(t, mDog, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInitFn()

			recorder := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(recorder)
			engine = InitRoutes(engine, tt.fields.dog)

			req := tt.getRequestFn()
			engine.ServeHTTP(recorder, req)
			tt.resultAssertionFn(recorder)
		})
	}
}

func TestDog_Matches(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	tokenProcessor := token.NewJWT("test-secret", time.Minute*5)

	mockDogUsecase := NewMockDogUsecase(controller)
	mockIdentityExtractor := NewMockIdentityExtractor(controller)
	mockPaginator := NewMockPaginator(controller)
	authMiddleware := NewAuthMiddleware(tokenProcessor)

	userID := uuid.New()
	dogID := uuid.New()

	dList := domain.DogList{
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Name:      "dog1",
			Sex:       "male",
			Age:       4,
			Breed:     "test-breed",
			Image:     "http://test.com/dog1.jpeg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Name:      "dog2",
			Sex:       "female",
			Age:       6,
			Breed:     "test-breed",
			Image:     "http://test.com/dog3.jpeg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mList := messages.DogListResponseBody{
		{
			ID:    dList[0].ID.String(),
			Name:  dList[0].Name,
			Sex:   dList[0].Sex.String(),
			Age:   dList[0].Age,
			Breed: dList[0].Breed,
			Image: dList[0].Image,
		},
		{
			ID:    dList[1].ID.String(),
			Name:  dList[1].Name,
			Sex:   dList[1].Sex.String(),
			Age:   dList[1].Age,
			Breed: dList[1].Breed,
			Image: dList[1].Image,
		},
	}

	type fields struct {
		dog *Dog
	}
	tests := []struct {
		name              string
		fields            fields
		mocksInitFn       func()
		getRequestFn      func() *http.Request
		resultAssertionFn func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "getting pagination error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.InvalidArgument, "testing-error")
				mockPaginator.EXPECT().GetPagination(gomock.Any()).Return(domain.Pagination{}, err)
			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/dog/wrong-dog-id-param/matches", nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "getting dog id from params error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				mockPaginator.EXPECT().GetPagination(gomock.Any()).Return(domain.Pagination{Page: 1, PerPage: 2}, nil)

			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/dog/wrong-dog-id-param/matches", nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "identity extractor error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Internal, "testing error")

				mockPaginator.EXPECT().GetPagination(gomock.Any()).Return(domain.Pagination{Page: 1, PerPage: 2}, nil)
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(uuid.Nil, err)

			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/dog/%s/matches", dogID.String()), nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Internal, "testing error")
				pag := domain.Pagination{Page: 1, PerPage: 2}

				mockPaginator.EXPECT().GetPagination(gomock.Any()).Return(pag, nil)
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockDogUsecase.EXPECT().Matches(gomock.Any(), userID, dogID, pag).Return(nil, err)

			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/dog/%s/matches", dogID.String()), nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "success",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				pag := domain.Pagination{Page: 1, PerPage: 2}

				mockPaginator.EXPECT().GetPagination(gomock.Any()).Return(pag, nil)
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockDogUsecase.EXPECT().Matches(gomock.Any(), userID, dogID, pag).Return(dList, nil)

			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/dog/%s/matches", dogID.String()), nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				var resp messages.DogListResponseBody
				if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
					assert.Error(t, err)
				}

				assert.Equal(t, mList, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInitFn()

			recorder := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(recorder)
			engine = InitRoutes(engine, tt.fields.dog)

			engine.ServeHTTP(recorder, tt.getRequestFn())
			tt.resultAssertionFn(recorder)
		})
	}
}

func TestDog_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	tokenProcessor := token.NewJWT("test-secret", time.Minute*5)

	mockDogUsecase := NewMockDogUsecase(controller)
	mockIdentityExtractor := NewMockIdentityExtractor(controller)
	mockPaginator := NewMockPaginator(controller)
	authMiddleware := NewAuthMiddleware(tokenProcessor)

	userID := uuid.New()
	wrongDogRequestBody := messages.CreateOrUpdateDogRequestBody{
		Name:  "dog1",
		Sex:   "unknown",
		Age:   40,
		Breed: "test",
		Image: "http://test.com/dog1.jpeg",
	}

	validDogRequestBody := messages.CreateOrUpdateDogRequestBody{
		Name:  "dog1",
		Sex:   "male",
		Age:   15,
		Breed: "test",
		Image: "http://test.com/dog1.jpeg",
	}

	domainDogIN := domain.Dog{
		UserID: userID,
		Name:   validDogRequestBody.Name,
		Sex:    domain.DogSex(validDogRequestBody.Sex),
		Age:    validDogRequestBody.Age,
		Breed:  validDogRequestBody.Breed,
		Image:  validDogRequestBody.Image,
	}

	domainDogOut := domain.Dog{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      domainDogIN.Name,
		Sex:       domainDogIN.Sex,
		Age:       domainDogIN.Age,
		Image:     domainDogIN.Image,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	responseBody := messages.DogResponseBody{
		ID:    domainDogOut.ID.String(),
		Name:  domainDogOut.Name,
		Sex:   domainDogOut.Sex.String(),
		Age:   domainDogOut.Age,
		Breed: domainDogOut.Breed,
		Image: domainDogOut.Image,
	}

	type fields struct {
		dog *Dog
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
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {

			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(wrongDogRequestBody)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/dog", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "getting pagination error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Internal, "testing-error")
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(uuid.Nil, err)
			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(validDogRequestBody)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/dog", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Internal, "testing-error")
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockDogUsecase.EXPECT().Create(gomock.Any(), gomock.Eq(domainDogIN)).Return(domain.Dog{}, err)
			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(validDogRequestBody)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/dog", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockDogUsecase.EXPECT().Create(gomock.Any(), gomock.Eq(domainDogIN)).Return(domainDogOut, nil)
			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(validDogRequestBody)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/dog", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				var resp messages.DogResponseBody
				if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
					assert.Error(t, err)
				}

				assert.Equal(t, responseBody, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInitFn()

			recorder := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(recorder)
			engine = InitRoutes(engine, tt.fields.dog)

			engine.ServeHTTP(recorder, tt.getRequestFn())
			tt.resultAssertionFn(recorder)
		})
	}
}

func TestDog_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	tokenProcessor := token.NewJWT("test-secret", time.Minute*5)

	mockDogUsecase := NewMockDogUsecase(controller)
	mockIdentityExtractor := NewMockIdentityExtractor(controller)
	mockPaginator := NewMockPaginator(controller)
	authMiddleware := NewAuthMiddleware(tokenProcessor)

	userID := uuid.New()
	dogID := uuid.New()
	wrongDogRequestBody := messages.CreateOrUpdateDogRequestBody{
		Name:  "dog1",
		Sex:   "unknown",
		Age:   40,
		Breed: "test",
		Image: "http://test.com/dog1.jpeg",
	}

	validDogRequestBody := messages.CreateOrUpdateDogRequestBody{
		Name:  "dog1",
		Sex:   "male",
		Age:   15,
		Breed: "test",
		Image: "http://test.com/dog1.jpeg",
	}

	domainDogIN := domain.Dog{
		UserID: userID,
		Name:   validDogRequestBody.Name,
		Sex:    domain.DogSex(validDogRequestBody.Sex),
		Age:    validDogRequestBody.Age,
		Breed:  validDogRequestBody.Breed,
		Image:  validDogRequestBody.Image,
	}

	domainDogOut := domain.Dog{
		ID:        dogID,
		UserID:    userID,
		Name:      domainDogIN.Name,
		Sex:       domainDogIN.Sex,
		Age:       domainDogIN.Age,
		Image:     domainDogIN.Image,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	responseBody := messages.DogResponseBody{
		ID:    domainDogOut.ID.String(),
		Name:  domainDogOut.Name,
		Sex:   domainDogOut.Sex.String(),
		Age:   domainDogOut.Age,
		Breed: domainDogOut.Breed,
		Image: domainDogOut.Image,
	}

	type fields struct {
		dog *Dog
	}
	tests := []struct {
		name              string
		fields            fields
		mocksInitFn       func()
		getRequestFn      func() *http.Request
		resultAssertionFn func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "getting dog id from params error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {

			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(wrongDogRequestBody)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPut, "/api/dog/wrong-dog-id", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "request body validation error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {

			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(wrongDogRequestBody)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/dog/%s", dogID.String()), bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "identity extractor error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Internal, "testing-error")
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(uuid.Nil, err)
			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(validDogRequestBody)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/dog/%s", dogID.String()), bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Internal, "testing-error")
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockDogUsecase.EXPECT().Update(gomock.Any(), gomock.Eq(dogID), gomock.Eq(domainDogIN)).Return(domain.Dog{}, err)
			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(validDogRequestBody)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/dog/%s", dogID.String()), bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockDogUsecase.EXPECT().Update(gomock.Any(), gomock.Eq(dogID), gomock.Eq(domainDogIN)).Return(domainDogOut, nil)
			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(validDogRequestBody)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/dog/%s", dogID.String()), bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				var resp messages.DogResponseBody
				if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
					assert.Error(t, err)
				}

				assert.Equal(t, responseBody, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInitFn()

			recorder := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(recorder)
			engine = InitRoutes(engine, tt.fields.dog)

			engine.ServeHTTP(recorder, tt.getRequestFn())
			tt.resultAssertionFn(recorder)
		})
	}
}

func TestDog_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	tokenProcessor := token.NewJWT("test-secret", time.Minute*5)

	mockDogUsecase := NewMockDogUsecase(controller)
	mockIdentityExtractor := NewMockIdentityExtractor(controller)
	mockPaginator := NewMockPaginator(controller)
	authMiddleware := NewAuthMiddleware(tokenProcessor)

	userID := uuid.New()
	dogID := uuid.New()

	type fields struct {
		dog *Dog
	}
	tests := []struct {
		name              string
		fields            fields
		mocksInitFn       func()
		getRequestFn      func() *http.Request
		resultAssertionFn func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "getting dog id from params error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodDelete, "/api/dog/wrong-dog-id", nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "identity extractor error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Internal, "testing-error")
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(uuid.Nil, err)
			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/dog/%s", dogID.String()), nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Internal, "testing-error")
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockDogUsecase.EXPECT().Delete(gomock.Any(), gomock.Eq(dogID), gomock.Eq(userID)).Return(err)
			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/dog/%s", dogID.String()), nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockDogUsecase.EXPECT().Delete(gomock.Any(), gomock.Eq(dogID), gomock.Eq(userID)).Return(nil)
			},
			getRequestFn: func() *http.Request {
				req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/dog/%s", dogID.String()), nil)
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

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
			engine = InitRoutes(engine, tt.fields.dog)

			engine.ServeHTTP(recorder, tt.getRequestFn())
			tt.resultAssertionFn(recorder)
		})
	}
}

func TestDog_Reaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	tokenProcessor := token.NewJWT("test-secret", time.Minute*5)

	mockDogUsecase := NewMockDogUsecase(controller)
	mockIdentityExtractor := NewMockIdentityExtractor(controller)
	mockPaginator := NewMockPaginator(controller)
	authMiddleware := NewAuthMiddleware(tokenProcessor)

	userID := uuid.New()
	//dogID := uuid.New()

	wrongReaction := messages.ReactionRequestBody{
		Liker:  "wrong-id",
		Liked:  "wrong-id",
		Action: "wrong-action",
	}

	likerID := uuid.New()
	likedID := uuid.New()

	validReaction := messages.ReactionRequestBody{
		Liker:  likerID.String(),
		Liked:  likedID.String(),
		Action: "like",
	}

	domainReaction := domain.Reaction{
		Liker:  likerID,
		Liked:  likedID,
		Action: domain.Like,
	}

	type fields struct {
		dog *Dog
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
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(wrongReaction)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/dog/reaction", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "identity extractor error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.InvalidArgument, "test-error")
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(uuid.Nil, err)
			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(validReaction)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/dog/reaction", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				err := ierr.New(ierr.Internal, "testing-error")

				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockDogUsecase.EXPECT().AddReaction(gomock.Any(), gomock.Eq(userID), gomock.Eq(domainReaction)).
					Return(err)
			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(validReaction)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/dog/reaction", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

				return req
			},
			resultAssertionFn: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "usecase error",
			fields: fields{
				dog: NewDog(mockDogUsecase, mockIdentityExtractor, mockPaginator, authMiddleware.Auth),
			},
			mocksInitFn: func() {
				mockIdentityExtractor.EXPECT().ExtractFromContext(gomock.Any()).Return(userID, nil)
				mockDogUsecase.EXPECT().AddReaction(gomock.Any(), gomock.Eq(userID), gomock.Eq(domainReaction)).
					Return(nil)
			},
			getRequestFn: func() *http.Request {
				b, err := json.Marshal(validReaction)
				if err != nil {
					assert.Error(t, err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/dog/reaction", bytes.NewReader(b))
				if err != nil {
					assert.Error(t, err)
				}

				st, err := tokenProcessor.Generate(userID)
				if err != nil {
					assert.Error(t, err)
				}

				req.Header.Set(AuthorizationHeaderName, fmt.Sprintf("%s%s", bearerPrefix, st))

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
			engine = InitRoutes(engine, tt.fields.dog)

			engine.ServeHTTP(recorder, tt.getRequestFn())
			tt.resultAssertionFn(recorder)
		})
	}
}
