package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/lotusMind/meditation/db/sqlc"
	"github.com/lotusMind/meditation/util"
)

type createSessionRequest struct {
	ID          int64  `uri:"user_id"  binding:"required,min=1"`
	SessionType string `uri:"session_type"  binding:"required,min=1,max=50"`
}

func (server *Server) createSession(ctx *gin.Context) {
	var req createSessionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	sessionPlatform := util.GetPlatformTypeBasedOnSessionType(req.SessionType)
	if len(sessionPlatform) <= 1 {
		ctx.JSON(http.StatusBadRequest, errorResponse(InvalidSessionType))
	}

	createSessionLogArgs := db.CreateSessionLogParams{
		UserID:          req.ID,
		SessionType:     req.SessionType,
		SessionPlatform: sessionPlatform,
	}
	sessionLogResult, err := server.store.CreateSessionLogTransaction(ctx, createSessionLogArgs)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, sessionLogResult)
}
