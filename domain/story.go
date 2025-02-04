package domain

import "fmt"

type Chapter struct {
	Title      string        `json:"title"`
	Paragraphs []string      `json:"story"`
	Options    []StoryOption `json:"options"`
}

type StoryOption struct {
	Text string `json:"text"`
	Ref  string `json:"arc"`
}

func (o StoryOption) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, o.Ref)
}

type ChapterRef string
