package db

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lotusMind/meditation/util"
	"github.com/stretchr/testify/require"
)

func TestCreateSessionLogTransaction(t *testing.T) {
	store := NewStore(testDB)

	sessionType := "tibetan_singing_bowl_mr"

	args := CreateSessionLogParams{
		UserID:          1,
		SessionType:     sessionType,
		SessionPlatform: "mr",
	}

	var wg sync.WaitGroup
	errs := make(chan error, 5)
	results := make(chan CreateSessionLogTransactionResult, 5)

	n := 5
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			result, err := store.CreateSessionLogTransaction(context.Background(), args)
			errs <- err
			results <- result
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errs)
	close(results)

	// Check results and errors after they have been returned from goroutines
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		parsedUUID, err := uuid.Parse(result.UUID)

		TSBSession, err := store.GetTibetanSingingBowlMrByUuid(context.Background(), parsedUUID)
		require.NoError(t, err)
		require.NotEmpty(t, TSBSession)

		SessionLog, err := store.GetSessionLogByUUID(context.Background(), parsedUUID)
		require.NoError(t, err)
		require.NotEmpty(t, SessionLog)
		require.Equal(t, result.UUID, TSBSession.Uuid.String())
	}
}

func TestCreateUserTransaction(t *testing.T) {
	store := NewStore(testDB)

	// parse date
	format := "2006-01-02"
	birthDate, err := time.Parse(format, "1990-01-01")
	require.NoError(t, err)

	// create a a mr user
	args := CreateUserTransactiontArgs{
		Email:          util.RandomEmail(),
		FirstName:      util.RandomString(5),
		LastName:       util.RandomString(5),
		Gender:         util.RandomGender(),
		Birthdate:      birthDate,
		Country:        util.RandomCountryCode(),
		Platform:       "mr",
		Goal:           "I want improve my spiritual level",
		HashedPassword: "1234567",
	}
	userResult, err := store.CreateUserTransaction(context.Background(), args)

	require.NotEmpty(t, userResult)
	require.NoError(t, err)

	userMrProfile, err := store.GetUserProfileMrByUserId(context.Background(), userResult.ID)
	require.NoError(t, err)
	require.NotEmpty(t, userMrProfile)
	require.Equal(t, userMrProfile.UserID, userResult.ID)

	userMobileProfile, err := store.GetUserProfileMobileByUserId(context.Background(), userResult.ID)
	require.Error(t, err)
	require.Empty(t, userMobileProfile)
}

func TestUpdateSessionFinishTransaction(t *testing.T) {
	store := NewStore(testDB)
	var wg sync.WaitGroup
	// testing tibetansinging bowl mr

	// what do we want to test?
	// 1. create 5 different mr users
	createUserResults := make(chan CreateUserResult, 5)
	createUserResultsErrs := make(chan error, 5)
	// 2. create 5 different session logs for each user - tibetan singing bowl
	createSessionLogResults := make(chan CreateSessionLogTransactionResult, 5)
	createSessionLogResultsErrs := make(chan error, 5)
	// 3. update tibetan singing bowl with starting mood for the 5 users
	tbStartingMoodUpdateErrs := make(chan error, 5)
	// 4. update tibetan singing bowl with finishing mood for the 5 users using UpdateSessionFinishTransaction
	tbFinishMoodUpdateErrs := make(chan error, 5)
	// 5. check if the users mrProfile has the correct time

	n := 5
	for i := 0; i < n; i++ {
		wg.Add(1) //
		go func() {
			defer wg.Done() // Decrement the counter when done
			// create a a mr user
			args := CreateUserTransactiontArgs{
				Email:          util.RandomEmail(),
				FirstName:      util.RandomString(5),
				LastName:       util.RandomString(5),
				Gender:         util.RandomGender(),
				Birthdate:      time.Time{},
				Country:        util.RandomCountryCode(),
				Platform:       "mr",
				Goal:           "I want improve my spiritual level",
				HashedPassword: "1234567",
			}
			createUserResult, err := store.CreateUserTransaction(context.Background(), args)

			createUserResultsErrs <- err
			createUserResults <- createUserResult

			// create a tibetan singing bowl session log for the user
			sessionType := "tibetan_singing_bowl_mr"
			createSessionLogArgs := CreateSessionLogParams{
				UserID:          createUserResult.ID,
				SessionType:     sessionType,
				SessionPlatform: "mr",
			}
			createSessionLogResult, err := store.CreateSessionLogTransaction(context.Background(), createSessionLogArgs)
			createSessionLogResultsErrs <- err
			createSessionLogResults <- createSessionLogResult

			// update the tibetan singing bowl session log with starting mood
			updateSessionStartingMoodArgs := UpdateTibetanSingingBowlMrStartingMoodByUuidParams{
				Uuid:            uuid.MustParse(createSessionLogResult.UUID),
				StartMoodRating: int16(5),
				StartMood:       util.RandomString(5),
			}
			_, tbStartingMoodUpdateErr := store.UpdateTibetanSingingBowlMrStartingMoodByUuid(context.Background(), updateSessionStartingMoodArgs)
			tbStartingMoodUpdateErrs <- tbStartingMoodUpdateErr

			// update the tibetan singing bowl session log with finishing mood
			updateSessionFinishMoodArgs := UpdateSessionFinishTransactionParams{
				Uuid:             uuid.MustParse(createSessionLogResult.UUID),
				FinishMoodRating: int16(5),
				FinishMood:       util.RandomString(5),
				SessionCompleted: int16(1),
				SessionType:      sessionType,
				EndsAt:           time.Now().Add(time.Minute * time.Duration(util.RandomInt(1, 60))),
			}
			tbFinishMoodUpdateErr := store.UpdateSessionFinishTransaction(context.Background(), updateSessionFinishMoodArgs)
			tbFinishMoodUpdateErrs <- tbFinishMoodUpdateErr
		}()
	}
	wg.Wait()
	close(createUserResultsErrs)
	close(createUserResults)
	for i := 0; i < n; i++ {
		// check user
		createUserResultsErr := <-createUserResultsErrs
		require.NoError(t, createUserResultsErr)
		createUserResult := <-createUserResults
		require.NotEmpty(t, createUserResult)
		require.NotEmpty(t, createUserResult.ID)

		// check session log created correctly
		createSessionLogResult := <-createSessionLogResults
		require.NotEmpty(t, createSessionLogResult)
		require.NotEmpty(t, createSessionLogResult.UUID)

		// check no erros with updating starting mood
		tbSessionStartingMoodUpdateErr := <-tbStartingMoodUpdateErrs
		require.NoError(t, tbSessionStartingMoodUpdateErr)

		// check no erros with updating finishing mood
		tbSessionFinishMoodUpdateErr := <-tbFinishMoodUpdateErrs
		require.NoError(t, tbSessionFinishMoodUpdateErr)

		// check if the users mrProfile has the correct time
		userMrProfile, err := store.GetUserProfileMrByUserId(context.Background(), createUserResult.ID)
		require.NoError(t, err)
		require.NotEmpty(t, userMrProfile)

		TibetanSingingBowlMrResultArray, err := store.GetTibetanSingingBowlMrByUserID(context.Background(), createUserResult.ID)
		require.NoError(t, err)
		require.NotEmpty(t, TibetanSingingBowlMrResultArray)
		firstTibetanSingingBowlMrResult := TibetanSingingBowlMrResultArray[0]
		timeSpentInMins := int64(firstTibetanSingingBowlMrResult.EndsAt.Sub(firstTibetanSingingBowlMrResult.StartedAt).Minutes())
		require.Equal(t, timeSpentInMins, userMrProfile.TotalTimeSpentInMins)
	}
}
func TestCreateUserForTestingDeletion(t *testing.T) {
	store := NewStore(testDB)

	// Use WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup
	errs := make(chan error, 5) // Buffered channel to prevent blocking
	results := make(chan CreateUserForTestingDeletionResult, 5)

	n := 5 // number of concurrent goroutines
	for i := 0; i < n; i++ {
		wg.Add(1) // Increment the wait group counter
		go func() {
			defer wg.Done() // Decrement the counter when done

			result, err := store.CreateUserForTestingDeletion(context.Background())
			errs <- err
			results <- result
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errs)
	close(results)

	// Check results and errors after they have been returned from goroutines
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// first check if both session logs have the correct user id - for tummo breathing and tibetan singing bowl
		require.Equal(t, result.UserId, result.TibetanSingingBowlMrSessionLogUserId)
		require.Equal(t, result.UserId, result.TummoBreathingMrSessionLogUserId)

		// then check the uuid in session logs match session tables
		require.Equal(t, result.TibetanSingingBowlMrSessionLogUuid, result.TibetanSingingBowlMrUuid)
		require.Equal(t, result.TummoBreathingMrSessionLogUuid, result.TummoBreathingMrUuid)
	}
}
