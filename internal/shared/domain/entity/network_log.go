package entity

type NetworkLog struct {
	URL             string
	RequestHeaders  map[string]interface{}
	ResponseHeaders map[string]interface{}
	StatusCode      int
}
