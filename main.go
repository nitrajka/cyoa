package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"chooseYourAdventure/domain"
	handlers "chooseYourAdventure/handlers/http"
	"chooseYourAdventure/repositories"
	"chooseYourAdventure/services"
	"github.com/sirupsen/logrus"
)

const PORT = "8080"

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	var fileName, port string
	flag.StringVar(&port, "p", PORT, "Port to run APP on.")
	flag.StringVar(&fileName, "f", "story.json", "json file name to load story from")
	flag.Parse()

	storyParts, err := loadStoryParts(fileName)
	if err != nil {
		logger.WithError(err).Fatal("failed to load story parts")
	}

	repo := repositories.NewStoryRepository(storyParts)
	service := services.NewStoryTeller(repo)
	handler := handlers.NewStoryHandler(logger, service)

	logger.WithField("port", port).Info("Starting server...")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
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
