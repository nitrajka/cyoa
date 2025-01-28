package repositories

import "chooseYourAdventure/domain"

type storyRepository struct {
	storyParts map[domain.StoryRef]domain.Story
}

func NewStoryRepository(parts map[domain.StoryRef]domain.Story) domain.StoryRepository {
	return &storyRepository{storyParts: parts}
}

func (s *storyRepository) GetSubStory(ref domain.StoryRef) (*domain.Story, error) { // todo: remove error
	if storyRepo, ok := s.storyParts[ref]; ok {
		return &storyRepo, nil
	}
	return nil, nil
}
