package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUserProfileMobile(t *testing.T) {
	CreateRandomMobileProfile(t, testQueries)
}

func TestGetUserProfileMobileByUserId(t *testing.T) {
	userId := CreateRandomMobileProfile(t, testQueries)
	userMobileProfile, err := testQueries.GetUserProfileMobileByUserId(context.Background(), userId)
	require.NoError(t, err)
	require.NotEmpty(t, userMobileProfile)
}

func TestGetUserProfileMobileTime(t *testing.T) {
	userId := CreateRandomMobileProfile(t, testQueries)
	totalTimeSpent, err := testQueries.GetUserProfileMobileTime(context.Background(), userId)
	require.NoError(t, err)
	// // 0 because user just got created
	require.Equal(t, int64(0), totalTimeSpent)
}

func TestUpdateUserProfileMobilerTime(t *testing.T) {
	userId := CreateRandomMobileProfile(t, testQueries)
	updatedTime := 10
	args := UpdateUserProfileMobilerTimeParams{
		TotalTimeSpentInMins: int64(updatedTime),
		UserID:               int64(userId),
	}
	err := testQueries.UpdateUserProfileMobilerTime(context.Background(), args)
	require.NoError(t, err)

	totalTimeSpent, err := testQueries.GetUserProfileMobileTime(context.Background(), userId)
	require.NoError(t, err)
	// 0 because user just got created
	require.Equal(t, totalTimeSpent, int64(updatedTime))
}
