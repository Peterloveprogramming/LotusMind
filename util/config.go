package util

import (
	"time"

	"github.com/spf13/viper"
)

// stores all the configuration of the application
// values will then be read by viper.
type Config struct {
	APP_ENVIROMENT         string        `mapstructure:"APP_ENVIROMENT"`
	CHAKARA_REPORT_API_URL string        `mapstructure:"CHAKARA_REPORT_API_URL"`
	DBDriver               string        `mapstructure:"DB_DRIVER"`
	DBSource               string        `mapstructure:"DB_SOURCE"`
	ServerAddress          string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	AWSRegion              string        `mapstructure:"AWS_REGION"`
	AWSAccessKeyID         string        `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey     string        `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AWSBucketName          string        `mapstructure:"AWS_BUCKET_NAME"`
	Email                  string        `mapstructure:"EMAIL"`
	EmailPassword          string        `mapstructure:"EMAIL_PASSWORD"`
	EmailSmtp              string        `mapstructure:"EMAIL_SMTP"`
	EmailSmtpAddress       string        `mapstructure:"EMAIL_SMTP_ADDR"`
	FrontEndUrlDev         string        `mapstructure:"FRONTEND_URL_DEV"`
	FrontEndUrlProd        string        `mapstructure:"FRONTEND_URL_PROD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // can be json or xml

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
