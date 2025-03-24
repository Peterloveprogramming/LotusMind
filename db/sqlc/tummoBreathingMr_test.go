package db

import (
	"context"
	"testing"

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
