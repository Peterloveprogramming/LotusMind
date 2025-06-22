// package main

// import (
// 	"database/sql"
// 	"log"

// 	_ "github.com/lib/pq"
// 	"github.com/lotusMind/meditation/api"
// 	db "github.com/lotusMind/meditation/db/sqlc"
// 	"github.com/lotusMind/meditation/util"
// )

// func main() {
// 	//Load configuration, use "." because main.go is in the same directory
// 	config, err := util.LoadConfig(".")
// 	if err != nil {
// 		log.Fatal("can not load configuration")
// 	}

// 	// create a new database connection
// 	conn, err := sql.Open(config.DBDriver, config.DBSource)

// if err != nil {
// 	log.Fatal("Cannot connect", err)
// }

//		store := db.NewStore(conn)
//		server, err := api.NewServer(config, store)
//		if err != nil {
//			log.Fatal("can not create server", err)
//		}
//		err = server.Start(config.ServerAddress)
//		if err != nil {
//			log.Fatal("can not start the server", err)
//		}
//	}
package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	db "github.com/lotusMind/meditation/db/sqlc"
	lambdaServerless "github.com/lotusMind/meditation/lambda_serverless"
	"github.com/lotusMind/meditation/util"
)

var lambdaWrapper *lambdaServerless.Lambda

// how to establish database connection in lambda?
func handler(ctx context.Context, event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	switch {
	case event.HTTPMethod == "GET" && event.Path == "/test":
		return lambdaWrapper.Test(ctx, event)
	case event.HTTPMethod == "POST" && event.Path == "/register_email":
		return lambdaWrapper.RegisterEmail(ctx, event)
	case event.HTTPMethod == "GET" && event.Path == "/chakra/results/:email":
		return lambdaWrapper.RegisterEmail(ctx, event)
	default:
		return lambdaWrapper.RequestNotFound(ctx, event)
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
