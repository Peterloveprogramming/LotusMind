package db

import (
	"context"
	"testing"
	"time"

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

func TestUpdateTibetanSingingBowlMrStartingMoodByUuid(t *testing.T) {
	tibetanSingingBowl := CreateRandomTibetanSingingBowlMr(t, testQueries)
	updatedStartingMoodRating := 2
	updatedStartingMood := "Sad"
	args := UpdateTibetanSingingBowlMrStartingMoodByUuidParams{
		Uuid:            tibetanSingingBowl.Uuid,
		StartMoodRating: int16(updatedStartingMoodRating),
		StartMood:       updatedStartingMood,
	}
	updatedTibetanSingingBowl, err := testQueries.UpdateTibetanSingingBowlMrStartingMoodByUuid(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTibetanSingingBowl)
	require.Equal(t, updatedStartingMoodRating, int(updatedTibetanSingingBowl.StartMoodRating))
	require.Equal(t, updatedStartingMood, updatedTibetanSingingBowl.StartMood)
}

func TestUpdateTibetanSingingBowlMrFinishingMoodByUuid(t *testing.T) {
	tibetanSingingBowl := CreateRandomTibetanSingingBowlMr(t, testQueries)
	updatedFinishingMoodRating := 10
	updatedFinishingMood := "Happy"
	args := UpdateTibetanSingingBowlMrFinishingMoodByUuidParams{
		Uuid:             tibetanSingingBowl.Uuid,
		FinishMoodRating: int16(updatedFinishingMoodRating),
		FinishMood:       updatedFinishingMood,
		SessionCompleted: int16(1),
	}
	updatedTibetanSingingBowl, err := testQueries.UpdateTibetanSingingBowlMrFinishingMoodByUuid(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTibetanSingingBowl)
	require.Equal(t, updatedFinishingMoodRating, int(updatedTibetanSingingBowl.FinishMoodRating))
	require.Equal(t, updatedFinishingMood, updatedTibetanSingingBowl.FinishMood)
	require.Equal(t, int16(1), updatedTibetanSingingBowl.SessionCompleted)

}

func TestUpdateTibetanSingingBowlMrQuitByUuid(t *testing.T) {
	tibetanSingingBowl := CreateRandomTibetanSingingBowlMr(t, testQueries)
	endsTime := time.Now()
	args := UpdateTibetanSingingBowlMrQuitByUuidParams{
		Uuid:   tibetanSingingBowl.Uuid,
		EndsAt: endsTime,
	}
	updatedTibetanSingingBowl, err := testQueries.UpdateTibetanSingingBowlMrQuitByUuid(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTibetanSingingBowl)
	require.Equal(t, endsTime.UTC(), updatedTibetanSingingBowl.EndsAt.UTC())
}

func TestTibetanSingingBowlMrByUuid(t *testing.T) {
	tibetanSingingBowl := CreateRandomTibetanSingingBowlMr(t, testQueries)
	fetchedTibetanSingingBowl, err := testQueries.GetTibetanSingingBowlMrByUuid(context.Background(), tibetanSingingBowl.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedTibetanSingingBowl)
	require.Equal(t, tibetanSingingBowl.UniqueID, fetchedTibetanSingingBowl.UniqueID)
	require.Equal(t, tibetanSingingBowl.Uuid, fetchedTibetanSingingBowl.Uuid)

}

func TestGetTibetanSingingBowlMrByUniqueID(t *testing.T) {
	tibetanSingingBowl := CreateRandomTibetanSingingBowlMr(t, testQueries)
	fetchedTibetanSingingBowl, err := testQueries.GetTibetanSingingBowlMrByUniqueID(context.Background(), tibetanSingingBowl.UniqueID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedTibetanSingingBowl)
	require.Equal(t, tibetanSingingBowl.UniqueID, fetchedTibetanSingingBowl.UniqueID)
	require.Equal(t, tibetanSingingBowl.Uuid, fetchedTibetanSingingBowl.Uuid)

}

func TestGetTibetanSingingBowlMrByUserID(t *testing.T) {
	// first create a tibetanSingingBowl.
	// session log and user entries also get created along the way.
	tibetanSingingBowl := CreateRandomTibetanSingingBowlMr(t, testQueries)
	// fetch session log based on tibetanSingingBowl.uuid
	sessionLog, err := testQueries.GetSessionLogByUUID(context.Background(), tibetanSingingBowl.Uuid)
	require.NotEmpty(t, sessionLog)
	require.NoError(t, err)

	// returns an array
	fetchedTibetanSingingBowls, err := testQueries.GetTibetanSingingBowlMrByUserID(context.Background(), sessionLog.UserID)
	require.NotEmpty(t, fetchedTibetanSingingBowls)
	require.NoError(t, err)

	fetchedTibetanSingingBowl := fetchedTibetanSingingBowls[0]

	require.Equal(t, fetchedTibetanSingingBowl.UniqueID, tibetanSingingBowl.UniqueID)
	require.Equal(t, fetchedTibetanSingingBowl.Uuid, tibetanSingingBowl.Uuid)
}

func TestDeleteTibetanSingingBowlMrByUniqueID(t *testing.T) {
	tibetanSingingBowl := CreateRandomTibetanSingingBowlMr(t, testQueries)

	deletedTibetanSingingBowl, err := testQueries.SoftDeleteTibetanSingingBowlMrByUniqueID(context.Background(), tibetanSingingBowl.UniqueID)
	require.NotEmpty(t, deletedTibetanSingingBowl)
	require.NoError(t, err)
	require.NotEqual(t, tibetanSingingBowl.DeletedAt.UTC(), deletedTibetanSingingBowl.DeletedAt.UTC())

}
