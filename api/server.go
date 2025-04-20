package api

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lotusMind/meditation/chakareport"
	db "github.com/lotusMind/meditation/db/sqlc"
	"github.com/lotusMind/meditation/storage"
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
	config             util.Config
	store              *db.Store
	router             *gin.Engine
	tokenMaker         token.Maker
	chakaraReportMaker chakareport.Maker
	storageMaker       storage.Maker
}

// set up api routes for that server
func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	chakaraReportMaker, err := chakareport.ChakraMaker(config.APP_ENVIROMENT, config.CHAKARA_REPORT_API_URL)
	if err != nil {
		return nil, fmt.Errorf("can not create chakaraReportMaker: %w", err)
	}

	storageMaker, err := storage.StorageMaker(config.APP_ENVIROMENT, config.AWSRegion, config.AWSAccessKeyID, config.AWSSecretAccessKey, config.AWSBucketName)
	if err != nil {
		return nil, fmt.Errorf("can not create storageMaker: %w", err)
	}

	server := &Server{
		store:              store,
		tokenMaker:         tokenMaker,
		config:             config,
		chakaraReportMaker: chakaraReportMaker,
		storageMaker:       storageMaker,
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

	// 添加 CORS 中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.POST("/user/login", server.loginUser)
	router.POST("/user/create", server.createUser)
	router.POST("/user/register_email", server.registEmail)
	// 添加获取脉轮测试结果的路由
	router.GET("/chakra/results/:email", server.getChakraTestResults)
	router.GET("/chakra/results/:email/:testNum", server.getChakraTestResults)
	router.GET("/chakra/results/getByCode/:code", server.getReportByCode)
	// 添加创建脉轮测试结果的路由
	// router.POST("/chakra/results/create", server.createChakraTestResult)
	router.POST("/chakra/results/create_batch", server.createChakraTestResults)

	router.POST("/chakra/results/getChakraReport", server.getChakraReport)

	router.GET("/test", func(c *gin.Context) {
		err := server.storageMaker.SaveChakaraReportAsText("y.peter2223@gmail.com", "djj", "dj")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	router.GET("/test1", func(c *gin.Context) {
		report, err := server.storageMaker.GetChakaraReportByUniqueCode("y.peter2223@gmail.com", "djj")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"message": report,
		})
	})

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
