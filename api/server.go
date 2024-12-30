package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	db "github.com/lotusMind/meditation/db/sqlc"
)

var (
	InvalidSessionType     = errors.New("the session type is invalid")
	UnsupportedSessionType = errors.New("unsupported session type")
	InvalidUuid            = errors.New("uuid is invalid")
)

// servers all the http requests for lotus mind
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// set up api routes for that server
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	// Route definition for Session
	router.POST("/session/create/:user_id/:session_type", server.createSession)
	router.POST("/session/update/start/:session_uuid/:session_type", server.updateSessionStartingMood)
	router.POST("/session/update/finish/:session_uuid/:session_type", server.updateSessionFinishingMood)
	router.POST("/session/update/quit/:session_uuid/:session_type", server.updateSessionQuit)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
