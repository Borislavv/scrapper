package livenessinterface

type Service interface {
	IsAlive() bool
}
