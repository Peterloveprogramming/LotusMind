package api

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/lotusMind/meditation/db/mock"
	"github.com/lotusMind/meditation/util"
)

var testServer *Server
var testStore *mockdb.MockStore
var testConfig util.Config
var testController *gomock.Controller
var err error

func TestMain(m *testing.M) {
	// set it to test mode so it dosent output lots of logs when running package tests.
	gin.SetMode(gin.TestMode)

	//Load configuration
	testConfig, err = util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testController = gomock.NewController(&testing.T{})
	testStore = mockdb.NewMockStore(testController)

	testServer, err = NewServer(testConfig, testStore)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	code := m.Run()
	testController.Finish()
	os.Exit(code)
}
