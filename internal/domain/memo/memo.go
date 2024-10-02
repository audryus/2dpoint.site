package domain

type Memo struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Url    string `json:"url"`
}
