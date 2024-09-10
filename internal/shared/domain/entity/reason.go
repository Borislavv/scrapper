package entity

type Reason struct {
	ComparedVersion int `bson:"comparedVersion"`
	// Fields is a map of diff. fields (key is a field name).
	Fields map[string]*Field `bson:"fields"`
}

func NewReason(comparedVersion int, fields map[string]*Field) *Reason {
	return &Reason{ComparedVersion: comparedVersion, Fields: fields}
}

func (r *Reason) GetComparedVersion() int {
	return r.ComparedVersion
}

func (r *Reason) GetFields() map[string]*Field {
	return r.Fields
}
