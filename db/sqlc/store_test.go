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
