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

func TestUpdateUser(t *testing.T) {
	user := CreateRandomMrUser(t, testQueries)
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
	updatedUser, err := testQueries.UpdateUser(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.FirstName, updatedUser.FirstName)
	require.Equal(t, user.LastName, updatedUser.LastName)
	require.Equal(t, user.Gender, updatedUser.Gender)
	require.Equal(t, user.BirthDate.UTC(), updatedUser.BirthDate.UTC())
	require.Equal(t, user.Country, updatedUser.Country)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, user.Goals, updatedUser.Goals)
	require.Equal(t, int16(1), updatedUser.IsMobileUser)
	require.Equal(t, int16(0), updatedUser.IsMrUser)
}

func TestDeleteUser(t *testing.T) {
	createdUser := CreateRandomMrUser(t, testQueries)

	//delete user
	err := testQueries.DeleteUser(context.Background(), createdUser.ID)
	require.NoError(t, err)

	// now get the deleted suer.
	deletedUser, err := testQueries.GetUserById(context.Background(), createdUser.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedUser)

	require.NotEmpty(t, deletedUser.DeletedAt)

}
