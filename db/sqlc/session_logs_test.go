package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateRandomSession(t *testing.T) {
	// need to make sure that user with id 1 exists in the database.
	createRandomSessionLog(t, testQueries)
}

func TestGetSessionLogByUuid(t *testing.T) {
	sessionLog := createRandomSessionLog(t, testQueries)
	retrievedSessionLog, err := testQueries.GetSessionLogByUUID(context.Background(), sessionLog.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedSessionLog)
}

func TestGetUserIdFromSessionLogUuid(t *testing.T) {
	sessionLog := createRandomSessionLog(t, testQueries)

	userId, err := testQueries.GetUserIdFromSessionLogUuid(context.Background(), sessionLog.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, userId)

	User, error := testQueries.GetUserById(context.Background(), userId)
	require.NotEmpty(t, User)
	require.NoError(t, error)
	require.Equal(t, User.ID, userId)
}

func TestGetSessionLogByUserId(t *testing.T) {
	// first create 10 random Session logs with user id 1
	for i := 0; i < 10; i++ {
		createRandomSessionLog(t, testQueries)
	}

	// need to make sure that user with id 1 exists in the database.
	SessionLogs, err := testQueries.GetSessionsLogByUserId(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, SessionLogs)

	for i := 0; i < len(SessionLogs); i++ {
		currentSessionLog := SessionLogs[i]
		require.Equal(t, currentSessionLog.UserID, int64(1))
	}
}

func TestGetSessionLogByOffset(t *testing.T) {
	// Page 1: Use OFFSET 0 and LIMIT 10 to fetch the first 10 rows.
	// Page 2: Use OFFSET 10 and LIMIT 10 to fetch the next 10 rows.

	for i := 0; i < 10; i++ {
		createRandomSessionLog(t, testQueries)
	}
	args := GetSessionsLogByUserIdWithOffsetParams{
		UserID: 1,
		Limit:  5,
		Offset: 0,
	}

	//get the first 5 session Logs
	firstSessionLogs, err := testQueries.GetSessionsLogByUserIdWithOffset(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, firstSessionLogs)
	require.Equal(t, len(firstSessionLogs), 5)

	// get the next 5 session logs
	args = GetSessionsLogByUserIdWithOffsetParams{
		UserID: 1,
		Limit:  5,
		Offset: 5,
	}
	secondSessionLogs, err := testQueries.GetSessionsLogByUserIdWithOffset(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, secondSessionLogs)
	require.Equal(t, len(secondSessionLogs), 5)

	for i := 0; i < 5; i++ {
		sessionLogFromTheFirstBatch := firstSessionLogs[i]
		sessionLogFromTheSecondBatch := secondSessionLogs[i]

		require.NotEqual(t, sessionLogFromTheFirstBatch.Uuid, sessionLogFromTheSecondBatch.Uuid)
		require.Greater(t, sessionLogFromTheSecondBatch.CreatedAt, sessionLogFromTheFirstBatch.CreatedAt)
	}
}

// Delete
func TestDeleteSessionLog(t *testing.T) {

	createdSessionLog := createRandomSessionLog(t, testQueries)

	// deleting session log
	err := testQueries.DeleteSessionLog(context.Background(), createdSessionLog.Uuid)
	require.NoError(t, err)

	// get the deleted session log
	deletedSessionLog, err := testQueries.GetSessionLogByUUID(context.Background(), createdSessionLog.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, deletedSessionLog)

	// make sure the deleted_at is not nil since we are implementing a soft delete
	require.NotEmpty(t, deletedSessionLog.DeletedAt)
}
