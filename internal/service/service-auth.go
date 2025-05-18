package service

import (
	api "backend_athlet_monitoring/.api_athlet_monitoring/go"
	"backend_athlet_monitoring/internal/models"
	"backend_athlet_monitoring/internal/storage"
	"backend_athlet_monitoring/internal/utils"
	"context"
	"fmt"
	"time"
)

type ServiceAuth struct {
	storage *storage.Storage
}

func NewServiceAuth(storage *storage.Storage) (*ServiceAuth, error) {
	return &ServiceAuth{
		storage: storage,
	}, nil
}

func (s *ServiceAuth) AuthChangePasswordPost(ctx context.Context, request api.AuthChangePasswordPostRequest) (api.ImplResponse, error) {
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

	if err := s.storage.ChangePassword(accessClaims.UserID, request.OldPassword, request.NewPassword); err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка изменения пароля: %v", err)
	}

	return api.ImplResponse{
		Code: 200,
		Body: "Пароль успешно изменен",
	}, nil
}

func (s *ServiceAuth) AuthLoginPost(ctx context.Context, request api.LoginRequest) (api.ImplResponse, error) {
	user, err := s.storage.Login(request.Email, request.Password)
	if err != nil {
		return api.ImplResponse{
			Code: 401,
		}, fmt.Errorf("Неверный логин или пароль")
	}

	accessToken, refreshToken, err := utils.CreateJWTTokens(user)
	if err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка создания JWT токенов: %v", err)
	}

	headers := map[string][]string{
		"Set-Cookie": {
			fmt.Sprintf("access_token=%s; HttpOnly; Secure; Path=/; Max-Age=900", accessToken),
			fmt.Sprintf("refresh_token=%s; HttpOnly; Secure; Path=/; Max-Age=86400", refreshToken),
		},
	}

	tokens := api.LoginResponse{
		Role: user.Role,
	}

	return api.ImplResponse{
		Code:    200,
		Headers: headers,
		Body:    tokens,
	}, nil
}

func (s *ServiceAuth) AuthRegisterAthletePost(ctx context.Context, request api.AthleteRegisterRequest) (api.ImplResponse, error) {
	passwordHash, err := utils.CreatePasswordHash(request.Password)
	if err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка хеширования пароля: %v", err)
	}

	dateOfBirth, err := utils.ParseDateYYYYMMDD(request.DateOfBirth)
	if err != nil {
		return api.ImplResponse{
			Code: 400,
		}, fmt.Errorf("ошибка в дате рождения: %v", err)
	}

	user := &models.User{
		Email:        request.Email,
		PasswordHash: passwordHash,
		Role:         "athlete",
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	if _, err := s.storage.CreateUser(user); err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка создания пользователя: %v", err)
	}

	athlete := models.Athlete{
		AthleteID:   user.ID,
		Email:       request.Email,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		MiddleName:  request.MiddleName,
		Gender:      request.Gender,
		DateOfBirth: dateOfBirth,
		Phone:       request.Phone,
	}

	if err := s.storage.CreateAthlete(athlete); err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка создания профиля атлета: %v", err)
	}

	return api.ImplResponse{
		Code: 201,
	}, nil
}

func (s *ServiceAuth) AuthRegisterInvitePost(ctx context.Context, request api.InviteRegisterRequest) (api.ImplResponse, error) {
	err := s.storage.CheckInviteCode(request.InviteCode)
	if err != nil {
		return api.ImplResponse{
			Code: 400,
		}, fmt.Errorf("неверный код приглашения: %v", err)
	}

	passwordHash, err := utils.CreatePasswordHash(request.Password)
	if err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка хеширования пароля: %v", err)
	}

	user := &models.User{
		Email:        request.Email,
		PasswordHash: passwordHash,
		Role:         request.Role,
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	userId, err := s.storage.CreateUser(user)
	if err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка создания пользователя: %v", err)
	}

	switch request.Role {
	case "coach":
		coach := models.Coach{
			CoachID:         userId,
			Email:           request.Email,
			Phone:           request.Phone,
			FirstName:       request.FirstName,
			LastName:        request.LastName,
			MiddleName:      request.MiddleName,
			PasswordHash:    passwordHash,
			LicenseNumber:   request.LicenseNumber,
			ExperienceLevel: request.ExperienceLevel,
			SportTypeID:     int(request.SportTypeId),
		}
		err = s.storage.CreateCoach(coach)
		if err != nil {
			return api.ImplResponse{
				Code: 500,
			}, fmt.Errorf("ошибка создания профиля тренера: %v", err)
		}
	case "medical":
		medical := models.MedicalStaff{
			MedicalStaffID: userId,
			Email:          request.Email,
			FirstName:      request.FirstName,
			LastName:       request.LastName,
			MiddleName:     request.MiddleName,
			PasswordHash:   passwordHash,
			Specialization: request.Specialization,
			LicenseNumber:  request.LicenseNumber,
			Verified:       true,
			OrganizationID: int(request.OrganizationId),
		}

		err = s.storage.CreateMedicalStaff(medical)
		if err != nil {
			return api.ImplResponse{
				Code: 500,
			}, fmt.Errorf("ошибка создания профиля медицинского персонала: %v", err)
		}
	}

	err = s.storage.UseInviteCode(request.InviteCode)
	if err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка использования приглашения: %v", err)
	}

	return api.ImplResponse{
		Code: 201,
	}, nil
}

func (s *ServiceAuth) AuthCreateInviteCodePost(ctx context.Context, request api.InviteCreateRequest) (api.ImplResponse, error) {
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

	if accessClaims.Role != "admin" {
		return api.ImplResponse{
			Code: 403,
		}, fmt.Errorf("недостаточно прав для создания приглашения")
	}

	invite := models.CreateInvite{
		Email:          request.Email,
		Role:           request.Role,
		FullName:       request.FullName,
		LicenseNumber:  request.LicenseNumber,
		OrganizationID: int(request.OrganizationId),
	}

	inviteCode := utils.GenerateInviteCode()

	if err := s.storage.CreateInvite(invite, inviteCode, accessClaims.UserID); err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка создания приглашения: %v", err)
	}

	response := api.InviteCreateResponse{
		InviteCode: inviteCode,
		Url:        "http://localhost:5173/register-with-invite/" + inviteCode,
	}

	return api.ImplResponse{
		Code: 201,
		Body: response,
	}, nil
}

func (s *ServiceAuth) AuthGetInviteDetailsPost(ctx context.Context, request api.InviteDetailsRequest) (api.ImplResponse, error) {
	err := s.storage.CheckInviteCode(request.InviteCode)
	if err != nil {
		return api.ImplResponse{
			Code: 400,
		}, fmt.Errorf("неверный код приглашения: %v", err)
	}

	inviteDetails, err := s.storage.GetInviteDetails(request.InviteCode)
	if err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка получения деталей приглашения: %v", err)
	}

	response := api.InviteDetailsResponse{
		InviteCode:     inviteDetails.InviteCode,
		Email:          inviteDetails.Email,
		Role:           inviteDetails.Role,
		LicenseNumber:  inviteDetails.LicenseNumber,
		OrganizationId: int32(inviteDetails.OrganizationID),
	}

	return api.ImplResponse{
		Code: 200,
		Body: response,
	}, nil
}

func (s *ServiceAuth) AuthAthleteProfilePost(ctx context.Context, request api.AthleteProfileRequest) (api.ImplResponse, error) {
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

	userId := int(request.AthleteId)

	if userId == 0 {
		userId = accessClaims.UserID
	}

	profileInfo, err := s.storage.GetAthleteProfile(userId)
	if err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка получения профиля спортсмена: %v", err)
	}

	response := api.AthleteProfileResponse{
		Email:       profileInfo.Email,
		Phone:       profileInfo.Phone,
		FirstName:   profileInfo.FirstName,
		MiddleName:  profileInfo.MiddleName,
		LastName:    profileInfo.LastName,
		DateOfBirth: profileInfo.DateOfBirth.Format("2006-01-02"),
	}

	return api.ImplResponse{
		Code: 200,
		Body: response,
	}, nil
}

func (s *ServiceAuth) AuthAthleteTeamPost(ctx context.Context, request api.AthleteProfileRequest) (api.ImplResponse, error) {
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

	return api.ImplResponse{}, nil
}

func (s *ServiceAuth) AuthGetAllInvitesGet(ctx context.Context) (api.ImplResponse, error) {
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

	if accessClaims.Role != "admin" {
		return api.ImplResponse{
			Code: 403,
		}, fmt.Errorf("недостаточно прав для создания приглашения")
	}

	invites, err := s.storage.GetAllInvites()
	if err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка получения всех приглашений: %v", err)
	}

	response := api.InviteAllResponse{
		Invites: invites,
	}

	return api.ImplResponse{
		Code: 200,
		Body: response,
	}, nil
}

func (s *ServiceAuth) AuthCancelInvitePost(ctx context.Context, request api.InviteCancelRequest) (api.ImplResponse, error) {
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

	if accessClaims.Role != "admin" {
		return api.ImplResponse{
			Code: 403,
		}, fmt.Errorf("недостаточно прав для создания приглашения")
	}

	if err := s.storage.CancelInvite(int(request.InviteId)); err != nil {
		return api.ImplResponse{
			Code: 500,
		}, fmt.Errorf("ошибка отмены приглашения: %v", err)
	}

	return api.ImplResponse{
		Code: 200,
		Body: "Приглашение успешно отменено",
	}, nil
}
