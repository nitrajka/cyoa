package domain

type Story struct {
	Title   string        `json:"title"`
	Story   []string      `json:"story"`
	Options []StoryOption `json:"options"`
}

type StoryOption struct {
	Text string `json:"text"`
	Ref  string `json:"arc"`
}

type StoryRef string

type ref StoryRef
