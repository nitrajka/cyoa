package domain

type StoryRepository interface {
	GetSubStory(storyRef StoryRef) (*Story, error)
}
