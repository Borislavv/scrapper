package entityinterface

type Field interface {
	GetMessage() string
	GetPreviousValue() interface{}
	GetCurrentValue() interface{}
}
