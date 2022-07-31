package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	addr       string
	logger     Logger
	httpServer *http.Server
}

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
	LogHTTPRequest(r *http.Request, statusCode int, requestDuration time.Duration)
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, id uuid.UUID, event storage.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvent(ctx context.Context, id uuid.UUID) (storage.Event, error)
	SelectOnDayEvents(ctx context.Context, t time.Time) ([]storage.Event, error)
	SelectOnWeekEvents(ctx context.Context, t time.Time) ([]storage.Event, error)
	SelectOnMonthEvents(ctx context.Context, t time.Time) ([]storage.Event, error)
}

func NewRouter(logger Logger, app Application) http.Handler {
	handler := &HTTPHandler{app: app}
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/events/", loggingMiddleware(handler.GetEvents, logger)).Methods("GET")
	router.HandleFunc("/events/", loggingMiddleware(handler.CreateEvent, logger)).Methods("POST")
	router.HandleFunc("/events/{uuid}/", loggingMiddleware(handler.GetEvent, logger)).Methods("GET")
	router.HandleFunc("/events/{uuid}/", loggingMiddleware(handler.UpdateEvent, logger)).Methods("POST")
	router.HandleFunc("/events/{uuid}/", loggingMiddleware(handler.DeleteEvent, logger)).Methods("DELETE")

	return router
}

func NewServer(addr string, logger Logger, app Application) *Server {
	s := &Server{
		addr:   addr,
		logger: logger,
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: NewRouter(logger, app),
	}

	s.httpServer = httpServer

	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("HTTP server started: %s", s.addr))
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("HTTP server stopped: %s", s.addr))

	return s.httpServer.Shutdown(ctx)
}

type HTTPHandler struct {
	app Application
}

func (h *HTTPHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var event storage.Event
	err = json.Unmarshal(data, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.app.CreateEvent(r.Context(), event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *HTTPHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	uuid, err := h.getUUID(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var event storage.Event
	err = json.Unmarshal(data, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.app.UpdateEvent(r.Context(), uuid, event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	uuid, err := h.getUUID(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.app.DeleteEvent(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	uuid, err := h.getUUID(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	event, err := h.app.GetEvent(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *HTTPHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	date, err := time.Parse("2006-01-02 15:04:05", r.URL.Query().Get("date"))
	if err != nil {
		date = time.Now()
	}

	period := r.URL.Query().Get("period")

	var events []storage.Event

	switch period {
	case "day":
		events, err = h.app.SelectOnDayEvents(r.Context(), date)
	case "week":
		events, err = h.app.SelectOnWeekEvents(r.Context(), date)
	case "month":
		events, err = h.app.SelectOnMonthEvents(r.Context(), date)
	default:
		events, err = h.app.SelectOnDayEvents(r.Context(), date)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(&events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (*HTTPHandler) getUUID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(mux.Vars(r)["uuid"])
}
