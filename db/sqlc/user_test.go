package db

import (
	"context"
	"testing"
	"time"

	"github.com/lotusMind/meditation/util"
	"github.com/stretchr/testify/require"
)

func createRandomMrUser(t *testing.T) User {
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

	user, err := testQueries.CreateUser(context.Background(), args)
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

func createRandomUsers(t *testing.T, numberOfUsers int) {
	for i := 0; i < numberOfUsers; i++ {
		createRandomMrUser(t)
	}
}

// Crud Testing
// Create
func TestCreateUser(t *testing.T) {
	createRandomMrUser(t)
}

func TestGetUsersByCountry(t *testing.T) {
	// first, create 10 users to make sure we have neough data.
	createRandomUsers(t, 10)

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
	createdUser := createRandomMrUser(t)

	// get user based on their email
	retrievedUser, err := testQueries.GetUserByEmail(context.Background(), createdUser.Email)
	require.NoError(t, err)
	require.NotEmpty(t, createdUser)

	// make sure email match
	require.Equal(t, createdUser.Email, retrievedUser.Email)
}

func TestGetUserByID(t *testing.T) {
	createdUser := createRandomMrUser(t)

	// get user based on their id
	retrievedUser, err := testQueries.GetUserById(context.Background(), createdUser.ID)
	require.NoError(t, err)
	require.NotEmpty(t, createdUser)

	// make sure id match
	require.Equal(t, createdUser.ID, retrievedUser.ID)
}
