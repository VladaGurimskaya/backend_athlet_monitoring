package service

import (
	api "backend_athlet_monitoring/.api_athlet_monitoring/go"
	"backend_athlet_monitoring/internal/storage"
	"context"
)

type ServiceTeams struct {
	storage *storage.Storage
}

func NewServiceTeams(storage *storage.Storage) (*ServiceTeams, error) {
	return &ServiceTeams{
		storage: storage,
	}, nil
}

func (s *ServiceTeams) TeamsJoinPost(ctx context.Context, request api.TeamsJoinPostRequest) (api.ImplResponse, error) {
	return api.ImplResponse{}, nil
}

func (s *ServiceTeams) TeamsJoinRequestIdApprovePost(ctx context.Context, id int32) (api.ImplResponse, error) {
	return api.ImplResponse{}, nil
}

func (s *ServiceTeams) TeamsJoinRequestIdRejectPost(ctx context.Context, id int32) (api.ImplResponse, error) {
	return api.ImplResponse{}, nil
}
