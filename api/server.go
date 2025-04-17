package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lotusMind/meditation/chakareport"
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
	config             util.Config
	store              *db.Store
	router             *gin.Engine
	tokenMaker         token.Maker
	chakaraReportMaker chakareport.Maker
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
	server := &Server{
		store:              store,
		tokenMaker:         tokenMaker,
		config:             config,
		chakaraReportMaker: chakaraReportMaker,
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
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(server.config.AWSAccessKeyID, server.config.AWSSecretAccessKey, "")),
		)

		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}
		// // Create an S3 client
		s3Client := s3.NewFromConfig(cfg)

		// // Use the S3 client to perform operations
		// // For example, list buckets

		bucketName := "chakara-report"
		// List all buckets
		result, err := s3Client.ListBuckets(c, &s3.ListBucketsInput{})
		if err != nil {
			log.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Check if the "report" bucket exists
		bucketExists := false
		for _, bucket := range result.Buckets {
			if aws.ToString(bucket.Name) == bucketName {
				bucketExists = true
				break
			}
		}

		if bucketExists {
			// user folder - pretend user id is 2
			userFolder := "4/"

			// List objects with the folder prefix
			result, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
				Bucket: &bucketName,
				Prefix: &userFolder,
			})
			if err != nil {
				log.Fatalf("Error listing objects: %v", err)
			}
			// Check if  thefolder eixsts
			folderExists := len(result.Contents) > 0
			fmt.Printf("Folder '%s' exists: %t\n", userFolder, folderExists)

			if !folderExists {
				// Create the folder by creating an empty object with the folder prefix as the key.
				_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
					Bucket: &bucketName,
					Key:    aws.String(userFolder), // Note: No filename here, just the folder path
				})
				if err != nil {
					log.Printf("Couldn't create folder '%s' in bucket '%s'. Here's why: %v\n", userFolder, bucketName, err)
					c.JSON(http.StatusInternalServerError, errorResponse(err))
					return
				}
				log.Printf("Created folder '%s' in bucket '%s'\n", userFolder, bucketName)
			}

			// Generate a UUID
			uniqueID := uuid.New()
			reportName := fmt.Sprintf("%s.txt", uniqueID.String()) // Unique filename
			reportContent := "Your Root Chakra is balanced, grounding you in stability and security, while your Sacral Chakra shows signs of creative energy but mild emotional blockage. The Solar Plexus Chakra is underactive, suggesting a need for increased confidence and self-discipline. Meanwhile, your Heart and Throat Chakras are open, allowing compassion and honest communication to flow naturally."

			// Upload the file to the bucket
			_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket: &bucketName,
				Key:    aws.String(userFolder + reportName), // Full path including folder
				Body:   strings.NewReader(reportContent),
			})
			if err != nil {
				log.Printf("Couldn't upload file '%s' to bucket '%s'. Here's why: %v\n", reportName, bucketName, err)
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			log.Printf("Uploaded file '%s' to bucket '%s'\n", reportName, bucketName)

			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Bucket '%s' exists.", bucketName)})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Bucket '%s' does not exist.", bucketName)})
		}

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
