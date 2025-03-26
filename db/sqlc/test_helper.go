package db

import (
	"context"
	"testing"
	"time"

	"github.com/lotusMind/meditation/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomMrUser(t *testing.T, q *Queries) User {
	// parse date
	format := "2006-01-02"
	birthDate, err := time.Parse(format, "1990-01-01")
	require.NoError(t, err)

	// hashing the password
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	args := CreateUserParams{
		Email:          util.RandomString(10) + "@example.com",
		FirstName:      util.RandomString(5),
		LastName:       util.RandomString(5),
		Gender:         util.RandomGender(),
		BirthDate:      birthDate,
		Country:        util.RandomCountryCode(),
		IsMrUser:       1,
		IsMobileUser:   0,
		Goals:          "I want improve my spiritual level",
		HashedPassword: hashedPassword,
	}

	user, err := q.CreateUser(context.Background(), args)
	require.NotEmpty(t, user)
	require.NoError(t, err)
	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.FirstName, user.FirstName)
	require.Equal(t, args.LastName, user.LastName)
	require.Equal(t, args.Gender, user.Gender)
	require.Equal(t, args.BirthDate.UTC(), user.BirthDate.UTC()) // Compare in UTC
	require.Equal(t, args.Country, user.Country)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.Goals, user.Goals)
	return user
}

func CreateRandomMobileUser(t *testing.T, q *Queries) User {
	user := CreateRandomMrUser(t, q)
	updateParams := UpdateUserParams{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Gender:         user.Gender,
		BirthDate:      user.BirthDate,
		Country:        user.Country,
		HashedPassword: user.HashedPassword,
		Goals:          user.Goals,
		IsMobileUser:   1,
		IsMrUser:       0,
	}
	updatedUser, err := q.UpdateUser(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, updatedUser.ID, user.ID)
	require.Equal(t, updatedUser.Email, user.Email)
	require.Equal(t, updatedUser.FirstName, user.FirstName)
	require.Equal(t, updatedUser.LastName, user.LastName)
	require.Equal(t, updatedUser.Gender, user.Gender)
	require.Equal(t, updatedUser.BirthDate.UTC(), user.BirthDate.UTC()) // Compare in UTC
	require.Equal(t, updatedUser.Country, user.Country)
	require.Equal(t, updatedUser.HashedPassword, user.HashedPassword)
	require.Equal(t, updatedUser.Goals, user.Goals)
	require.Equal(t, int16(1), updatedUser.IsMobileUser)
	require.Equal(t, int16(0), updatedUser.IsMrUser)
	return updatedUser
}

func CreateRandomUsers(t *testing.T, numberOfUsers int, q *Queries) {
	for i := 0; i < numberOfUsers; i++ {
		CreateRandomMrUser(t, q)
	}
}

func createRandomSessionLog(t *testing.T, q *Queries, sessionTypeOptional ...string) SessionLog {
	// first create a user otherwise it will throw an error
	user := CreateRandomMrUser(t, q)

	sessionType := util.RandomSessionType()
	if len(sessionTypeOptional) > 0 {
		sessionType = sessionTypeOptional[0]
	}
	sessionPlatform := util.GetPlatformTypeBasedOnSessionType(sessionType)
	arg := CreateSessionLogParams{
		UserID:          user.ID,
		SessionType:     sessionType,
		SessionPlatform: sessionPlatform,
	}

	sessionLog, err := q.CreateSessionLog(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sessionLog)
	require.Equal(t, arg.UserID, sessionLog.UserID)
	require.Equal(t, arg.SessionType, sessionLog.SessionType)
	require.Equal(t, arg.SessionPlatform, sessionLog.SessionPlatform)

	return sessionLog
}

func createSessionWithUserID(t *testing.T, q *Queries, userId int64) SessionLog {
	sessionType := util.RandomSessionType()
	sessionPlatform := util.GetPlatformTypeBasedOnSessionType(sessionType)
	arg := CreateSessionLogParams{
		UserID:          userId,
		SessionType:     sessionType,
		SessionPlatform: sessionPlatform,
	}

	sessionLog, err := q.CreateSessionLog(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sessionLog)
	require.Equal(t, arg.UserID, sessionLog.UserID)
	require.Equal(t, arg.SessionType, sessionLog.SessionType)
	require.Equal(t, arg.SessionPlatform, sessionLog.SessionPlatform)

	return sessionLog
}

func CreateRandomMrProfile(t *testing.T, q *Queries) int64 {
	user := CreateRandomMrUser(t, q)

	err := q.CreateUserProfileMr(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	return user.ID
}

func CreateRandomMobileProfile(t *testing.T, q *Queries) int64 {
	user := CreateRandomMobileUser(t, q)

	err := q.CreateUserProfileMobile(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	return user.ID
}

func CreateRandomTibetanSingingBowlMr(t *testing.T, q *Queries) TibetanSingingBowlMr {
	sessionLog := createRandomSessionLog(t, q, "tibetan_singing_bowl_mr")

	args := CreateTibetanSingingBowlMrParams{
		//convert to uuid
		Uuid:             sessionLog.Uuid,
		StartMoodRating:  0,
		StartMood:        "N/A",
		FinishMoodRating: 0,
		FinishMood:       "N/A",
		SessionCompleted: 0,
	}
	tibetanSingingBowl, err := q.CreateTibetanSingingBowlMr(context.Background(), args)
	require.NotEmpty(t, tibetanSingingBowl)
	require.NoError(t, err)
	return tibetanSingingBowl
}

func CreateRandomTummoBreathingMr(t *testing.T, q *Queries) TummoBreathingMr {
	sessionLog := createRandomSessionLog(t, q, "tummo_breathing_mr")

	args := CreateTummoBreathingMrParams{
		//convert to uuid
		Uuid:             sessionLog.Uuid,
		StartMoodRating:  0,
		StartMood:        "N/A",
		FinishMoodRating: 0,
		FinishMood:       "N/A",
		SessionCompleted: 0,
	}
	tummoBreathingMr, err := q.CreateTummoBreathingMr(context.Background(), args)
	require.NotEmpty(t, tummoBreathingMr)
	require.NoError(t, err)
	return tummoBreathingMr
}
