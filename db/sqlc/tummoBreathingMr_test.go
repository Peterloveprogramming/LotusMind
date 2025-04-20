package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateTummoBreathingMr(t *testing.T) {
	CreateRandomTummoBreathingMr(t, testQueries)
}

func TestUpdateTummoBreathingMrByUniqueID(t *testing.T) {
	tummoBreathingMr := CreateRandomTummoBreathingMr(t, testQueries)
	args := UpdateTummoBreathingMrByUniqueIDParams{
		UniqueID:         tummoBreathingMr.UniqueID,
		Uuid:             tummoBreathingMr.Uuid,
		StartMoodRating:  1,
		StartMood:        "Happy",
		FinishMoodRating: 1,
		FinishMood:       "Happy",
		SessionCompleted: 1,
		StartedAt:        tummoBreathingMr.StartedAt,
		EndsAt:           tummoBreathingMr.EndsAt,
		DeletedAt:        tummoBreathingMr.DeletedAt,
	}
	updatedTummoBreathing, err := testQueries.UpdateTummoBreathingMrByUniqueID(context.Background(), args)
	require.NotEmpty(t, updatedTummoBreathing)
	require.NoError(t, err)
	require.Equal(t, updatedTummoBreathing.UniqueID, tummoBreathingMr.UniqueID)
	require.Equal(t, updatedTummoBreathing.Uuid, tummoBreathingMr.Uuid)
	require.Equal(t, updatedTummoBreathing.StartMoodRating, int16(1))
	require.Equal(t, updatedTummoBreathing.StartMood, "Happy")
	require.Equal(t, updatedTummoBreathing.FinishMood, "Happy")
	require.Equal(t, updatedTummoBreathing.SessionCompleted, int16(1))
	require.Equal(t, updatedTummoBreathing.StartedAt, tummoBreathingMr.StartedAt)
	require.Equal(t, updatedTummoBreathing.EndsAt, tummoBreathingMr.EndsAt)
	require.Equal(t, updatedTummoBreathing.DeletedAt, tummoBreathingMr.DeletedAt)
}

func TestUpdateTummoBreathingMrByUuid(t *testing.T) {
	tummoBreathingMr := CreateRandomTummoBreathingMr(t, testQueries)
	args := UpdateTummoBreathingMrByUuidParams{
		Uuid:             tummoBreathingMr.Uuid,
		StartMoodRating:  2,
		StartMood:        "Sad",
		FinishMoodRating: 2,
		FinishMood:       "MoreSad",
		SessionCompleted: 1,
		StartedAt:        tummoBreathingMr.StartedAt,
		EndsAt:           tummoBreathingMr.EndsAt,
		DeletedAt:        tummoBreathingMr.DeletedAt,
	}
	updatedTummoBreathing, err := testQueries.UpdateTummoBreathingMrByUuid(context.Background(), args)
	require.NotEmpty(t, updatedTummoBreathing)
	require.NoError(t, err)
	require.Equal(t, updatedTummoBreathing.UniqueID, tummoBreathingMr.UniqueID)
	require.Equal(t, updatedTummoBreathing.Uuid, tummoBreathingMr.Uuid)
	require.Equal(t, updatedTummoBreathing.StartMoodRating, int16(2))
	require.Equal(t, updatedTummoBreathing.StartMood, "Sad")
	require.Equal(t, updatedTummoBreathing.FinishMood, "MoreSad")
	require.Equal(t, updatedTummoBreathing.SessionCompleted, int16(1))
	require.Equal(t, updatedTummoBreathing.StartedAt, tummoBreathingMr.StartedAt)
	require.Equal(t, updatedTummoBreathing.EndsAt, tummoBreathingMr.EndsAt)
	require.Equal(t, updatedTummoBreathing.DeletedAt, tummoBreathingMr.DeletedAt)
}

func TestUpdateTummoBreathingMrStartingMoodByUuid(t *testing.T) {
	tummoBreathingMr := CreateRandomTummoBreathingMr(t, testQueries)
	updatedStartingMoodRating := 2
	updatedStartingMood := "Sad"
	args := UpdateTummoBreathingMrStartingMoodByUuidParams{
		Uuid:            tummoBreathingMr.Uuid,
		StartMoodRating: int16(updatedStartingMoodRating),
		StartMood:       updatedStartingMood,
	}
	updatedTummoBreathing, err := testQueries.UpdateTummoBreathingMrStartingMoodByUuid(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTummoBreathing)
	require.Equal(t, updatedStartingMoodRating, int(updatedTummoBreathing.StartMoodRating))
	require.Equal(t, updatedStartingMood, updatedTummoBreathing.StartMood)
}

func TestUpdateTummoBreathingMrFinishingMoodByUuid(t *testing.T) {
	tummoBreathingMr := CreateRandomTummoBreathingMr(t, testQueries)
	updatedFinishingMoodRating := 10
	updatedFinishingMood := "Happy"
	args := UpdateTummoBreathingMrFinishingMoodByUuidParams{
		Uuid:             tummoBreathingMr.Uuid,
		FinishMoodRating: int16(updatedFinishingMoodRating),
		FinishMood:       updatedFinishingMood,
		SessionCompleted: int16(1),
	}
	updatedTummoBreathing, err := testQueries.UpdateTummoBreathingMrFinishingMoodByUuid(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTummoBreathing)
	require.Equal(t, updatedFinishingMoodRating, int(updatedTummoBreathing.FinishMoodRating))
	require.Equal(t, updatedFinishingMood, updatedTummoBreathing.FinishMood)
	require.Equal(t, int16(1), updatedTummoBreathing.SessionCompleted)
}

func TestUpdateTummoBreathingMrQuitByUuid(t *testing.T) {
	tummoBreathingMr := CreateRandomTummoBreathingMr(t, testQueries)
	endsTime := time.Now()
	args := UpdateTummoBreathingMrQuitByUuidParams{
		Uuid:   tummoBreathingMr.Uuid,
		EndsAt: endsTime,
	}
	updatedTummoBreathingMr, err := testQueries.UpdateTummoBreathingMrQuitByUuid(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTummoBreathingMr)
	require.WithinDuration(t, endsTime.UTC(), updatedTummoBreathingMr.EndsAt.UTC(), time.Microsecond)
}

func TestGetTummoBreathingMrByUniqueID(t *testing.T) {
	tummoBreathingMr := CreateRandomTummoBreathingMr(t, testQueries)
	fetchedTummoBreathingMr, err := testQueries.GetTummoBreathingMrByUniqueID(context.Background(), tummoBreathingMr.UniqueID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedTummoBreathingMr)
	require.Equal(t, tummoBreathingMr.UniqueID, fetchedTummoBreathingMr.UniqueID)
	require.Equal(t, tummoBreathingMr.Uuid, fetchedTummoBreathingMr.Uuid)
}

func TestGetTummoBreathingMrByUuid(t *testing.T) {
	tummoBreathingMr := CreateRandomTummoBreathingMr(t, testQueries)
	fetchedTummoBreathingMr, err := testQueries.GetTummoBreathingMrByUuid(context.Background(), tummoBreathingMr.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedTummoBreathingMr)
	require.Equal(t, tummoBreathingMr.UniqueID, fetchedTummoBreathingMr.UniqueID)
	require.Equal(t, tummoBreathingMr.Uuid, fetchedTummoBreathingMr.Uuid)

}

func TestGetTummoBreathingMrByUserID(t *testing.T) {
	// first create a tummBreathing
	// session log and user entries also get created along the way.
	tummoBreathingMr := CreateRandomTummoBreathingMr(t, testQueries)
	// fetch session log based on tibetanSingingBowl.uuid
	sessionLog, err := testQueries.GetSessionLogByUUID(context.Background(), tummoBreathingMr.Uuid)
	require.NotEmpty(t, sessionLog)
	require.NoError(t, err)

	// returns an array
	fetchedTummoBreathings, err := testQueries.GetTummoBreathingMrByUserID(context.Background(), sessionLog.UserID)
	require.NotEmpty(t, fetchedTummoBreathings)
	require.NoError(t, err)

	fetchedTummoBreathing := fetchedTummoBreathings[0]

	require.Equal(t, fetchedTummoBreathing.UniqueID, tummoBreathingMr.UniqueID)
	require.Equal(t, fetchedTummoBreathing.Uuid, tummoBreathingMr.Uuid)
}

func TestSoftDeleteTummoBreathingMrByUniqueID(t *testing.T) {
	tummoBreathingMr := CreateRandomTummoBreathingMr(t, testQueries)

	deletedTummoBreathingMr, err := testQueries.SoftDeleteTummoBreathingMrByUniqueID(context.Background(), tummoBreathingMr.UniqueID)
	require.NotEmpty(t, deletedTummoBreathingMr)
	require.NoError(t, err)
	require.NotEqual(t, tummoBreathingMr.DeletedAt.UTC(), deletedTummoBreathingMr.DeletedAt.UTC())

}
