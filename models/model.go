package models

type Question struct {
	Title    string   `json:"title"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Desc     string   `json:"desc"`
	Options  []Option `json:"options,omitempty"`

	Index int `json:"index"`
}

type Option struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Extended bool   `json:"extended"`

	Index int `json:"index"`
}
