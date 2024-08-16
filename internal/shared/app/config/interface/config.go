package sharedconfiginterface

type Configurator interface {
	/* logger */
	GetLoggerLevel() string
	GetLoggerOutput() string
	GetLoggerFormatter() string
	GetLoggerContextExtraFields() []string
	/* mongo */
	GetMongoHost() string
	GetMongoPort() int
	GetMongoLogin() string
	GetMongoPassword() string
	GetMongoDatabase() string
}
