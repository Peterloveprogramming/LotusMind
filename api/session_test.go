package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/lotusMind/meditation/db/mock"
	db "github.com/lotusMind/meditation/db/sqlc"
	"github.com/lotusMind/meditation/util"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	sessionLogTibetanSingingBowl := randomTibetanSingingBowlrSessionLog(t)
	// create a user so we can use their email to generate authorization token.
	user := randomMrUser(t)
	testCases := []struct {
		name          string
		userID        int64
		sessionType   string
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			userID:      sessionLogTibetanSingingBowl.UserID,
			sessionType: sessionLogTibetanSingingBowl.SessionType,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateSessionLogTransaction(gomock.Any(), gomock.AssignableToTypeOf(db.CreateSessionLogParams{})).
					Times(1).
					Return(db.CreateSessionLogTransactionResult{
						UUID:        sessionLogTibetanSingingBowl.Uuid.String(),
						UserId:      sessionLogTibetanSingingBowl.UserID,
						SessionType: sessionLogTibetanSingingBowl.SessionType,
					}, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				var SessionLogResult db.CreateSessionLogTransactionResult
				err := json.NewDecoder(recorder.Body).Decode(&SessionLogResult)
				require.NoError(t, err)
				require.Equal(t, sessionLogTibetanSingingBowl.Uuid.String(), SessionLogResult.UUID)
				require.Equal(t, sessionLogTibetanSingingBowl.UserID, SessionLogResult.UserId)
				require.Equal(t, sessionLogTibetanSingingBowl.SessionType, SessionLogResult.SessionType)
			},
		},
		{
			name:        "invalid session type",
			userID:      sessionLogTibetanSingingBowl.UserID,
			sessionType: "invalid",
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateSessionLogTransaction(gomock.Any(), gomock.AssignableToTypeOf(db.CreateSessionLogParams{})).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
		{
			name:        "user does not exists",
			userID:      util.RandomInt(0, 1000),
			sessionType: sessionLogTibetanSingingBowl.SessionType,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateSessionLogTransaction(gomock.Any(), gomock.AssignableToTypeOf(db.CreateSessionLogParams{})).
					Times(1).
					Return(db.CreateSessionLogTransactionResult{}, sql.ErrNoRows) // Return an error

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

			url := fmt.Sprintf("/session/create/%d/%s", tc.userID, tc.sessionType)
			request, err := http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)
			// Add the Authorization header
			addAuthorization(t, request, testServer.tokenMaker, authorizationTypeBearer, user.Email, testServer.config.AccessTokenDuration)

			testServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateSessionStartingMood(t *testing.T) {
	sessionLogTibetanSingingBowl := randomTibetanSingingBowlrSessionLog(t)
	user := randomMrUser(t)

	testCases := []struct {
		name          string
		SessionUuid   string
		sessionType   string
		body          gin.H // Request body
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			SessionUuid: sessionLogTibetanSingingBowl.Uuid.String(),
			sessionType: sessionLogTibetanSingingBowl.SessionType,
			body: gin.H{
				"start_mood_rating": 5,
				"start_mood":        "Happy",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateTibetanSingingBowlMrStartingMoodByUuid(gomock.Any(), gomock.AssignableToTypeOf(db.UpdateTibetanSingingBowlMrStartingMoodByUuidParams{})).
					Times(1).
					Return(db.TibetanSingingBowlMr{}, nil) // Return an empty struct and nil error
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
				// Add more assertions here if needed, e.g., check if the database was updated.
			},
		},
		{
			name:        "Invalid sessiontype",
			SessionUuid: sessionLogTibetanSingingBowl.Uuid.String(),
			sessionType: "invalid",
			body: gin.H{
				"start_mood_rating": 5,
				"start_mood":        "Happy",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateTibetanSingingBowlMrStartingMoodByUuid(gomock.Any(), gomock.AssignableToTypeOf(db.UpdateTibetanSingingBowlMrStartingMoodByUuidParams{})).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// Add more assertions here if needed, e.g., check if the database was updated.
			},
		},
		{
			name:        "Invalid uuid",
			SessionUuid: "123",
			sessionType: sessionLogTibetanSingingBowl.SessionType,
			body: gin.H{
				"start_mood_rating": 5,
				"start_mood":        "Happy",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateTibetanSingingBowlMrStartingMoodByUuid(gomock.Any(), gomock.AssignableToTypeOf(db.UpdateTibetanSingingBowlMrStartingMoodByUuidParams{})).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// Add more assertions here if needed, e.g., check if the database was updated.
			},
		},
		{
			name:        "server error",
			SessionUuid: sessionLogTibetanSingingBowl.Uuid.String(),
			sessionType: sessionLogTibetanSingingBowl.SessionType,
			body: gin.H{
				"start_mood_rating": 5,
				"start_mood":        "Happy",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateTibetanSingingBowlMrStartingMoodByUuid(gomock.Any(), gomock.AssignableToTypeOf(db.UpdateTibetanSingingBowlMrStartingMoodByUuidParams{})).
					Times(1).Return(db.TibetanSingingBowlMr{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				// Add more assertions here if needed, e.g., check if the database was updated.
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStub(testStore)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/session/update/start/%s/%s", tc.SessionUuid, tc.sessionType)

			// Marshal the request body to JSON
			bodyData, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyData))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json") // Set content type
			addAuthorization(t, request, testServer.tokenMaker, authorizationTypeBearer, user.Email, testServer.config.AccessTokenDuration)

			testServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateSessionFinishingMood(t *testing.T) {
	sessionLogTibetanSingingBowl := randomTibetanSingingBowlrSessionLog(t)
	user := randomMrUser(t)

	testCases := []struct {
		name          string
		SessionUuid   string
		sessionType   string
		body          gin.H // Request body
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			SessionUuid: sessionLogTibetanSingingBowl.Uuid.String(),
			sessionType: sessionLogTibetanSingingBowl.SessionType,
			body: gin.H{
				"finish_mood":        "Happy",
				"finish_mood_rating": 123,
				"ends_at":            "2024-12-28T12:34:56.789Z",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateSessionFinishTransaction(gomock.Any(), gomock.AssignableToTypeOf(db.UpdateSessionFinishTransactionParams{})).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name:        "Invalid sessiontype",
			SessionUuid: sessionLogTibetanSingingBowl.Uuid.String(),
			sessionType: "invalid",
			body: gin.H{
				"finish_mood":        "Happy",
				"finish_mood_rating": 123,
				"ends_at":            "2024-12-28T12:34:56.789Z",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateSessionFinishTransaction(gomock.Any(), gomock.AssignableToTypeOf(db.UpdateSessionFinishTransactionParams{})).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// Add more assertions here if needed, e.g., check if the database was updated.
			},
		},
		// {
		// 	name:        "Invalid uuid",
		// 	SessionUuid: "123",
		// 	sessionType: sessionLogTibetanSingingBowl.SessionType,
		// 	body: gin.H{
		// 		"finish_mood":        "Happy",
		// 		"finish_mood_rating": 123,
		// 		"ends_at":            "2024-12-28T12:34:56.789Z",
		// 	},
		// 	buildStub: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			UpdateSessionFinishTransaction(gomock.Any(), gomock.AssignableToTypeOf(db.UpdateSessionFinishTransactionParams{})).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 		// Add more assertions here if needed, e.g., check if the database was updated.
		// 	},
		// },
		{
			name:        "server error",
			SessionUuid: sessionLogTibetanSingingBowl.Uuid.String(),
			sessionType: sessionLogTibetanSingingBowl.SessionType,
			body: gin.H{
				"finish_mood":        "Happy",
				"finish_mood_rating": 123,
				"ends_at":            "2024-12-28T12:34:56.789Z",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateSessionFinishTransaction(gomock.Any(), gomock.AssignableToTypeOf(db.UpdateSessionFinishTransactionParams{})).
					Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				// Add more assertions here if needed, e.g., check if the database was updated.
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStub(testStore)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/session/update/finish/%s/%s", tc.SessionUuid, tc.sessionType)

			// Marshal the request body to JSON
			bodyData, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyData))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json") // Set content type
			addAuthorization(t, request, testServer.tokenMaker, authorizationTypeBearer, user.Email, testServer.config.AccessTokenDuration)

			testServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateSessionQuit(t *testing.T) {
	sessionLogTibetanSingingBowl := randomTibetanSingingBowlrSessionLog(t)
	user := randomMrUser(t)

	testCases := []struct {
		name          string
		SessionUuid   string
		sessionType   string
		body          gin.H // Request body
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			SessionUuid: sessionLogTibetanSingingBowl.Uuid.String(),
			sessionType: sessionLogTibetanSingingBowl.SessionType,
			body: gin.H{
				"ends_at": "2024-12-28T12:34:56.789Z",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateSessionFinishTransaction(gomock.Any(), gomock.AssignableToTypeOf(db.UpdateSessionFinishTransactionParams{})).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name:        "Invalid sessiontype",
			SessionUuid: sessionLogTibetanSingingBowl.Uuid.String(),
			sessionType: "invalid",
			body: gin.H{
				"ends_at": "2024-12-28T12:34:56.789Z",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateSessionFinishTransaction(gomock.Any(), gomock.AssignableToTypeOf(db.UpdateSessionFinishTransactionParams{})).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// Add more assertions here if needed, e.g., check if the database was updated.
			},
		},
		{
			name:        "server error",
			SessionUuid: sessionLogTibetanSingingBowl.Uuid.String(),
			sessionType: sessionLogTibetanSingingBowl.SessionType,
			body: gin.H{
				"ends_at": "2024-12-28T12:34:56.789Z",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateSessionFinishTransaction(gomock.Any(), gomock.AssignableToTypeOf(db.UpdateSessionFinishTransactionParams{})).
					Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStub(testStore)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/session/update/quit/%s/%s", tc.SessionUuid, tc.sessionType)

			// Marshal the request body to JSON
			bodyData, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyData))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json") // Set content type
			addAuthorization(t, request, testServer.tokenMaker, authorizationTypeBearer, user.Email, testServer.config.AccessTokenDuration)

			testServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
