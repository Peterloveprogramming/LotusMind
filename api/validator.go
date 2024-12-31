package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/lotusMind/meditation/util"
)

var validatePlatform validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if platform, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedPlatform(platform)
	}
	return false
}

var validateDate validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if date, ok := fieldLevel.Field().Interface().(string); ok {
		return util.ValidateDateFormat(date)
	}
	return false
}

var validateSessionType validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if sessionType, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedPlatform(sessionType)
	}
	return false
}
