package root

type APIResponse struct {
	Status int         `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Count  int         `json:"count,omitempty"`
	Type   string      `json:"type,omitempty"`
}
