package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUserProfileMr(t *testing.T) {
	CreateRandomMrProfile(t, testQueries)
}

func TestGetUserProfileMrByUserId(t *testing.T) {
	userId := CreateRandomMrProfile(t, testQueries)
	userMrProfile, err := testQueries.GetUserProfileMrByUserId(context.Background(), userId)
	require.NoError(t, err)
	require.NotEmpty(t, userMrProfile)
}

func TestGetUserProfileMrTime(t *testing.T) {
	userId := CreateRandomMrProfile(t, testQueries)
	totalTimeSpent, err := testQueries.GetUserProfileMrTime(context.Background(), userId)
	require.NoError(t, err)
	// 0 because user just got created
	require.Equal(t, int64(0), totalTimeSpent)
}

func TestUpdateUserProfileMrTime(t *testing.T) {
	userId := CreateRandomMrProfile(t, testQueries)
	updatedTime := 10
	args := UpdateUserProfileMrTimeParams{
		TotalTimeSpentInMins: int64(updatedTime),
		UserID:               int64(userId),
	}
	err := testQueries.UpdateUserProfileMrTime(context.Background(), args)
	require.NoError(t, err)

	totalTimeSpent, err := testQueries.GetUserProfileMrTime(context.Background(), userId)
	require.NoError(t, err)
	// 0 because user just got created
	require.Equal(t, totalTimeSpent, int64(updatedTime))
}
