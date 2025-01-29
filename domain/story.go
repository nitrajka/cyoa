package domain

import "fmt"

type Story struct {
	Title   string        `json:"title"`
	Story   []string      `json:"story"`
	Options []StoryOption `json:"options"`
}

type StoryOption struct {
	Text string `json:"text"`
	Ref  string `json:"arc"`
}

func (o StoryOption) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, o.Ref)
}

type StoryRef string
