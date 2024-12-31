package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/lotusMind/meditation/db/sqlc"
	"github.com/lotusMind/meditation/util"
)

// createSession
type createSessionRequestParams struct {
	ID          int64  `uri:"user_id"  binding:"required,min=1"`
	SessionType string `uri:"session_type"  binding:"required,session_type"`
}

func (server *Server) createSession(ctx *gin.Context) {
	//verify params
	var req createSessionRequestParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// verify session platform
	sessionPlatform := util.GetPlatformTypeBasedOnSessionType(req.SessionType)
	if len(sessionPlatform) <= 1 {
		ctx.JSON(http.StatusBadRequest, errorResponse(InvalidSessionType))
		return
	}

	createSessionLogArgs := db.CreateSessionLogParams{
		UserID:          req.ID,
		SessionType:     req.SessionType,
		SessionPlatform: sessionPlatform,
	}
	sessionLogResult, err := server.store.CreateSessionLogTransaction(ctx, createSessionLogArgs)

	if err != nil {
		if strings.Contains(err.Error(), "foreign_key_violation") {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("user does not exist in database")))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, sessionLogResult)
}

// updateSessionStartingMood
type updateSessionStartingMoodParams struct {
	SessionUuid string `uri:"session_uuid"  binding:"required,min=1"`
	SessionType string `uri:"session_type"  binding:"required,session_type"`
}
type updateSessionStartingMoodBody struct {
	StartingMoodRating int16  `json:"start_mood_rating" binding:"required,min=1"`
	StartingMood       string `json:"start_mood" binding:"required,min=1,max=20"`
}

func (server *Server) updateSessionStartingMood(ctx *gin.Context) {

	//verify params
	var reqParam updateSessionStartingMoodParams
	if err := ctx.ShouldBindUri(&reqParam); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//verify body
	var reqBody updateSessionStartingMoodBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert uuid from stirng to uuid type
	sessionUuid, err := uuid.Parse(reqParam.SessionUuid)
	if err != nil {
		// Handle the error, e.g., return or log it
		ctx.JSON(http.StatusBadRequest, errorResponse(InvalidUuid))
		return
	}

	updateSessionStartMoodArgs := updateSessionStartMoodParams{
		Uuid:            sessionUuid,
		sessionType:     reqParam.SessionType,
		StartMoodRating: reqBody.StartingMoodRating,
		StartMood:       reqBody.StartingMood,
	}

	err = updateSessionStartMood(server, ctx, updateSessionStartMoodArgs)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// updateSessionFinishingMood
type updateSessionFinishingMoodParams struct {
	SessionUuid string `uri:"session_uuid"  binding:"required,min=1"`
	SessionType string `uri:"session_type"  binding:"required,session_type"`
}
type updateSessionFinishingMoodBody struct {
	FinishMoodRating int16  `json:"finish_mood_rating" binding:"required,min=1"`
	FinishMood       string `json:"finish_mood" binding:"required,min=1,max=20"`
	EndsAt           string `json:"ends_at" binding:"required,min=1`
}

func (server *Server) updateSessionFinishingMood(ctx *gin.Context) {

	//verify params
	var reqParam updateSessionFinishingMoodParams
	if err := ctx.ShouldBindUri(&reqParam); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//verify body
	var reqBody updateSessionFinishingMoodBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert uuid from stirng to uuid type
	sessionUuid, err := uuid.Parse(reqParam.SessionUuid)
	if err != nil {
		// Handle the error, e.g., return or log it
		ctx.JSON(http.StatusBadRequest, errorResponse(InvalidUuid))
		return
	}

	// convert ends at to time.Time format
	time, err := time.Parse(time.RFC3339, reqBody.EndsAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updateSessionFinishMoodArgs := updateSessionFinishMoodParams{
		Uuid:             sessionUuid,
		sessionType:      reqParam.SessionType,
		FinishMoodRating: reqBody.FinishMoodRating,
		FinishMood:       reqBody.FinishMood,
		EndsAt:           time,
		SessionCompleted: 1,
	}

	err = updateSessionFinishMood(server, ctx, updateSessionFinishMoodArgs)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// updateSessionQuit
type sessionQuitParams struct {
	SessionUuid string `uri:"session_uuid"  binding:"required,min=1"`
	SessionType string `uri:"session_type"  binding:"required,session_type"`
}
type sessionQuitBody struct {
	EndsAt string `json:"ends_at" binding:"required,min=1`
}

func (server *Server) updateSessionQuit(ctx *gin.Context) {

	//verify params
	var reqParam sessionQuitParams
	if err := ctx.ShouldBindUri(&reqParam); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//verify body
	var reqBody sessionQuitBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert uuid from stirng to uuid type
	sessionUuid, err := uuid.Parse(reqParam.SessionUuid)
	if err != nil {
		// Handle the error, e.g., return or log it
		ctx.JSON(http.StatusBadRequest, errorResponse(InvalidUuid))
		return
	}

	// convert ends at to time.Time format
	time, err := time.Parse(time.RFC3339, reqBody.EndsAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updateSessionQuitArgs := updateSessionQuitParams{
		Uuid:        sessionUuid,
		sessionType: reqParam.SessionType,
		EndsAt:      time,
	}

	err = updateSessionQuit(server, ctx, updateSessionQuitArgs)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
