package sharedconfiginterface

import "time"

type Configurator interface {
	GetLoggerLevel() string
	GetLoggerOutput() string
	GetLoggerFormatter() string
	GetLoggerContextExtraFields() []string

	GetMongoHost() string
	GetMongoPort() int
	GetMongoLogin() string
	GetMongoPassword() string
	GetMongoDatabase() string

	GetServerName() string
	GetServerPort() string
	GetServerRequestTimeout() time.Duration
	GetServerShutDownTimeout() time.Duration
}
