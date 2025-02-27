package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lotusMind/meditation/util"
)

// store is a wrapper that provides all functions to execute db queries and transactions
type Store struct {
	*Queries         // *Queries contains all the methods for executing specific SQL Queries
	db       *sql.DB // Database connection and used for creating transaction.
}

// Newstore creates a new store and returns a pointer to Store since its big.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// use default level isolation by passing nil
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v,rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// ExecTx executes a function within a database transaction
func (store *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type CreateSessionLogTransactionResult struct {
	UserId      int64  `json:"user_id"`
	SessionType string `json:"session_type"`
	UUID        string `json:"uuid"`
}

// CreateSessionLogTransaction creates a session log and initializes session-specific data
func (store *Store) CreateSessionLogTransaction(ctx context.Context, args CreateSessionLogParams) (CreateSessionLogTransactionResult, error) {
	var result CreateSessionLogTransactionResult

	err := store.execTx(ctx, func(q *Queries) error {
		sessionLog, err := q.CreateSessionLog(ctx, args)
		if err != nil {
			if prErr, ok := err.(*pq.Error); ok {
				switch prErr.Code.Name() {
				case "foreign_key_violation":
					return fmt.Errorf("foreign_key_violation")
				}
				// log.Println(prErr.Code.Name())
			}
			return fmt.Errorf("failed to create session log: %w", err)
		}

		switch args.SessionType {
		case "tibetan_singing_bowl_mr":
			params := CreateTibetanSingingBowlMrParams{
				Uuid:             sessionLog.Uuid,
				StartMoodRating:  0,
				StartMood:        "N/A",
				FinishMoodRating: 0,
				FinishMood:       "N/A",
				SessionCompleted: 0,
			}
			_, err = q.CreateTibetanSingingBowlMr(ctx, params)

		case "tummo_breathing_mr":
			params := CreateTummoBreathingMrParams{
				Uuid:             sessionLog.Uuid,
				StartMoodRating:  0,
				StartMood:        "N/A",
				FinishMoodRating: 0,
				FinishMood:       "N/A",
				SessionCompleted: 0,
			}
			_, err = q.CreateTummoBreathingMr(ctx, params)

		default:
			return fmt.Errorf("unsupported session type: %s", args.SessionType)
		}

		if err != nil {
			return fmt.Errorf("failed to create session-specific data: %w", err)
		}

		result.UUID = sessionLog.Uuid.String()
		result.UserId = args.UserID
		result.SessionType = args.SessionType
		return nil
	})

	return result, err
}

// update the time for total_time_spent_in_min in users table everytime user finishes or quites a session
// type UpdateMeditationTimeParams struct {
// 	UserId      int64  `json:"user_id"`
// 	SessionType string `json:"session_type"`
// 	UUID        string `json:"uuid"`
// }

// create a user
type CreateUserTransactiontArgs struct {
	Email          string
	FirstName      string
	LastName       string
	Gender         string
	Birthdate      time.Time
	Country        string
	Goal           string
	Platform       string
	HashedPassword string
}

type CreateUserResult struct {
	ID int64
}

// CreateUserTransaction handles the creation of a user transaction based on the platform.
func (store *Store) CreateUserTransaction(ctx context.Context, args CreateUserTransactiontArgs) (CreateUserResult, error) {
	var result CreateUserResult

	err := store.execTx(ctx, func(q *Queries) error {
		existsOnMr := 0
		existsOnMobile := 0

		// Determine if the user is from MR or Mobile
		switch args.Platform {
		case "mobile":
			existsOnMobile = 1
		case "mr":
			existsOnMr = 1
		}

		// Create user parameters
		createUserArgs := CreateUserParams{
			Email:          args.Email,
			FirstName:      args.FirstName,
			LastName:       args.LastName,
			Gender:         args.Gender,
			BirthDate:      args.Birthdate,
			Country:        args.Country,
			Goals:          args.Goal,
			IsMrUser:       int16(existsOnMr),
			IsMobileUser:   int16(existsOnMobile),
			HashedPassword: args.HashedPassword,
		}

		user, err := q.CreateUser(ctx, createUserArgs)
		if err != nil {
			if prErr, ok := err.(*pq.Error); ok {
				fmt.Println(prErr.Code.Name())
				switch prErr.Code.Name() {
				case "unique_violation":
					return fmt.Errorf("unique_violation")
				}
				// log.Println(prErr.Code.Name())
			}
			return fmt.Errorf("error occurred while creating a user: %w", err)
		}

		// now create user in its respective profile table
		switch args.Platform {
		case "mobile":
			err := q.CreateUserProfileMobile(ctx, user.ID)
			if err != nil {
				return fmt.Errorf("error occurred while creating a user mobile profile: %w", err)
			}
		case "mr":
			err := q.CreateUserProfileMr(ctx, user.ID)
			if err != nil {
				return fmt.Errorf("error occurred while creating a user mr profile: %w", err)
			}
		}
		result.ID = user.ID
		return nil
	})

	if err != nil {
		return CreateUserResult{}, err
	}

	return result, nil
}

// create a user
type UpdateSessionFinishTransactionParams struct {
	Uuid             uuid.UUID
	SessionType      string
	FinishMoodRating int16
	FinishMood       string
	SessionCompleted int16
	EndsAt           time.Time
}

// CreateUserTransaction handles the creation of a user transaction based on the platform.
func (store *Store) UpdateSessionFinishTransaction(ctx context.Context, args UpdateSessionFinishTransactionParams) error {

	err := store.execTx(ctx, func(q *Queries) error {
		var sessionStartTime time.Time
		var sessionFinishTime time.Time

		switch args.SessionType {
		case util.TIBETAN_SINGING_BOWL_MR:
			TSBparams := UpdateTibetanSingingBowlMrFinishingMoodByUuidParams{
				Uuid:             args.Uuid,
				FinishMoodRating: args.FinishMoodRating,
				FinishMood:       args.FinishMood,
				SessionCompleted: args.SessionCompleted,
				EndsAt:           args.EndsAt,
			}
			TSBLog, err := q.UpdateTibetanSingingBowlMrFinishingMoodByUuid(ctx, TSBparams)

			if err != nil {
				return err
			}
			if !TSBLog.DeletedAt.IsZero() {
				return fmt.Errorf("can not update a session that has already been deleted. uuid is %v", args.Uuid)
			}
			sessionStartTime = TSBLog.StartedAt
			sessionFinishTime = TSBLog.EndsAt
		case util.TUMMO_BREATHING_MR:
			TBparams := UpdateTummoBreathingMrFinishingMoodByUuidParams{
				Uuid:             args.Uuid,
				FinishMoodRating: args.FinishMoodRating,
				FinishMood:       args.FinishMood,
				SessionCompleted: args.SessionCompleted,
				EndsAt:           args.EndsAt,
			}
			TBLog, err := q.UpdateTummoBreathingMrFinishingMoodByUuid(ctx, TBparams)
			if err != nil {
				return err
			}
			if !TBLog.DeletedAt.IsZero() {
				return fmt.Errorf("can not update a session that has already been deleted. uuid is %v", args.Uuid)
			}
			sessionStartTime = TBLog.StartedAt
			sessionFinishTime = TBLog.EndsAt
		default:
			return fmt.Errorf("unsupported session type")
		}

		// Calculate the difference between the two times
		duration := int64(sessionFinishTime.Sub(sessionStartTime).Minutes())

		platform := util.GetPlatformTypeBasedOnSessionType(args.SessionType)

		//now get the user
		userId, err := q.GetUserIdFromSessionLogUuid(ctx, args.Uuid)

		if err != nil {
			return err
		}

		// create a variable to store the time user has already meditated before
		var existingMeditatedTime int64
		switch platform {
		case util.MR:

			existingMeditatedTime, err = q.GetUserProfileMrTime(ctx, userId)
			if err != nil {
				return err
			}

			updateArgs := UpdateUserProfileMrTimeParams{
				UserID:               userId,
				TotalTimeSpentInMins: existingMeditatedTime + duration,
			}
			err = q.UpdateUserProfileMrTime(ctx, updateArgs)
			if err != nil {
				return err
			}
		case util.MOBILE:
			existingMeditatedTime, err = q.GetUserProfileMobileTime(ctx, userId)
			if err != nil {
				return err
			}
			updateArgs := UpdateUserProfileMobilerTimeParams{
				UserID:               userId,
				TotalTimeSpentInMins: existingMeditatedTime + duration,
			}
			err := q.UpdateUserProfileMobilerTime(ctx, updateArgs)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("platform type is not supported")
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// testing only
type CreateUserForTestingDeletionResult struct {
	UserId                               int64
	TibetanSingingBowlMrSessionLogUserId int64
	TibetanSingingBowlMrSessionLogUuid   uuid.UUID
	TibetanSingingBowlMrUuid             uuid.UUID
	TummoBreathingMrSessionLogUserId     int64
	TummoBreathingMrSessionLogUuid       uuid.UUID
	TummoBreathingMrUuid                 uuid.UUID
}

func (store *Store) CreateUserForTestingDeletion(ctx context.Context) (CreateUserForTestingDeletionResult, error) {
	var result CreateUserForTestingDeletionResult

	err := store.execTx(ctx, func(q *Queries) error {
		// first create a user
		format := "2006-01-02"
		birthDate, err := time.Parse(format, "1990-01-01")

		if err != nil {
			return fmt.Errorf("error occurred while attempting to parse time %v", err)
		}

		createUserArgs := CreateUserParams{
			Email:          util.RandomString(10) + "@example.com", // Generate a random email
			FirstName:      util.RandomString(5),                   // Generate a random first name
			LastName:       util.RandomString(5),                   // Generate a random last name
			Gender:         util.RandomGender(),                    // Randomly choose gender
			BirthDate:      birthDate,                              // Generate a random birth date
			Country:        util.RandomCountryCode(),               // Randomly choose a country code
			HashedPassword: util.RandomString(12),                  // Assuming a random string for hashed password
		}

		createdUser, err := q.CreateUser(ctx, createUserArgs)

		if err != nil {
			return fmt.Errorf("error occurred while attempting to create user %v", err)
		}
		result.UserId = createdUser.ID

		// create session log for tibetan singing bowl based on user id
		createSessionLogTSBArgs := CreateSessionLogParams{
			UserID:          createdUser.ID,
			SessionType:     "tibetan_singing_bowl_mr",
			SessionPlatform: "mr",
		}
		TSBSessionLog, err := q.CreateSessionLog(context.Background(), createSessionLogTSBArgs)
		if err != nil {
			return fmt.Errorf("error occurred while attempting to create session log %v", err)
		}
		result.TibetanSingingBowlMrSessionLogUuid = TSBSessionLog.Uuid
		result.TibetanSingingBowlMrSessionLogUserId = TSBSessionLog.UserID

		// create tibetan singing bowl based on session log uuid

		createTibetanSingingBowlArgs := CreateTibetanSingingBowlMrParams{
			Uuid:             TSBSessionLog.Uuid,
			StartMoodRating:  int16(util.RandomInt(0, 10)),
			StartMood:        util.RandomString(5),
			FinishMoodRating: int16(util.RandomInt(0, 10)),
			FinishMood:       util.RandomString(5),
			SessionCompleted: 0,
		}

		tibetanSingingBowl, err := q.CreateTibetanSingingBowlMr(context.Background(), createTibetanSingingBowlArgs)
		if err != nil {
			return fmt.Errorf("error occurred while attempting to create tibetan singing bowl  %v", err)
		}
		result.TibetanSingingBowlMrUuid = tibetanSingingBowl.Uuid

		// create session log for tummo breathing mr based on user id
		createSessionLogTBArgs := CreateSessionLogParams{
			UserID:          createdUser.ID,
			SessionType:     "tummo_breathing_mr",
			SessionPlatform: "mr",
		}

		TBSessionLog, err := q.CreateSessionLog(context.Background(), createSessionLogTBArgs)

		if err != nil {
			return fmt.Errorf("error occurred while attempting to create session log %v", err)
		}
		result.TummoBreathingMrSessionLogUuid = TBSessionLog.Uuid
		result.TummoBreathingMrSessionLogUserId = TBSessionLog.UserID

		// create tummo breathing based on session log uuid
		createTummoBreathingArgs := CreateTummoBreathingMrParams{
			Uuid:             TBSessionLog.Uuid,
			StartMoodRating:  int16(util.RandomInt(0, 10)),
			StartMood:        util.RandomString(5),
			FinishMoodRating: int16(util.RandomInt(0, 10)),
			FinishMood:       util.RandomString(5),
			SessionCompleted: 0,
		}

		tummoBreathing, err := q.CreateTummoBreathingMr(context.Background(), createTummoBreathingArgs)
		if err != nil {
			return fmt.Errorf("error occurred while attempting to create tibetan singing bowl  %v", err)
		}
		result.TummoBreathingMrUuid = tummoBreathing.Uuid

		return nil
	})

	return result, err
}

type CreateUserEmailTransactiontArgs struct {
	Email string
}

func (store *Store) CreateUserEmailTransaction(ctx context.Context, args CreateUserEmailTransactiontArgs) error {
	err := store.execTx(ctx, func(q *Queries) error {
		params := CreateUserEmailParams{
			UniqueId: uuid.New(),
			Email:    args.Email,
		}
		_, err := q.CreateUserEmail(ctx, params)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (store *Store) CreateChakraTestResult(ctx context.Context, arg CreateChakraTestResultParams) (ChakraTestResult, error) {
	return store.Queries.CreateChakraTestResult(ctx, arg)
}
