package root

// NoteSearch struct
type NoteSearch struct {
	NetName    string `json:"net_name,omitempty"`
	NetType    string `json:"net_type,omitempty"`
	DataType   string `json:"data_type,omitempty"`
	SearchText string `json:"search_text,omitempty"`
	SearchType string `json:"search_type,omitempty"`
	Page       int    `json:"page,omitempty"`
	From       int    `json:"from,omitempty"`
	Count      int    `json:"count,omitempty"`
}
