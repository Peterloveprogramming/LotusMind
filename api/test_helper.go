package api

import (
	"testing"
	"time"

	db "github.com/lotusMind/meditation/db/sqlc"
	"github.com/lotusMind/meditation/util"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) db.User {
	// parse date
	format := "2006-01-02"
	birthDate, err := time.Parse(format, "1990-01-01")
	require.NoError(t, err)

	// hashing the password
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	return db.User{
		ID:             int64(util.RandomInt(1, 1000)),
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
}
