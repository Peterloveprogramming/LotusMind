package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/lotusMind/meditation/db/mock"
	"github.com/lotusMind/meditation/token"
	"github.com/lotusMind/meditation/util"
	"github.com/stretchr/testify/require"
)

func TestFetchUserInfoByIdParams(t *testing.T) {
	//Load configuration
	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	user := randomUser(t)

	ctrl := gomock.NewController(t)
	// checks to see if all methods expected to be called were called
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	// build stubs
	store.EXPECT().GetUserById(gomock.Any(), user.ID).Times(1).Return(user, nil)

	//start test server and send in requests.
	server, err := NewServer(config, store)
	require.NoError(t, err, store)

	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/user/get_info/%d", user.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Add the Authorization header
	addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, server.config.AccessTokenDuration)

	server.router.ServeHTTP(recorder, request)

	// require response status codes match
	require.Equal(t, http.StatusOK, recorder.Code)

	// require body match
	result, err := ioutil.ReadAll(recorder.Body)
	var userResult userInfoResult
	err = json.Unmarshal(result, &userResult)
	require.NoError(t, err)
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
