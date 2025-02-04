package handlers

import (
	"chooseYourAdventure/domain"
	mappers "chooseYourAdventure/handlers"
	"fmt"
	"github.com/manifoldco/promptui"
	"strings"
)

func StoryTeller(service domain.StoryTellerService) {
	nextStoryPartId := domain.ChapterRef(mappers.PathToSubStoryId("/"))
	var data *domain.Chapter
	var err error
	for nextStoryPartId != "" {
		data, err = service.FetchSubStory(nextStoryPartId)
		if err != nil {
			fmt.Printf("Error fetching sub story %s\n", nextStoryPartId)
			return
		}
		nextStoryPartId = displayStoryPartAndChooseNextOption(data)
	}
}

func displayStoryPartAndChooseNextOption(storyPart *domain.Chapter) domain.ChapterRef {
	fmt.Println(storyPart.Title + ":     " + strings.Join(storyPart.Paragraphs, "        "))

	templates := &promptui.SelectTemplates{
		Label:    "Select Chapter",
		Active:   "\U0001F336 {{ .Ref | red }}",
		Inactive: "   {{ .Ref | cyan }}",
		// Selected: "\U0001F336 {{ .Text | red | cyan }}",
		Details: `
		--------- Chapter continuation - {{ .Ref }} ----------
		{{ .Text }}`,
	}

	searcher := func(input string, index int) bool {
		pepper := storyPart.Options[index]
		name := strings.Replace(strings.ToLower(pepper.Text), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     storyPart,
		Items:     storyPart.Options,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %w\n", err)
		return ""
	}

	fmt.Printf("You choose number %s\n", result)
	return domain.ChapterRef(result)

}
