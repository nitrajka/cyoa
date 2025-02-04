package handlers

import (
	"chooseYourAdventure/domain"
	mappers "chooseYourAdventure/handlers"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strings"
)

var tmpl = template.Must(template.ParseFiles("templates/template.html"))

type storyHandler struct {
	logger             logrus.FieldLogger
	storyTellerService domain.StoryTellerService
	template           *template.Template
	pathParserF        func(string) domain.ChapterRef
}

type HandlerOption func(handler *storyHandler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(handler *storyHandler) {
		handler.template = t
	}
}

func WithPathFunc(f func(urlPath string) domain.ChapterRef) HandlerOption {
	return func(handler *storyHandler) {
		handler.pathParserF = f
	}
}

func NewStoryHandler(log logrus.FieldLogger, storyTellerService domain.StoryTellerService, options ...HandlerOption) http.Handler {
	h := &storyHandler{
		logger:             log,
		storyTellerService: storyTellerService,
		template:           tmpl,
		pathParserF:        mappers.PathToSubStoryId,
	}
	for _, optionF := range options {
		optionF(h)
	}
	return h
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
	if strings.HasPrefix(r.URL.Path, "/public/") {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/public/")
		http.FileServer(http.Dir("./templates/static/")).ServeHTTP(w, r)
		return
	}

	//subStoryId := mappers.PathToSubStoryId(r.URL.Path)
	chapterId := s.pathParserF(r.URL.Path)
	story, err := s.storyTellerService.FetchSubStory(chapterId)
	if err != nil {
		s.logger.WithError(err).WithField("substoryRef", chapterId).Error("failed to fetch sub story")
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}

	if story != nil {
		err := s.template.Execute(w, story)
		if err != nil {
			s.logger.WithError(err).Error("failed to execute template")
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}
