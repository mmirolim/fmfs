package log

import (
	logger "github.com/Sirupsen/logrus"
)

func Start() {
	// set default logrus formatter
	logger.SetFormatter(&logger.TextFormatter{})
	// set default log level
	logger.SetLevel(logger.DebugLevel)
}

// alias for map
type Flds map[string]interface{}

// print with Fields
func PrintWf(scope, action string, flds map[string]interface{}) *logger.Entry {
	// set  extra fields func, msg and err
	flds["action"] = action
	flds["scope"] = scope
	return logger.WithFields(flds)
}

// print without fields
func Print(scope, action string) *logger.Entry {
	flds := make(map[string]interface{})
	// set extra fields
	flds["action"] = action
	flds["scope"] = scope
	return logger.WithFields(flds)
}
