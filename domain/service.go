package domain

type StoryTellerService interface {
	FetchSubStory(StoryRef) (*Story, error)
}
