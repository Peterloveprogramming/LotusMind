package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// Crud Testing
// Create
func TestCreateUser(t *testing.T) {
	CreateRandomMrUser(t, testQueries)
}

func TestGetUsersByCountry(t *testing.T) {
	// first, create 10 users to make sure we have neough data.
	CreateRandomUsers(t, 10, testQueries)

	countryToCheck := "AF"

	// get all the users who are on the mobile platform
	users, err := testQueries.GetUsersByCountry(context.Background(), countryToCheck)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	// loop through each user to make sure their platform is mobile
	for i := 0; i < len(users); i++ {
		currentUser := users[i]
		require.Equal(t, currentUser.Country, countryToCheck)
	}
}

func TestGetUserByEmail(t *testing.T) {
	createdUser := CreateRandomMrUser(t, testQueries)

	// get user based on their email
	retrievedUser, err := testQueries.GetUserByEmail(context.Background(), createdUser.Email)
	require.NoError(t, err)
	require.NotEmpty(t, createdUser)

	// make sure email match
	require.Equal(t, createdUser.Email, retrievedUser.Email)
}

func TestGetUserByID(t *testing.T) {
	createdUser := CreateRandomMrUser(t, testQueries)

	// get user based on their id
	retrievedUser, err := testQueries.GetUserById(context.Background(), createdUser.ID)
	require.NoError(t, err)
	require.NotEmpty(t, createdUser)

	// make sure id match
	require.Equal(t, createdUser.ID, retrievedUser.ID)
}
