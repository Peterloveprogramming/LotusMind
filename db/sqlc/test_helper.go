package db

import (
	"context"
	"testing"
	"time"

	"github.com/lotusMind/meditation/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomMrUser(t *testing.T, q *Queries) User {
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

	user, err := q.CreateUser(context.Background(), args)
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

func CreateRandomUsers(t *testing.T, numberOfUsers int, q *Queries) {
	for i := 0; i < numberOfUsers; i++ {
		CreateRandomMrUser(t, q)
	}
}
