package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	mockdb "github.com/lotusMind/meditation/db/mock"
	db "github.com/lotusMind/meditation/db/sqlc"
	"github.com/lotusMind/meditation/token"
	"github.com/lotusMind/meditation/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	password := "marvel998"
	user := randomMrUser(t, password)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		body          gin.H
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"gender":     user.Gender,
				"birth_date": user.BirthDate.Format(util.GetDateFormat()),
				"country":    user.Country,
				"goals":      user.Goals,
				"platform":   "mr",
				"password":   password,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUserTransaction(gomock.Any(), gomock.Any()).Times(1).Return(db.CreateUserResult{ID: user.ID}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// require response status codes match
				require.Equal(t, http.StatusCreated, recorder.Code)

				// require body match
				result, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var userResult db.CreateUserResult
				err = json.Unmarshal(result, &userResult)
				require.NoError(t, err)
				require.Equal(t, user.ID, userResult.ID)
			},
		},
		{
			name: "UniqueViolation",
			body: gin.H{
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"gender":     user.Gender,
				"birth_date": user.BirthDate.Format(util.GetDateFormat()),
				"country":    user.Country,
				"goals":      user.Goals,
				"platform":   "mr",
				"password":   password,
			},
			buildStub: func(store *mockdb.MockStore) {
				pqError := &pq.Error{
					Code: "23505", // Unique violation error code
				}
				store.EXPECT().CreateUserTransaction(gomock.Any(), gomock.Any()).Times(1).Return(db.CreateUserResult{}, pqError)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidPlatform",
			body: gin.H{
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"gender":     user.Gender,
				"birth_date": user.BirthDate.Format(util.GetDateFormat()),
				"country":    user.Country,
				"goals":      user.Goals,
				"platform":   "invalid", // invalid platform
				"password":   password,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUserTransaction(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			tc.buildStub(testStore)

			recorder := httptest.NewRecorder()

			// marshal the response body
			bodyData, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/user/create")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyData))
			require.NoError(t, err)

			testServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestFetchUserTime(t *testing.T) {
	user := randomMrUser(t)

	testCases := []struct {
		name          string
		userId        int64
		platform      string
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			userId:   user.ID,
			platform: "mr",
			buildStub: func(store *mockdb.MockStore) {
				var result int64
				store.EXPECT().GetUserProfileMrTime(gomock.Any(), user.ID, "mr").Times(1).Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// require response status codes match
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStub(testStore)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/user/get_time/%d/%s", tc.userId, tc.platform)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			addAuthorization(t, request, testServer.tokenMaker, authorizationTypeBearer, user.Email, testServer.config.AccessTokenDuration)
			testServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
func TestFetchUserInfoByIdParams(t *testing.T) {
	//Load configuration
	// config, err := util.LoadConfig("../")
	// require.NoError(t, err)

	user := randomMrUser(t)

	testCases := []struct {
		name          string
		accountID     int64
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: user.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), user.ID).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// require response status codes match
				require.Equal(t, http.StatusOK, recorder.Code)

				// require body match
				result, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)
				var userResult userInfoResult
				err = json.Unmarshal(result, &userResult)
				require.NoError(t, err)
				require.Equal(t, user.ID, userResult.ID)
				require.Equal(t, user.Email, userResult.Email)
				require.Equal(t, user.FirstName, userResult.FirstName)
				require.Equal(t, user.LastName, userResult.LastName)
				require.Equal(t, user.Gender, userResult.Gender)
				require.Equal(t, user.BirthDate.Format(util.GetDateFormat()), userResult.BirthDate)
				require.Equal(t, user.Country, userResult.Country)
				require.Equal(t, user.IsMrUser, userResult.IsMrUser)
				require.Equal(t, user.IsMobileUser, userResult.IsMobileUser)
				require.Equal(t, user.Goals, userResult.Goals)
			},
		},
		{
			name:      "NotFound",
			accountID: user.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), user.ID).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// require response status codes match
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: user.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), user.ID).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// require response status codes match
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidId",
			accountID: 0,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// require response status codes match
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// ctrl := gomock.NewController(t)
			// checks to see if all methods expected to be called were called
			// defer ctrl.Finish()

			// store := mockdb.NewMockStore(testStore)
			tc.buildStub(testStore)

			//start test server and send in requests.
			// server, err := NewServer(config, store)
			// require.NoError(t, err)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/user/get_info/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add the Authorization header
			addAuthorization(t, request, testServer.tokenMaker, authorizationTypeBearer, user.Email, testServer.config.AccessTokenDuration)

			testServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	userEmail string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(userEmail, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}
