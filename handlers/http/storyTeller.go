package handlers

import (
	"html/template"
	"net/http"

	"chooseYourAdventure/domain"
	mappers "chooseYourAdventure/handlers"
	"github.com/sirupsen/logrus"
)

type storyHandler struct {
	logger             logrus.FieldLogger
	storyTellerService domain.StoryTellerService
}

func NewStoryHandler(log logrus.FieldLogger, storyTellerService domain.StoryTellerService) http.Handler {
	return &storyHandler{logger: log, storyTellerService: storyTellerService}
}

func (s *storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("Method not allowed"))
		if err != nil {
			s.logger.WithError(err).Error("could not write response")
		}
		return
	}
	subStoryId := mappers.PathToSubStoryId(r.URL.Path)

	story, err := s.storyTellerService.FetchSubStory(subStoryId)
	if err != nil {
		s.logger.WithError(err).WithField("substoryRef", subStoryId).Error("failed to fetch sub story")
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}

	// my own hot reloader
	tmpl := template.Must(template.ParseFiles("templates/template.html"))
	if story != nil {
		err := tmpl.Execute(w, story)
		if err != nil {
			s.logger.WithError(err).Error("failed to execute template")
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}
