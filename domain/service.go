package domain

type StoryTellerService interface {
	FetchSubStory(ChapterRef) (*Chapter, error)
}
