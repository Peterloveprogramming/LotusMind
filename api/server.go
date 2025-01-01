package api

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/lotusMind/meditation/db/sqlc"
	"github.com/lotusMind/meditation/token"
	"github.com/lotusMind/meditation/util"
)

var (
	InvalidSessionType     = errors.New("the session type is invalid")
	UnsupportedSessionType = errors.New("unsupported session type")
	InvalidUuid            = errors.New("uuid is invalid")
	DatabaseIntegrityError = errors.New("database integrity has been violated")
)

// servers all the http requests for lotus mind
type Server struct {
	config     util.Config
	store      *db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// set up api routes for that server
func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	//register validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("platform", validatePlatform)
		v.RegisterValidation("date", validateDate)
		v.RegisterValidation("session_type", validateSessionType)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/user/login", server.loginUser)
	router.POST("/user/create", server.createUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// Route definition for User
	authRoutes.GET("/user/get_info/:id", server.fetchUserInfoById)
	authRoutes.GET("/user/get_time/:platform/:id", server.fetchUserTime)

	// Route definition for Session
	authRoutes.POST("/session/create/:user_id/:session_type", server.createSession)
	authRoutes.POST("/session/update/start/:session_uuid/:session_type", server.updateSessionStartingMood)
	authRoutes.POST("/session/update/finish/:session_uuid/:session_type", server.updateSessionFinishingMood)
	authRoutes.POST("/session/update/quit/:session_uuid/:session_type", server.updateSessionQuit)

	server.router = router

}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
