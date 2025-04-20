package service

import (
	api "backend_athlet_monitoring/.api_athlet_monitoring/go"
	"backend_athlet_monitoring/internal/storage"
	"context"
)

type ServiceTraining struct {
	storage *storage.Storage
}

func NewServiceTraining(storage *storage.Storage) (*ServiceTraining, error) {
	return &ServiceTraining{
		storage: storage,
	}, nil
}

func (s *ServiceTraining) TrainingPlansItemsItemIdDelete(ctx context.Context, itemId int32) (api.ImplResponse, error) {
	return api.ImplResponse{}, nil
}

func (s *ServiceTraining) TrainingPlansItemsItemIdOverridePost(ctx context.Context, itemId int32, request api.TrainingPlansItemsItemIdOverridePostRequest) (api.ImplResponse, error) {
	return api.ImplResponse{}, nil
}

func (s *ServiceTraining) TrainingPlansPlanIdAssignPost(ctx context.Context, planId int32, request api.TrainingPlansPlanIdAssignPostRequest) (api.ImplResponse, error) {
	return api.ImplResponse{}, nil
}

func (s *ServiceTraining) TrainingPlansPlanIdItemsPost(ctx context.Context, planId int32, request api.TrainingPlansPlanIdItemsPostRequest) (api.ImplResponse, error) {
	return api.ImplResponse{}, nil
}

func (s *ServiceTraining) TrainingPlansPost(ctx context.Context, request api.TrainingPlansPostRequest) (api.ImplResponse, error) {
	return api.ImplResponse{}, nil
}
