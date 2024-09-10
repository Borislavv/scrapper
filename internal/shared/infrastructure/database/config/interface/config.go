package databaseconfiginterface

type Configurator interface {
	GetMongoHost() string
	GetMongoPort() int
	GetMongoLogin() string
	GetMongoPassword() string
	GetMongoDatabase() string
}
