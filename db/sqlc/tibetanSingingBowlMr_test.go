package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTibetanSingingBowlMr(t *testing.T) {
	CreateRandomTibetanSingingBowlMr(t, testQueries)
}

func TestUpdateTibetanSingingBowlMrByUniqueID(t *testing.T) {
	tibetanSingingBowl := CreateRandomTibetanSingingBowlMr(t, testQueries)
	args := UpdateTibetanSingingBowlMrByUniqueIDParams{
		UniqueID:         tibetanSingingBowl.UniqueID,
		Uuid:             tibetanSingingBowl.Uuid,
		StartMoodRating:  1,
		StartMood:        "Happy",
		FinishMoodRating: 1,
		FinishMood:       "Happy",
		SessionCompleted: 1,
		StartedAt:        tibetanSingingBowl.StartedAt,
		EndsAt:           tibetanSingingBowl.EndsAt,
		DeletedAt:        tibetanSingingBowl.DeletedAt,
	}
	updatedTibetanSingingBowl, err := testQueries.UpdateTibetanSingingBowlMrByUniqueID(context.Background(), args)
	require.NotEmpty(t, updatedTibetanSingingBowl)
	require.NoError(t, err)
	require.Equal(t, updatedTibetanSingingBowl.UniqueID, tibetanSingingBowl.UniqueID)
	require.Equal(t, updatedTibetanSingingBowl.Uuid, tibetanSingingBowl.Uuid)
	require.Equal(t, updatedTibetanSingingBowl.StartMoodRating, int16(1))
	require.Equal(t, updatedTibetanSingingBowl.StartMood, "Happy")
	require.Equal(t, updatedTibetanSingingBowl.FinishMood, "Happy")
	require.Equal(t, updatedTibetanSingingBowl.SessionCompleted, int16(1))
	require.Equal(t, updatedTibetanSingingBowl.StartedAt, tibetanSingingBowl.StartedAt)
	require.Equal(t, updatedTibetanSingingBowl.EndsAt, tibetanSingingBowl.EndsAt)
	require.Equal(t, updatedTibetanSingingBowl.DeletedAt, tibetanSingingBowl.DeletedAt)
}

func TestUpdateTibetanSingingBowlMrByUuid(t *testing.T) {
	tibetanSingingBowl := CreateRandomTibetanSingingBowlMr(t, testQueries)
	args := UpdateTibetanSingingBowlMrByUuidParams{
		Uuid:             tibetanSingingBowl.Uuid,
		StartMoodRating:  2,
		StartMood:        "Sad",
		FinishMoodRating: 2,
		FinishMood:       "MoreSad",
		SessionCompleted: 1,
		StartedAt:        tibetanSingingBowl.StartedAt,
		EndsAt:           tibetanSingingBowl.EndsAt,
		DeletedAt:        tibetanSingingBowl.DeletedAt,
	}
	updatedTibetanSingingBowl, err := testQueries.UpdateTibetanSingingBowlMrByUuid(context.Background(), args)
	require.NotEmpty(t, updatedTibetanSingingBowl)
	require.NoError(t, err)
	require.Equal(t, updatedTibetanSingingBowl.UniqueID, tibetanSingingBowl.UniqueID)
	require.Equal(t, updatedTibetanSingingBowl.Uuid, tibetanSingingBowl.Uuid)
	require.Equal(t, updatedTibetanSingingBowl.StartMoodRating, int16(2))
	require.Equal(t, updatedTibetanSingingBowl.StartMood, "Sad")
	require.Equal(t, updatedTibetanSingingBowl.FinishMood, "MoreSad")
	require.Equal(t, updatedTibetanSingingBowl.SessionCompleted, int16(1))
	require.Equal(t, updatedTibetanSingingBowl.StartedAt, tibetanSingingBowl.StartedAt)
	require.Equal(t, updatedTibetanSingingBowl.EndsAt, tibetanSingingBowl.EndsAt)
	require.Equal(t, updatedTibetanSingingBowl.DeletedAt, tibetanSingingBowl.DeletedAt)
}
