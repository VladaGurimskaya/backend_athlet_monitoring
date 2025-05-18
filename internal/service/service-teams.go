package service

import (
	api "backend_athlet_monitoring/.api_athlet_monitoring/go"
	"backend_athlet_monitoring/internal/storage"
	"backend_athlet_monitoring/internal/utils"
	"context"
	"fmt"
	"log"
)

type ServiceTeams struct {
	storage *storage.Storage
}

func NewServiceTeams(storage *storage.Storage) (*ServiceTeams, error) {
	return &ServiceTeams{
		storage: storage,
	}, nil
}

func (s *ServiceTeams) TeamsCreatePost(ctx context.Context, request api.TeamCreateRequest) (api.ImplResponse, error) {
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

	teamId, err := s.storage.CreateTeam(request.TeamName, int(request.SportTypeId))
	if err != nil {
		log.Println("ошибка создания команды: ", err)
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка создания команды: %v", err)
	}

	if err := s.storage.AddCoachToTeam(accessClaims.UserID, int(teamId), "first_coach"); err != nil {
		log.Println("ошибка добавления главного тренера в команду: ", err)
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка добавления главного тренера в команду: %v", err)
	}

	return api.ImplResponse{
		Code: 201,
		Body: "Команда успешно создана",
	}, nil
}

func (s *ServiceTeams) TeamsGet(ctx context.Context) (api.ImplResponse, error) {
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

	if accessClaims.Role == "coach" {
		teams, err := s.storage.GetTeamsByCoach(accessClaims.UserID)
		if err != nil {
			log.Println("ошибка получения команд: ", err)
			return api.ImplResponse{
				Code: 500,
			}, fmt.Errorf("ошибка получения команд: %v", err)
		}

		return api.ImplResponse{
			Code: 200,
			Body: teams,
		}, nil
	}

	return api.ImplResponse{
		Code: 500,
		Body: "ошибка получения команд",
	}, nil
}

func (s *ServiceTeams) TeamsAthletesPost(ctx context.Context, request api.TeamRequest) (api.ImplResponse, error) {
	httpReq, err := utils.ExtractHTTPRequest(ctx)
	if err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка извлечения HTTP-запроса: %v", err)
	}

	_, err = utils.VerifyAccessToken(httpReq)
	if err != nil {
		return api.ImplResponse{
			Code: 401,
		}, fmt.Errorf("ошибка проверки JWT токена: %v", err)
	}

	athletes, err := s.storage.GetAthletesByTeam(int(request.TeamId))
	if err != nil {
		log.Println("ошибка получения участников команды: ", err)
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка получения участников команды: %v", err)
	}

	return api.ImplResponse{
		Code: 200,
		Body: athletes,
	}, nil
}

func (s *ServiceTeams) TeamAthleteRemovePost(ctx context.Context, request api.AthleteProfileRequest) (api.ImplResponse, error) {
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

	err = s.storage.RemoveAthleteFromTeam(int(request.AthleteId))
	if err != nil {
		log.Println("ошибка удаления участника команды: ", err)
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка удаления участника команды: %v", err)
	}

	return api.ImplResponse{
		Code: 200,
		Body: "Участник команды успешно удален",
	}, nil
}

func (s *ServiceTeams) TeamsCoachesPost(ctx context.Context, request api.TeamRequest) (api.ImplResponse, error) {
	httpReq, err := utils.ExtractHTTPRequest(ctx)
	if err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка извлечения HTTP-запроса: %v", err)
	}

	_, err = utils.VerifyAccessToken(httpReq)
	if err != nil {
		return api.ImplResponse{
			Code: 401,
		}, fmt.Errorf("ошибка проверки JWT токена: %v", err)
	}

	coaches, err := s.storage.GetCoachesByTeam(int(request.TeamId))
	if err != nil {
		log.Println("ошибка получения тренеров команды: ", err)
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка получения тренеров команды: %v", err)
	}

	return api.ImplResponse{
		Code: 200,
		Body: coaches,
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
