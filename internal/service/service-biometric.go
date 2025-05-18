package service

import (
	api "backend_athlet_monitoring/.api_athlet_monitoring/go"
	"backend_athlet_monitoring/internal/storage"
	"backend_athlet_monitoring/internal/utils"
	"context"
	"fmt"
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
	httpReq, err := utils.ExtractHTTPRequest(ctx)
	if err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка извлечения HTTP-запроса: %v", err)
	}

	accessClaims, err := utils.VerifyAccessToken(httpReq)
	if err != nil {
		return api.ImplResponse{
			Code: 401,
		}, fmt.Errorf("ошибка проверки JWT токена: %v", err)
	}

	if err := s.storage.CreateBiometricData(accessClaims.UserID, request); err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка создания биометрических данных: %v", err)
	}

	return api.ImplResponse{
		Code: 201,
		Body: "Биометрические данные успешно созданы",
	}, nil
}
