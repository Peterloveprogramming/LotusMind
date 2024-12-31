package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/lotusMind/meditation/db/sqlc"
)

// updateSessionStartMood
type updateSessionStartMoodParams struct {
	Uuid            uuid.UUID
	sessionType     string
	StartMoodRating int16
	StartMood       string
}

func updateSessionStartMood(server *Server, ctx *gin.Context, args updateSessionStartMoodParams) error {
	switch args.sessionType {
	case "tibetan_singing_bowl_mr":
		TSBparams := db.UpdateTibetanSingingBowlMrStartingMoodByUuidParams{
			Uuid:            args.Uuid,
			StartMoodRating: args.StartMoodRating,
			StartMood:       args.StartMood,
		}
		err := server.store.UpdateTibetanSingingBowlMrStartingMoodByUuid(ctx, TSBparams)

		if err != nil {
			return err
		}
		return nil
	case "tummo_breathing_mr":
		TBparams := db.UpdateTummoBreathingMrStartingMoodByUuidParams{
			Uuid:            args.Uuid,
			StartMoodRating: args.StartMoodRating,
			StartMood:       args.StartMood,
		}
		err := server.store.UpdateTummoBreathingMrStartingMoodByUuid(ctx, TBparams)

		if err != nil {
			return err
		}
		return nil
	default:
		return InvalidSessionType
	}
}

// updateSessionFinishMood
type updateSessionFinishMoodParams struct {
	Uuid             string
	Status           string
	sessionType      string
	FinishMoodRating int16
	FinishMood       string
	SessionCompleted int16
	EndsAt           string
}

func updateSessionFinishMood(server *Server, ctx *gin.Context, args updateSessionFinishMoodParams) error {

	// convert uuid from stirng to uuid type
	sessionUuid, err := uuid.Parse(args.Uuid)
	if err != nil {
		// Handle the error, e.g., return or log it
		return err
	}

	// // convert ends at to time.Time format
	time, err := time.Parse(time.RFC3339, args.EndsAt)
	if err != nil {
		return err
	}

	updateSessionFinishTransactionArgs := db.UpdateSessionFinishTransactionParams{
		Uuid:             sessionUuid,
		SessionType:      args.sessionType,
		FinishMoodRating: args.FinishMoodRating,
		FinishMood:       args.FinishMood,
		EndsAt:           time,
		SessionCompleted: args.SessionCompleted,
	}

	err = server.store.UpdateSessionFinishTransaction(ctx, updateSessionFinishTransactionArgs)

	if err != nil {
		return err
	}

	return nil
}

// updateSessionQuit
type updateSessionQuitParams struct {
	Uuid        uuid.UUID
	sessionType string
	EndsAt      time.Time
}

func updateSessionQuit(server *Server, ctx *gin.Context, args updateSessionQuitParams) error {
	switch args.sessionType {
	case "tibetan_singing_bowl_mr":
		TSBparams := db.UpdateTibetanSingingBowlMrQuitByUuidParams{
			Uuid:   args.Uuid,
			EndsAt: args.EndsAt,
		}
		err := server.store.UpdateTibetanSingingBowlMrQuitByUuid(ctx, TSBparams)

		if err != nil {
			return err
		}
		return nil
	case "tummo_breathing_mr":
		TBparams := db.UpdateTummoBreathingMrQuitByUuidParams{
			Uuid:   args.Uuid,
			EndsAt: args.EndsAt,
		}
		err := server.store.UpdateTummoBreathingMrQuitByUuid(ctx, TBparams)

		if err != nil {
			return err
		}
		return nil
	default:
		return InvalidSessionType
	}
}

// updatefinishmood
// server use updatefinishmood, pass in the same parameters for the two functions

// updatefinishmood and updatefinishquit both need to do further work to build the arguments
//updatefinishmood add stage only, status = "finish"
//updatefinishquit need to add missing params, and status = "quit"

// helper checks if its quit or finish.

// in both cases it needs to convert the format.

// then it calls the transaction, passes in completed arguments, also notify the session type.

// transaction first needs to first retrieve the session based on the session type.
// if no session, throw error
// if no start time throw error
// store the start time
// update session finish mood
// then get the profiel session
// calculate
// if result is not ok, roll back and throw error.
//update theo ther function too.
