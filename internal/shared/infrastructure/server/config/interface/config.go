package sharedserverconfiginterface

import "time"

type Configurator interface {
	GetServerName() string
	GetServerPort() string
	GetServerShutDownTimeout() time.Duration
	GetServerRequestTimeout() time.Duration
}
