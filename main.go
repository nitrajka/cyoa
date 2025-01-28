package main

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"

	"chooseYourAdventure/domain"
	handlers "chooseYourAdventure/handlers/http"
	"chooseYourAdventure/repositories"
	"chooseYourAdventure/services"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	var fileName string
	// todo: make port configurable as well
	flag.StringVar(&fileName, "f", "story.json", "json file name to load story from")
	flag.Parse()

	storyParts, err := loadStoryParts(fileName)
	if err != nil {
		logger.WithError(err).Fatal("failed to load story parts")
	}

	// todo: better naming, don't use single char names
	// example: repo, service, handler
	r := repositories.NewStoryRepository(storyParts)
	s := services.NewStoryTeller(r)
	st := handlers.NewStoryHandler(logger, s)

	logger.Info("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", st); err != nil {
		logger.WithError(err).Warn("ListenAndServe failed")
	}
}

func loadStoryParts(fileName string) (map[domain.StoryRef]domain.Story, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var storyData map[domain.StoryRef]domain.Story
	if err = json.Unmarshal(data, &storyData); err != nil {
		return nil, err
	}
	return storyData, nil
}
