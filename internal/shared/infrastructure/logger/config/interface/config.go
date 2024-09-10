package sharedloggerinterface

type Configurator interface {
	GetLoggerLevel() string
	GetLoggerOutput() string
	GetLoggerFormatter() string
	GetLoggerContextExtraFields() []string
}
