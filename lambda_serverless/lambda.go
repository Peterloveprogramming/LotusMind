package lambdaServerless

import (
	"fmt"

	"github.com/lotusMind/meditation/chakareport"
	db "github.com/lotusMind/meditation/db/sqlc"
	sendemail "github.com/lotusMind/meditation/email"
	"github.com/lotusMind/meditation/storage"
	"github.com/lotusMind/meditation/util"
)

type Lambda struct {
	config             util.Config
	store              db.Store
	chakaraReportMaker chakareport.Maker
	storageMaker       storage.Maker
	sendEmailMaker     sendemail.Maker
}

func NewLambda(config util.Config, store db.Store) (*Lambda, error) {

	chakaraReportMaker, err := chakareport.ChakraMaker(config.APP_ENVIROMENT, config.CHAKARA_REPORT_API_URL)
	if err != nil {
		return nil, fmt.Errorf("can not create chakaraReportMaker: %w", err)
	}

	storageMaker, err := storage.StorageMaker(config.APP_ENVIROMENT, config.AWSRegion, config.AWSAccessKeyID, config.AWSSecretAccessKey, config.AWSBucketName)
	if err != nil {
		return nil, fmt.Errorf("can not create storageMaker: %w", err)
	}

	var frontEndUrl string
	if config.APP_ENVIROMENT == "dev" {
		frontEndUrl = config.FrontEndUrlDev
	} else {
		frontEndUrl = config.FrontEndUrlProd
	}
	sendEmailMaker, err := sendemail.EmailMaker(frontEndUrl, config.Email, config.EmailPassword, config.EmailSmtp, config.EmailSmtpAddress)
	if err != nil {
		return nil, fmt.Errorf("can not create sendEmailMaker: %w", err)
	}
	lambda := &Lambda{
		store:              store,
		config:             config,
		chakaraReportMaker: chakaraReportMaker,
		storageMaker:       storageMaker,
		sendEmailMaker:     sendEmailMaker,
	}
	return lambda, nil
}
