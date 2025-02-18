package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
	"os"

	"chooseYourAdventure/domain"
	mappers "chooseYourAdventure/handlers"
	handlersCLI "chooseYourAdventure/handlers/cli"
	handlersHTTP "chooseYourAdventure/handlers/http"
	"chooseYourAdventure/repositories"
	"chooseYourAdventure/services"
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

func loadStoryParts(fileName string) (map[domain.ChapterRef]domain.Chapter, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// another option how to decode json from file
	//tmp := json.NewDecoder(f)
	//err := tmp.Decode(&domain.Chapter{})
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var storyData map[domain.ChapterRef]domain.Chapter
	if err = json.Unmarshal(data, &storyData); err != nil {
		return nil, err
	}
	return storyData, nil
}

func runInteractiveCLI(logger logrus.FieldLogger, storyParts map[domain.ChapterRef]domain.Chapter) {
	repo := repositories.NewStoryRepository(storyParts)
	service := services.NewStoryTeller(repo)

	handlersCLI.StoryTeller(service)
}

func runHTTP(logger logrus.FieldLogger, port string, storyParts map[domain.ChapterRef]domain.Chapter) {
	repo := repositories.NewStoryRepository(storyParts)
	service := services.NewStoryTeller(repo)
	tpl := template.Must(template.ParseFiles("templates/templateStyled.html"))
	handler := handlersHTTP.NewStoryHandler(logger, service,
		handlersHTTP.WithTemplate(tpl),
		handlersHTTP.WithPathFunc(mappers.PathToStoryIdWithPrefix),
	)

	mux := http.NewServeMux()
	// This story handler is using a custom function and template
	// Because we use /story/ (trailing slash) all web requests
	// whose path has the /story/ prefix will be routed here.
	mux.Handle("/story/", handler)
	// This story handler is using the default functions and templates
	// Because we use / (base path) all incoming requests not
	// mapped elsewhere will be sent here.
	mux.Handle("/", handlersHTTP.NewStoryHandler(logger, service))
	mux.HandleFunc("/health", healthHandler)
	logger.WithField("port", port).Info("Starting server...")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		logger.WithError(err).Warn("ListenAndServe failed")
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
