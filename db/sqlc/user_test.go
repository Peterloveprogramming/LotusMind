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

// Read
// func TestGetUsersByPlatform(t *testing.T) {
// 	// first, create 10 users to make sure we have neough data.
// 	createRandomUsers(t, 10)

// 	platformToCheck := "mobile"

// 	// get all the users who are on the mobile platform
// 	users, err := testQueries.GetUsersByPlatform(context.Background(), platformToCheck)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, users)

// 	// loop through each user to make sure their platform is mobile
// 	for i := 0; i < len(users); i++ {
// 		currentUser := users[i]
// 		require.Equal(t, currentUser.Platform, platformToCheck)
// 	}
// }

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

// Update
// func TestUpdateUser(t *testing.T) {
// 	format := "2006-01-02"
// 	birthDate, err := time.Parse(format, "1990-01-01")
// 	require.NoError(t, err)

// 	createdUser := createRandomUser(t)
// 	args := UpdateUserParams{
// 		ID:             createdUser.ID,
// 		Email:          util.RandomString(10) + "@example.com", // Generate a random email
// 		FirstName:      util.RandomString(5),                   // Generate a random first name
// 		LastName:       util.RandomString(5),                   // Generate a random last name
// 		Gender:         util.RandomGender(),                    // Randomly choose gender
// 		BirthDate:      birthDate,                              // Generate a random birth date
// 		Country:        util.RandomCountryCode(),               // Randomly choose a country code
// 		HashedPassword: util.RandomString(12),                  // Assuming a random string for hashed password
// 		Goals:          "I want to learn meditation!",          // Static goals
// 		Platform:       util.RandomPlatform(),                  // Randomly choose platform
// 	}

// 	updatedUser, err := testQueries.UpdateUser(context.Background(), args)
// 	require.NotEmpty(t, updatedUser)
// 	require.NoError(t, err)
// 	require.Equal(t, args.Email, updatedUser.Email)
// 	require.Equal(t, args.FirstName, updatedUser.FirstName)
// 	require.Equal(t, args.LastName, updatedUser.LastName)
// 	require.Equal(t, args.Gender, updatedUser.Gender)
// 	require.Equal(t, args.BirthDate.UTC(), updatedUser.BirthDate.UTC()) // Compare in UTC
// 	require.Equal(t, args.Country, updatedUser.Country)
// 	require.Equal(t, args.HashedPassword, updatedUser.HashedPassword)
// 	require.Equal(t, args.Goals, updatedUser.Goals)
// 	require.Equal(t, args.Platform, updatedUser.Platform)
// }

// // Delete
func TestDeleteUser(t *testing.T) {

	// first create a entry in 4 tables, users, session_logs, tibetan singing bowl and tummo breathing and make sure all records are linked
	store := NewStore(testDB)

	result, err := store.CreateUserForTestingDeletion(context.Background())
	require.NoError(t, err)

	require.NotEmpty(t, result)

	require.Equal(t, result.UserId, result.TibetanSingingBowlMrSessionLogUserId)
	require.Equal(t, result.UserId, result.TummoBreathingMrSessionLogUserId)
	require.Equal(t, result.TibetanSingingBowlMrSessionLogUuid, result.TibetanSingingBowlMrUuid)
	require.Equal(t, result.TummoBreathingMrSessionLogUuid, result.TummoBreathingMrUuid)

	// deleting user
	err = testQueries.DeleteUser(context.Background(), result.UserId)
	require.NoError(t, err)

	// get deleted user and make sure deleted_at is not empty
	deletedUser, err := testQueries.GetUserById(context.Background(), result.UserId)
	require.NoError(t, err)
	require.NotEmpty(t, deletedUser.DeletedAt)

	//  get deleted session log for tibetan singing bowl and make sure deleted_at is not empty
	deletedTSBSessionLog, err := testQueries.GetSessionLogByUUID(context.Background(), result.TibetanSingingBowlMrSessionLogUuid)
	require.NoError(t, err)
	require.NotEmpty(t, deletedTSBSessionLog.DeletedAt)

	//  get deleted session log for tummo breathing and make sure deleted_at is not empty
	deletedTBSessionLog, err := testQueries.GetSessionLogByUUID(context.Background(), result.TummoBreathingMrSessionLogUuid)
	require.NoError(t, err)
	require.NotEmpty(t, deletedTBSessionLog.DeletedAt)

	//  get the deleted tibetan singing bowl and make sure deleted_at is not empty
	deletedTibetanSingingBowl, err := testQueries.GetTibetanSingingBowlMrByUuid(context.Background(), result.TibetanSingingBowlMrUuid)
	require.NoError(t, err)
	require.NotEmpty(t, deletedTibetanSingingBowl.DeletedAt)

	//  get the deleted tummo breathing and make sure deleted_at is not empty
	deletedTummoBreathing, err := testQueries.GetTummoBreathingMrByUuid(context.Background(), result.TummoBreathingMrUuid)
	require.NoError(t, err)
	require.NotEmpty(t, deletedTummoBreathing.DeletedAt)

}
