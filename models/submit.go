package models

type SubmitData struct {
	Url  string `json:"url,omitempty"`
	Data []any  `json:"data,omitempty"`
	Num  int    `json:"num,omitempty"`
	Type string `json:"type,omitempty"`
}
