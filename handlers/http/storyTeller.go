package handlers

import (
	"html/template"
	"net/http"

	"chooseYourAdventure/domain"
	mappers "chooseYourAdventure/handlers"
	"github.com/sirupsen/logrus"
)

// todo: make this private, as initializing it would not make sense, since it has private fields
type StoryHandler struct {
	logger             logrus.FieldLogger
	storyTellerService domain.StoryTellerService
}

// todo: return interface instead
func NewStoryHandler(log logrus.FieldLogger, storyTellerService domain.StoryTellerService) *StoryHandler {
	return &StoryHandler{logger: log, storyTellerService: storyTellerService}
}

func (s *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		// todo: should we continue with error?
	}

	tmpl := template.Must(template.ParseFiles("templates/template.html"))
	if story != nil {
		err := tmpl.Execute(w, story)
		if err != nil {
			s.logger.WithError(err).Error("failed to execute template")
		}
	}
}
