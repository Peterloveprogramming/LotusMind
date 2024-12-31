package util

import "time"

const (
	//Platform
	MR     = "mr"
	MOBILE = "mobile"

	//Session type
	TIBETAN_SINGING_BOWL_MR = "tibetan_singing_bowl_mr"
	TUMMO_BREATHING_MR      = "tummo_breathing_mr"
)

func IsSupportedPlatform(platform string) bool {
	switch platform {
	case MR, MOBILE:
		return true
	}
	return false
}

func ValidateDateFormat(date string) bool {
	_, err := time.Parse(GetDateFormat(), date)
	return err == nil
}

func IsSessionTypeSupported(sessionType string) bool {
	switch sessionType {
	case TIBETAN_SINGING_BOWL_MR, TUMMO_BREATHING_MR:
		return true
	}
	return false
}
