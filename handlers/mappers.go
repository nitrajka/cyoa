package mappers

import (
	"strings"

	"chooseYourAdventure/domain"
)

const initStoryPart = "intro"

func PathToSubStoryId(path string) domain.ChapterRef {
	pathContent := domain.ChapterRef(initStoryPart)
	if path != "/" && path != "" {
		// [1:] because we want to remove "/" character
		pathContent = domain.ChapterRef(strings.TrimSpace(path[1:]))
	}
	return pathContent
}
