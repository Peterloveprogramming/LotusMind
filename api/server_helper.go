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
	Uuid             uuid.UUID
	sessionType      string
	FinishMoodRating int16
	FinishMood       string
	SessionCompleted int16
	EndsAt           time.Time
}

func updateSessionFinishMood(server *Server, ctx *gin.Context, args updateSessionFinishMoodParams) error {
	switch args.sessionType {
	case "tibetan_singing_bowl_mr":
		TSBparams := db.UpdateTibetanSingingBowlMrFinishingMoodByUuidParams{
			Uuid:             args.Uuid,
			FinishMoodRating: args.FinishMoodRating,
			FinishMood:       args.FinishMood,
			SessionCompleted: args.SessionCompleted,
			EndsAt:           args.EndsAt,
		}
		err := server.store.UpdateTibetanSingingBowlMrFinishingMoodByUuid(ctx, TSBparams)

		if err != nil {
			return err
		}
		return nil
	case "tummo_breathing_mr":
		TBparams := db.UpdateTummoBreathingMrFinishingMoodByUuidParams{
			Uuid:             args.Uuid,
			FinishMoodRating: args.FinishMoodRating,
			FinishMood:       args.FinishMood,
			SessionCompleted: args.SessionCompleted,
			EndsAt:           args.EndsAt,
		}
		err := server.store.UpdateTummoBreathingMrFinishingMoodByUuid(ctx, TBparams)

		if err != nil {
			return err
		}
		return nil
	default:
		return InvalidSessionType
	}
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
