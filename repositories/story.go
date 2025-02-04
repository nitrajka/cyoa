package repositories

import "chooseYourAdventure/domain"

type storyRepository struct {
	storyParts map[domain.ChapterRef]domain.Chapter
}

func NewStoryRepository(parts map[domain.ChapterRef]domain.Chapter) domain.StoryRepository {
	return &storyRepository{storyParts: parts}
}

func (s *storyRepository) GetSubStory(ref domain.ChapterRef) (*domain.Chapter, error) { // todo: remove error
	if storyRepo, ok := s.storyParts[ref]; ok {
		return &storyRepo, nil
	}
	return nil, nil
}
