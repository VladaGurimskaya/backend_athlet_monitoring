package service

import (
	api "backend_athlet_monitoring/.api_athlet_monitoring/go"
	"backend_athlet_monitoring/internal/storage"
	"backend_athlet_monitoring/internal/utils"
	"context"
	"fmt"
	"log"
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

func (s *ServiceTraining) TrainingPlansPlanIdItemsPost(ctx context.Context, planId int32, request api.TrainingPlanItemCreate) (api.ImplResponse, error) {
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

	if accessClaims.Role != "coach" {
		return api.ImplResponse{
			Code: 403,
		}, fmt.Errorf("недостаточно прав для создания команды")
	}

	return api.ImplResponse{
		Code: 201,
	}, nil
}

func (s *ServiceTraining) TrainingPlansPost(ctx context.Context, request api.TrainingPlanCreate) (api.ImplResponse, error) {
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

	if accessClaims.Role != "coach" {
		return api.ImplResponse{
			Code: 403,
		}, fmt.Errorf("недостаточно прав для создания команды")
	}

	err = s.storage.CreateTrainingPlan(request, accessClaims.UserID)
	if err != nil {
		log.Println("ошибка создания тренировочного плана: ", err)
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка создания тренировочного плана: %v", err)
	}

	return api.ImplResponse{
		Code: 201,
	}, nil
}

func (s *ServiceTraining) GetTrainingPlansPost(ctx context.Context) (api.ImplResponse, error) {
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

	if accessClaims.Role != "coach" {
		return api.ImplResponse{
			Code: 403,
		}, fmt.Errorf("недостаточно прав для создания команды")
	}

	trainingPlanList, err := s.storage.GetTrainingPlans(accessClaims.UserID)
	if err != nil {
		log.Println("ошибка получения тренировочных планов: ", err)
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка получения тренировочных планов: %v", err)
	}

	return api.ImplResponse{
		Code: 200,
		Body: trainingPlanList,
	}, nil
}
