//go:generate protoc EventService.proto --go_out=./pb/ --go-grpc_out=./pb/ --proto_path=../../../api

package internalgrpc

import (
	context "context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedEventServiceServer
	addr       string
	app        Application
	logger     Logger
	grpcServer *grpc.Server
}

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
	LogGRPCRequest(r interface{}, method string, requestDuration time.Duration)
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

func NewServer(addr string, logger Logger, app Application) *Server {
	s := &Server{
		addr:   addr,
		app:    app,
		logger: logger,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			loggingMiddleware(logger),
		),
	)

	s.grpcServer = grpcServer

	pb.RegisterEventServiceServer(s.grpcServer, s)

	return s
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("GRPC server started: %s", s.addr))

	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() error {
	s.logger.Info(fmt.Sprintf("GRPC server stopped: %s", s.addr))

	s.grpcServer.GracefulStop()

	return nil
}

func (s *Server) Create(ctx context.Context, r *pb.CreateEventRequest) (*pb.ResultResponse, error) {
	storageEvent, err := s.validateAndCreateStorageEvent(r.GetEvent())
	if err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}

	if err = s.app.CreateEvent(ctx, storageEvent); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}

	return &pb.ResultResponse{}, nil
}

func (s *Server) Update(ctx context.Context, r *pb.UpdateEventRequest) (*pb.ResultResponse, error) {
	id, err := uuid.Parse(r.GetId())
	if err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}

	storageEvent, err := s.validateAndCreateStorageEvent(r.GetEvent())
	if err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}

	if err = s.app.UpdateEvent(ctx, id, storageEvent); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}

	return &pb.ResultResponse{}, nil
}

func (s *Server) Delete(ctx context.Context, r *pb.DeleteEventRequest) (*pb.ResultResponse, error) {
	id, err := uuid.Parse(r.GetId())
	if err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}

	if err = s.app.DeleteEvent(ctx, id); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}

	return &pb.ResultResponse{}, nil
}

func (s *Server) Get(ctx context.Context, r *pb.GetEventRequest) (*pb.Event, error) {
	id, err := uuid.Parse(r.GetId())
	if err != nil {
		return &pb.Event{}, err
	}

	event, err := s.app.GetEvent(ctx, id)
	if err != nil {
		return &pb.Event{}, err
	}

	return createEvent(event), nil
}

func (s *Server) GetOnDayEvents(ctx context.Context, r *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	t := r.GetDate().AsTime()

	events, err := s.app.SelectOnDayEvents(ctx, t)
	if err != nil {
		return nil, err
	}

	return createGetEventsResponse(events), nil
}

func (s *Server) GetOnWeekEvents(ctx context.Context, r *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	t := r.GetDate().AsTime()

	events, err := s.app.SelectOnWeekEvents(ctx, t)
	if err != nil {
		return nil, err
	}

	return createGetEventsResponse(events), nil
}

func (s *Server) GetOnMonthEvents(ctx context.Context, r *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	t := r.GetDate().AsTime()

	events, err := s.app.SelectOnMonthEvents(ctx, t)
	if err != nil {
		return nil, err
	}

	return createGetEventsResponse(events), nil
}

func (*Server) validateAndCreateStorageEvent(event *pb.Event) (storage.Event, error) {
	id, err := uuid.Parse(event.GetId())
	if err != nil {
		return storage.Event{}, err
	}

	userID, err := uuid.Parse(event.GetUserId())
	if err != nil {
		return storage.Event{}, err
	}

	return storage.Event{
		ID:          id,
		Title:       event.GetTitle(),
		BeginAt:     event.GetBeginAt().AsTime(),
		EndAt:       event.GetEndAt().AsTime(),
		Description: event.GetDescription(),
		UserID:      userID,
		NotifyAt:    event.GetNotifyAt().AsTime(),
	}, nil
}

func createGetEventsResponse(events []storage.Event) *pb.GetEventsResponse {
	pbEvents := make([]*pb.Event, 0)

	for _, event := range events {
		pbEvents = append(pbEvents, createEvent(event))
	}

	return &pb.GetEventsResponse{
		Events: pbEvents,
	}
}

func createEvent(event storage.Event) *pb.Event {
	return &pb.Event{
		Id:          event.ID.String(),
		Title:       event.Title,
		BeginAt:     timestamppb.New(event.BeginAt),
		EndAt:       timestamppb.New(event.EndAt),
		Description: event.Description,
		UserId:      event.UserID.String(),
		NotifyAt:    timestamppb.New(event.NotifyAt),
	}
}
