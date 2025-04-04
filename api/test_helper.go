package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	db "github.com/lotusMind/meditation/db/sqlc"
	"github.com/lotusMind/meditation/token"
	"github.com/lotusMind/meditation/util"
	"github.com/stretchr/testify/require"
)

func randomMrUser(t *testing.T, password ...string) db.User {
	// parse date
	format := "2006-01-02"
	birthDate, err := time.Parse(format, "1990-01-01")
	require.NoError(t, err)

	// hashing the password
	var hashedPassword string
	if len(password) > 0 {
		hashedPassword, err = util.HashPassword(password[0])
	} else {
		hashedPassword, err = util.HashPassword(util.RandomString(6))
	}
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

func randomMrProfile(t *testing.T, userIdParam ...int64) db.UsersProfileMr {
	var userId *int64
	if len(userIdParam) > 0 {
		userId = &userIdParam[0]
	} else {
		user := randomMrUser(t)
		userId = &user.ID
	}
	return db.UsersProfileMr{
		UserID:               *userId,
		TotalTimeSpentInMins: int64(util.RandomInt(1, 1000)),
	}
}

func randomMobileProfile(t *testing.T, userIdParam ...int64) db.UsersProfileMobile {
	var userId *int64
	if len(userIdParam) > 0 {
		userId = &userIdParam[0]
	} else {
		user := randomMrUser(t)
		userId = &user.ID
	}
	return db.UsersProfileMobile{
		UserID:               *userId,
		TotalTimeSpentInMins: int64(util.RandomInt(1, 1000)),
	}
}

func randomTibetanSingingBowlrSessionLog(t *testing.T) db.SessionLog {
	mrUser := randomMrUser(t)
	return db.SessionLog{
		Uuid:            uuid.New(),
		UserID:          mrUser.ID,
		SessionType:     "tibetan_singing_bowl_mr",
		SessionPlatform: "mr",
	}
}

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	userEmail string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(userEmail, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}
