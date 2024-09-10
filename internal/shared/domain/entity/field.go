package entity

type Field struct {
	Message       string      `bson:"message"`
	PreviousValue interface{} `bson:"previousValue"`
	CurrentValue  interface{} `bson:"currentValue"`
}

func NewReasonField(message string, previous, current interface{}) *Field {
	return &Field{
		Message:       message,
		PreviousValue: previous,
		CurrentValue:  current,
	}
}

func (f *Field) GetMessage() string {
	return f.Message
}

func (f *Field) GetPreviousValue() interface{} {
	return f.PreviousValue
}

func (f *Field) GetCurrentValue() interface{} {
	return f.CurrentValue
}
