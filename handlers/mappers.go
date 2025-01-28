package mappers

import "chooseYourAdventure/domain"

const initStoryPart = "intro"

func PathToSubStoryId(path string) domain.StoryRef {
	pathContent := domain.StoryRef(initStoryPart)
	if path != "/" {
		pathContent = domain.StoryRef(path[1:])
	}
	return pathContent
}
