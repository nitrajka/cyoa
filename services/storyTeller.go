package services

import "chooseYourAdventure/domain"

type storyTeller struct {
	storyRepository domain.StoryRepository
}

func NewStoryTeller(repository domain.StoryRepository) domain.StoryTellerService {
	return &storyTeller{storyRepository: repository}
}

func (st *storyTeller) FetchSubStory(ref domain.ChapterRef) (*domain.Chapter, error) {
	return st.storyRepository.GetSubStory(ref)
}
