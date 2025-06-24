package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	db "github.com/lotusMind/meditation/db/sqlc"
	lambdaServerless "github.com/lotusMind/meditation/lambda_serverless"
	"github.com/lotusMind/meditation/util"
)

var lambdaWrapper *lambdaServerless.Lambda

// how to establish database connection in lambda?
func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch {
	case event.HTTPMethod == "GET" && event.Path == "/test":
		// {
		//   "httpMethod": "GET",
		//   "path": "/test"
		// }
		return lambdaWrapper.Test(ctx, event), nil
	case event.HTTPMethod == "POST" && event.Path == "/register_email":
		return lambdaWrapper.RegisterEmail(ctx, event), nil
	case event.HTTPMethod == "POST" && event.Path == "/chakra/results/getChakraReport":
		return lambdaWrapper.GetChakraReport(ctx, event), nil

	case event.HTTPMethod == "POST" && event.Path == "/chakra/results/":
		// {
		// "httpMethod": "POST",
		// "path": "/chakra/results/",
		// "body": "{\"email\":\"test@example.com\", \"testNum\":\"ABC123XYZ\"}"
		// }
		type ChakraRequest struct {
			Email   string `json:"email"`
			TestNum string `json:"testNum"`
		}
		var req ChakraRequest
		err := json.Unmarshal([]byte(event.Body), &req)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Invalid JSON body",
			}, nil
		}

		if req.Email == "" || req.TestNum == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Missing email or testNum",
			}, nil
		}

		return lambdaWrapper.GetChakraTestResults(ctx, event, req.Email, req.TestNum), nil
	default:
		return lambdaWrapper.RequestNotFound(ctx, event), nil
	}
}

func main() {
	//Load configuration, use "." because main.go is in the same directory
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load configuration")
	}

	// 	// create a new database connection
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect", err)
	}

	store := db.NewStore(conn)
	lambdaWrapper, err = lambdaServerless.NewLambda(config, store)
	if err != nil {
		log.Fatal("can not create server", err)
	}

	lambda.Start(handler)
}
