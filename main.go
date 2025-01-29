package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"chooseYourAdventure/domain"
	handlersCLI "chooseYourAdventure/handlers/cli"
	handlersHTTP "chooseYourAdventure/handlers/http"
	"chooseYourAdventure/repositories"
	"chooseYourAdventure/services"
	"github.com/sirupsen/logrus"
)

const PORT = "8080"

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	var fileName, port, appType string
	// port only mandatory when running http app type
	flag.StringVar(&port, "p", PORT, "Port to run APP on. needed only for http mode.")
	flag.StringVar(&fileName, "f", "story.json", "json file name to load story from")
	flag.StringVar(&appType, "t", "http", "how to run this app: cli|http")
	flag.Parse()

	storyParts, err := loadStoryParts(fileName)
	if err != nil {
		logger.WithError(err).Fatal("failed to load story parts")
	}

	switch appType {
	case "http":
		runHTTP(logger, port, storyParts)
	case "cli":
		runInteractiveCLI(logger, storyParts)
	default:
		logger.Fatalf("unknown app type: %s", appType)
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

func runInteractiveCLI(logger logrus.FieldLogger, storyParts map[domain.StoryRef]domain.Story) {
	repo := repositories.NewStoryRepository(storyParts)
	service := services.NewStoryTeller(repo)

	handlersCLI.StoryTeller(service)
}

func runHTTP(logger logrus.FieldLogger, port string, storyParts map[domain.StoryRef]domain.Story) {
	repo := repositories.NewStoryRepository(storyParts)
	service := services.NewStoryTeller(repo)
	handler := handlersHTTP.NewStoryHandler(logger, service)

	logger.WithField("port", port).Info("Starting server...")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
		logger.WithError(err).Warn("ListenAndServe failed")
	}
}
