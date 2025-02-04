package domain

type StoryRepository interface {
	GetSubStory(storyRef ChapterRef) (*Chapter, error)
}
