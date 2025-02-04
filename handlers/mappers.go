package mappers

import (
	"fmt"
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

func PathToStoryIdWithPrefix(path string) domain.ChapterRef {
	prefix := "story"
	pathContent := domain.ChapterRef(initStoryPart)
	if path != fmt.Sprintf("/%s", prefix) && path != fmt.Sprintf("/%s/", prefix) {
		pathContent = domain.ChapterRef(strings.TrimPrefix(strings.TrimSpace(path), fmt.Sprintf("/%s", prefix))[1:])
	}
	return pathContent
}
