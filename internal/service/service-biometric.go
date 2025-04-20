package service

import (
	api "backend_athlet_monitoring/.api_athlet_monitoring/go"
	"backend_athlet_monitoring/internal/storage"
	"context"
)

type ServiceBiometric struct {
	storage *storage.Storage
}

func NewServiceBiometric(storage *storage.Storage) (*ServiceBiometric, error) {
	return &ServiceBiometric{
		storage: storage,
	}, nil
}

func (s *ServiceBiometric) BiometricsPost(ctx context.Context, request api.BiometricInput) (api.ImplResponse, error) {
	return api.ImplResponse{}, nil
}
