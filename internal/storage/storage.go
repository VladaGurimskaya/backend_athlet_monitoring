package storage

import (
	api "backend_athlet_monitoring/.api_athlet_monitoring/go"
	"backend_athlet_monitoring/internal/database"
	"backend_athlet_monitoring/internal/models"
	"backend_athlet_monitoring/internal/utils"
	"fmt"
	"math"
	"time"
)

type Storage struct {
	db *database.Database
}

func NewStorage(db *database.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateAdmin() error {
	passwordHash, err := utils.CreatePasswordHash("admin")
	if err != nil {
		return err
	}

	user := &models.User{
		Email:        "admin@mail.ru",
		PasswordHash: passwordHash,
		Role:         "admin",
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	if _, err := s.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Login(email string, password string) (models.User, error) {
	query := "SELECT user_id, password_hash, role FROM users WHERE email = $1"

	var userId int
	var passwordHash, role string
	err := s.db.QueryRow(query, email).Scan(&userId, &passwordHash, &role)
	if err != nil {
		return models.User{}, fmt.Errorf("ошибка авторизации: %v", err)
	}

	err = utils.CheckPasswordHash(password, passwordHash)
	if err != nil {
		return models.User{}, fmt.Errorf("ошибка авторизации: %v", err)
	}

	user := models.User{
		ID:    userId,
		Email: email,
		Role:  role,
	}

	return user, nil
}

func (s *Storage) CreateUser(user *models.User) (int, error) {
	query := `
		INSERT INTO users (
			email, 
			password_hash, 
			role, 
			is_active, 
			created_at
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING user_id`

	err := s.db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.IsActive,
		user.CreatedAt,
	).Scan(&user.ID)

	if err != nil {
		return 0, fmt.Errorf("ошибка создания пользователя: %v", err)
	}

	return user.ID, nil
}

func (s *Storage) CreateAthlete(athlete models.Athlete) error {
	query := `
		INSERT INTO athletes (
			athlete_id,
			phone,
			first_name,
			middle_name,
			last_name,
			gender,
			date_of_birth
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)`

	_, err := s.db.Exec(
		query,
		athlete.AthleteID,
		athlete.Phone,
		athlete.FirstName,
		athlete.MiddleName,
		athlete.LastName,
		athlete.Gender,
		athlete.DateOfBirth,
	)

	if err != nil {
		return fmt.Errorf("ошибка создания атлета: %v", err)
	}

	return nil
}

func (s *Storage) CreateCoach(coach models.Coach) error {
	query := `
		INSERT INTO coaches (
			coach_id,
			phone,
			first_name,
			middle_name,
			last_name,
			license_id,
			experience_level,
			sport_type_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)`

	_, err := s.db.Exec(
		query,
		coach.CoachID,
		coach.Phone,
		coach.FirstName,
		coach.MiddleName,
		coach.LastName,
		coach.LicenseNumber,
		coach.ExperienceLevel,
		coach.SportTypeID,
	)

	if err != nil {
		return fmt.Errorf("ошибка создания тренера: %v", err)
	}

	return nil
}

func (s *Storage) CreateMedicalStaff(medicalStaff models.MedicalStaff) error {
	query := `
		INSERT INTO medicalstaff (
			medical_staff_id,
			first_name,
			middle_name,
			last_name,
			specialization,
			license_number,
			verified,
			organization_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)`

	_, err := s.db.Exec(
		query,
		medicalStaff.MedicalStaffID,
		medicalStaff.FirstName,
		medicalStaff.MiddleName,
		medicalStaff.LastName,
		medicalStaff.Specialization,
		medicalStaff.LicenseNumber,
		medicalStaff.Verified,
		medicalStaff.OrganizationID,
	)

	if err != nil {
		return fmt.Errorf("ошибка создания медицинского персонала: %v", err)
	}

	return nil
}

func (s *Storage) CreateInvite(invite models.CreateInvite, inviteCode string, userId int) error {
	query := `
		INSERT INTO admininvites (
			email,
			role,
			full_name,
			license_number,
			organization_id,
		    invite_code,
			created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)`

	_, err := s.db.Exec(
		query,
		invite.Email,
		invite.Role,
		invite.FullName,
		invite.LicenseNumber,
		invite.OrganizationID,
		inviteCode,
		userId,
	)

	if err != nil {
		return fmt.Errorf("ошибка создания приглашения: %v", err)
	}

	return nil
}

func (s *Storage) GetInviteDetails(inviteCode string) (models.DetailsInvite, error) {
	invite := models.DetailsInvite{}

	query := `
		SELECT
			invite_code,
			email,
			role,
			license_number,
			organization_id
		FROM admininvites
		WHERE invite_code = $1`

	err := s.db.QueryRow(query, inviteCode).Scan(
		&invite.InviteCode,
		&invite.Email,
		&invite.Role,
		&invite.LicenseNumber,
		&invite.OrganizationID,
	)

	if err != nil {
		return models.DetailsInvite{}, fmt.Errorf("ошибка получения деталей приглашения: %v", err)
	}

	return invite, nil
}

func (s *Storage) CheckInviteCode(inviteCode string) error {
	query := `
		SELECT is_used
		FROM admininvites
		WHERE invite_code = $1`

	var isUsed bool
	err := s.db.QueryRow(query, inviteCode).Scan(&isUsed)

	if err != nil {
		return fmt.Errorf("ошибка проверки приглашения: %v", err)
	}

	if isUsed {
		return fmt.Errorf("приглашение уже использовано")
	}

	return nil
}

func (s *Storage) UseInviteCode(inviteCode string) error {
	query := `
		UPDATE admininvites
		SET is_used = true
		WHERE invite_code = $1`

	_, err := s.db.Exec(query, inviteCode)

	if err != nil {
		return fmt.Errorf("ошибка использования приглашения: %v", err)
	}

	return nil
}

func (s *Storage) ChangePassword(userId int, oldPassword, newPassword string) error {
	query := `SELECT password_hash FROM users WHERE user_id = $1`

	var passwordHash string

	err := s.db.QueryRow(query, userId).Scan(&passwordHash)
	if err != nil {
		return fmt.Errorf("ошибка получения хешированного пароля: %v", err)
	}

	err = utils.CheckPasswordHash(oldPassword, passwordHash)
	if err != nil {
		return fmt.Errorf("неверный пароль: %v", err)
	}

	newPasswordHash, err := utils.CreatePasswordHash(newPassword)
	if err != nil {
		return fmt.Errorf("ошибка хеширования нового пароля: %v", err)
	}

	query = `UPDATE users SET password_hash = $1 WHERE user_id = $2`

	_, err = s.db.Exec(query, newPasswordHash, userId)
	if err != nil {
		return fmt.Errorf("ошибка обновления пароля: %v", err)
	}

	return nil
}

func (s *Storage) GetAthleteProfile(userId int) (models.Athlete, error) {
	var athleteInfo models.Athlete

	query := `
		SELECT
			email
		FROM users
		WHERE user_id = $1`

	err := s.db.QueryRow(query, userId).Scan(&athleteInfo.Email)
	if err != nil {
		return models.Athlete{}, fmt.Errorf("ошибка получения информации о пользователе: %v", err)
	}

	query = `
		SELECT
			phone,
			first_name,
			middle_name,
			last_name,
			date_of_birth
		FROM athletes
		WHERE athlete_id = $1`

	err = s.db.QueryRow(query, userId).Scan(
		&athleteInfo.Phone,
		&athleteInfo.FirstName,
		&athleteInfo.MiddleName,
		&athleteInfo.LastName,
		&athleteInfo.DateOfBirth,
	)
	if err != nil {
		return models.Athlete{}, fmt.Errorf("ошибка получения информации о пользователе: %v", err)
	}

	return athleteInfo, nil
}

func (s *Storage) GetAllInvites() ([]api.InviteDetailsResponse, error) {
	query := `
		SELECT
			invite_id,
			invite_code,
			email,
			role,
			license_number,
			organization_id,
			is_used
		FROM admininvites`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения приглашений: %v", err)
	}
	defer rows.Close()

	var invites []api.InviteDetailsResponse
	for rows.Next() {
		var invite api.InviteDetailsResponse
		if err := rows.Scan(
			&invite.InviteId,
			&invite.InviteCode,
			&invite.Email,
			&invite.Role,
			&invite.LicenseNumber,
			&invite.OrganizationId,
			&invite.IsUsed,
		); err != nil {
			return nil, fmt.Errorf("ошибка получения приглашений: %v", err)
		}
		invites = append(invites, invite)
	}

	return invites, nil
}

func (s *Storage) CancelInvite(inviteId int) error {
	query := `DELETE FROM admininvites WHERE invite_id = $1`

	_, err := s.db.Exec(query, inviteId)
	if err != nil {
		return fmt.Errorf("ошибка отмены приглашения: %v", err)
	}

	return nil
}

func (s *Storage) CreateBiometricData(athleteId int, request api.BiometricInput) error {
	query := `INSERT INTO biometricdata (
		athlete_id,
        date,
    	morning_pulse,
        evening_pulse,
    	hrv,
        weight
	) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(
		query,
		athleteId,
		request.Date,
		request.MorningPulse,
		request.EveningPulse,
		request.HRV,
		request.Weight,
	)
	if err != nil {
		return fmt.Errorf("ошибка создания биометрических данных: %v", err)
	}

	return nil
}

func (s *Storage) CreateTeam(teamName string, sportTypeId int) (int, error) {
	query := `INSERT INTO teams (team_name, sport_type_id) VALUES ($1, $2) RETURNING team_id`

	var teamId int
	err := s.db.QueryRow(query, teamName, sportTypeId).Scan(&teamId)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания команды: %v", err)
	}

	return teamId, nil
}

func (s *Storage) AddCoachToTeam(coachId, teamId int, role string) error {
	query := `INSERT INTO coachteamlink (coach_id, team_id, role_in_team) VALUES ($1, $2, $3)`

	_, err := s.db.Exec(query, coachId, teamId, role)
	if err != nil {
		return fmt.Errorf("ошибка добавления тренера в команду: %v", err)
	}

	return nil
}

func (s *Storage) GetAllAthletes(coachId int) (api.AllAthletesResponse, error) {
	query := `
		SELECT
			a.athlete_id,
			a.first_name,
			a.middle_name,
			a.last_name,
			a.date_of_birth,
			u.email,
			t.team_name,
			tjr.reviewed_at
		FROM athletes a
		JOIN users u ON u.user_id = a.athlete_id
		JOIN teams t ON t.team_id = a.team_id
		JOIN teamjoinrequests tjr ON tjr.athlete_id = a.athlete_id
		JOIN coachteamlink ctl ON ctl.team_id = t.team_id
		WHERE ctl.coach_id = $1
		ORDER BY tjr.reviewed_at DESC`

	rows, err := s.db.Query(query, coachId)
	if err != nil {
		return api.AllAthletesResponse{}, fmt.Errorf("ошибка получения участников: %v", err)
	}
	defer rows.Close()

	var athletes []api.AthleteProfileResponse
	for rows.Next() {
		var athlete api.AthleteProfileResponse
		if err := rows.Scan(
			&athlete.AthleteId,
			&athlete.FirstName,
			&athlete.MiddleName,
			&athlete.LastName,
			&athlete.DateOfBirth,
			&athlete.Email,
			&athlete.TeamName,
			&athlete.TeamSignedDate,
		); err != nil {
			return api.AllAthletesResponse{}, fmt.Errorf("ошибка получения участников: %v", err)
		}

		athletes = append(athletes, athlete)
	}

	return api.AllAthletesResponse{
		Athletes: athletes,
	}, nil
}

func (s *Storage) GetAthleteTeamStatus(athleteId int) (api.TeamAthleteStatusResponse, error) {
	query := `
		SELECT
			tjr.request_id,
			t.team_name,
			tjr.requested_at,
			tjr.status
		FROM teamjoinrequests tjr
		JOIN teams t ON t.team_id = tjr.team_id
		WHERE tjr.athlete_id = $1`

	var status api.TeamAthleteStatusResponse
	err := s.db.QueryRow(query, athleteId).Scan(
		&status.RequestId,
		&status.TeamName,
		&status.RequestedAt,
		&status.Status,
	)
	if err != nil {
		return api.TeamAthleteStatusResponse{}, fmt.Errorf("ошибка получения команды: %v", err)
	}

	return status, nil
}

func (s *Storage) GetTeamByAthlete(athleteId int) (api.AthleteTeamResponse, error) {
	query := `
		SELECT
			t.team_name,
			s.name AS sport_type
		FROM athletes a
		JOIN teams t ON a.team_id = t.team_id
		JOIN sporttypes s ON t.sport_type_id = s.sport_type_id
		WHERE a.athlete_id = $1`

	var teamInfo api.AthleteTeamResponse
	err := s.db.QueryRow(query, athleteId).Scan(
		&teamInfo.TeamName,
		&teamInfo.SportType,
	)
	if err != nil {
		return api.AthleteTeamResponse{}, fmt.Errorf("ошибка получения команды: %v", err)
	}

	return teamInfo, nil
}

func (s *Storage) GetTeamsByCoach(coachId int) (api.TeamResponse, error) {
	query := `
		SELECT
			t.team_id,
			t.team_name,
			s.name AS sport_type,
			COUNT(a.athlete_id) AS count_athletes
		FROM teams t
		JOIN coachteamlink ctl ON t.team_id = ctl.team_id
		JOIN sporttypes s ON t.sport_type_id = s.sport_type_id
		LEFT JOIN athletes a ON a.team_id = t.team_id
		WHERE ctl.coach_id = $1
		GROUP BY t.team_id, t.team_name, s.name
	`

	rows, err := s.db.Query(query, coachId)
	if err != nil {
		return api.TeamResponse{}, fmt.Errorf("ошибка получения команд: %v", err)
	}
	defer rows.Close()

	var teamList []api.Team
	for rows.Next() {
		var team api.Team
		if err := rows.Scan(
			&team.TeamId,
			&team.TeamName,
			&team.SportType,
			&team.CountAthletes,
		); err != nil {
			return api.TeamResponse{}, fmt.Errorf("ошибка сканирования команды: %v", err)
		}
		teamList = append(teamList, team)
	}

	return api.TeamResponse{Teams: teamList}, nil
}

func (s *Storage) GetAthletesByTeam(teamId int) (api.TeamGetAthletesResponse, error) {
	query := `
		SELECT
			a.athlete_id,
			a.first_name,
			a.last_name,
			a.middle_name,
			a.date_of_birth,
			a.phone,
			u.email
		FROM athletes a
		JOIN users u ON u.user_id = a.athlete_id
		WHERE a.team_id = $1
	`

	rows, err := s.db.Query(query, teamId)
	if err != nil {
		return api.TeamGetAthletesResponse{}, fmt.Errorf("ошибка получения атлетов: %v", err)
	}
	defer rows.Close()

	var athletes []api.AthleteProfileResponse
	for rows.Next() {
		var athlete api.AthleteProfileResponse
		if err := rows.Scan(
			&athlete.AthleteId,
			&athlete.FirstName,
			&athlete.LastName,
			&athlete.MiddleName,
			&athlete.DateOfBirth,
			&athlete.Phone,
			&athlete.Email,
		); err != nil {
			return api.TeamGetAthletesResponse{}, fmt.Errorf("ошибка сканирования атлета: %v", err)
		}

		athletes = append(athletes, athlete)
	}

	return api.TeamGetAthletesResponse{Athletes: athletes}, nil
}

func (s *Storage) RemoveAthleteFromTeam(athleteId, teamId int) error {
	query := `UPDATE athletes SET team_id = NULL WHERE athlete_id = $1`

	_, err := s.db.Exec(query, athleteId)
	if err != nil {
		return fmt.Errorf("ошибка удаления атлета из команды: %v", err)
	}

	query = `DELETE FROM teamjoinrequests WHERE athlete_id = $1 AND team_id = $2`
	_, err = s.db.Exec(query, athleteId, teamId)
	if err != nil {
		return fmt.Errorf("ошибка удаления запроса на вступление в команду: %v", err)
	}

	return nil
}

func (s *Storage) GetCoachesByTeam(teamId int) (api.TeamGetCoachesResponse, error) {
	query := `
		SELECT
			c.coach_id,
			c.first_name,
			c.last_name,
			c.middle_name,
			c.phone,
			u.email,
			ctl.role_in_team
		FROM coaches c
		JOIN users u ON u.user_id = c.coach_id
		JOIN coachteamlink ctl ON c.coach_id = ctl.coach_id
		WHERE ctl.team_id = $1
	`

	rows, err := s.db.Query(query, teamId)
	if err != nil {
		return api.TeamGetCoachesResponse{}, fmt.Errorf("ошибка получения тренеров: %v", err)
	}
	defer rows.Close()

	var coaches []api.CoachProfileResponse
	for rows.Next() {
		var coach api.CoachProfileResponse
		if err := rows.Scan(
			&coach.CoachId,
			&coach.FirstName,
			&coach.LastName,
			&coach.MiddleName,
			&coach.Phone,
			&coach.Email,
			&coach.RoleInTeam,
		); err != nil {
			return api.TeamGetCoachesResponse{}, fmt.Errorf("ошибка сканирования тренера: %v", err)
		}

		coaches = append(coaches, coach)
	}

	return api.TeamGetCoachesResponse{Coaches: coaches}, nil
}

func (s *Storage) GetAllTeams() (api.TeamResponse, error) {
	query := `SELECT team_id, team_name FROM teams`

	rows, err := s.db.Query(query)
	if err != nil {
		return api.TeamResponse{}, fmt.Errorf("ошибка получения команд: %v", err)
	}
	defer rows.Close()

	var teamList []api.Team
	for rows.Next() {
		var team api.Team
		if err := rows.Scan(
			&team.TeamId,
			&team.TeamName,
		); err != nil {
			return api.TeamResponse{}, fmt.Errorf("ошибка сканирования команды: %v", err)
		}
		teamList = append(teamList, team)
	}

	return api.TeamResponse{Teams: teamList}, nil
}

func (s *Storage) CreateJoinTeamRequest(teamId, athleteId int) error {
	query := `INSERT INTO teamjoinrequests (team_id, athlete_id) VALUES ($1, $2)`
	_, err := s.db.Exec(query, teamId, athleteId)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса на вступление в команду: %v", err)
	}

	return nil
}

func (s *Storage) GetTeamsJoinRequests(coachId int) (api.TeamJoinsListResponse, error) {
	query := `
		SELECT
			tjr.request_id,
			a.last_name,
			a.first_name,
			a.middle_name,
			u.email,
			t.team_name,
			tjr.requested_at,
			tjr.status
		FROM teamjoinrequests tjr
		JOIN athletes a ON a.athlete_id = tjr.athlete_id
		JOIN teams t ON t.team_id = tjr.team_id
		JOIN users u ON u.user_id = a.athlete_id
		JOIN coachteamlink ctl ON ctl.team_id = tjr.team_id
		WHERE ctl.coach_id = $1`

	rows, err := s.db.Query(query, coachId)
	if err != nil {
		return api.TeamJoinsListResponse{}, fmt.Errorf("ошибка получения запросов на вступление в команду: %v", err)
	}
	defer rows.Close()

	var requests []api.TeamJoin
	for rows.Next() {
		var request api.TeamJoin
		if err := rows.Scan(
			&request.RequestId,
			&request.SecondName,
			&request.FirstName,
			&request.MiddleName,
			&request.Email,
			&request.TeamName,
			&request.RequestDate,
			&request.Status,
		); err != nil {
			return api.TeamJoinsListResponse{}, fmt.Errorf("ошибка сканирования запроса на вступление в команду: %v", err)
		}
		requests = append(requests, request)
	}

	return api.TeamJoinsListResponse{Joins: requests}, nil
}

func (s *Storage) GetTeamsJoinRequestsByTeamId(teamId int) (api.TeamJoinsListResponse, error) {
	query := `
		SELECT
			tjr.request_id,
			a.last_name,
			a.first_name,
			a.middle_name,
			u.email,
			t.team_name,
			tjr.requested_at,
			tjr.status
		FROM teamjoinrequests tjr
		JOIN athletes a ON a.athlete_id = tjr.athlete_id
		JOIN teams t ON t.team_id = tjr.team_id
		JOIN users u ON u.user_id = a.athlete_id
		WHERE tjr.team_id = $1 AND tjr.status = 'pending'`

	rows, err := s.db.Query(query, teamId)
	if err != nil {
		return api.TeamJoinsListResponse{}, fmt.Errorf("ошибка получения запросов на вступление в команду: %v", err)
	}
	defer rows.Close()

	var requests []api.TeamJoin
	for rows.Next() {
		var request api.TeamJoin
		if err := rows.Scan(
			&request.RequestId,
			&request.SecondName,
			&request.FirstName,
			&request.MiddleName,
			&request.Email,
			&request.TeamName,
			&request.RequestDate,
			&request.Status,
		); err != nil {
			return api.TeamJoinsListResponse{}, fmt.Errorf("ошибка сканирования запроса на вступление в команду: %v", err)
		}
		requests = append(requests, request)
	}

	return api.TeamJoinsListResponse{Joins: requests}, nil
}

func (s *Storage) ApproveJoinTeamRequest(requestId, coachId int) error {
	query := `UPDATE teamjoinrequests SET status = 'approved', reviewed_by = $1, reviewed_at = CURRENT_TIMESTAMP WHERE request_id = $2`

	_, err := s.db.Exec(query, coachId, requestId)
	if err != nil {
		return fmt.Errorf("ошибка подтверждения запроса на вступление в команду: %v", err)
	}

	var teamId, athleteId int

	query = `SELECT team_id, athlete_id FROM teamjoinrequests WHERE request_id = $1`
	err = s.db.QueryRow(query, requestId).Scan(&teamId, &athleteId)
	if err != nil {
		return fmt.Errorf("ошибка получения данных запроса на вступление в команду: %v", err)
	}

	query = `UPDATE athletes SET team_id = $1 WHERE athlete_id = $2`
	_, err = s.db.Exec(query, teamId, athleteId)
	if err != nil {
		return fmt.Errorf("ошибка обновления данных атлета: %v", err)
	}

	return nil
}

func (s *Storage) RejectJoinTeamRequest(requestId int) error {
	query := `UPDATE teamjoinrequests SET status = 'rejected' WHERE request_id = $1`

	_, err := s.db.Exec(query, requestId)
	if err != nil {
		return fmt.Errorf("ошибка отклонения запроса на вступление в команду: %v", err)
	}

	return nil
}

func (s *Storage) CreateTrainingPlan(request api.TrainingPlanCreate, coachId int) error {
	query := `INSERT INTO trainingplan (title, team_id, description, start_date, end_date, status, period_type_id, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := s.db.Exec(query, request.Title, request.TeamId, request.Description, request.StartDate, request.EndDate, request.Status, request.PeriodTypeId, coachId)
	if err != nil {
		return fmt.Errorf("ошибка создания тренировочного плана: %v", err)
	}

	return nil
}

func (s *Storage) GetTrainingPlans(coachId int) (api.TrainingPlanListResponse, error) {
	query := `
		SELECT
			tp.plan_id,
			tp.title,
			t.team_name,
			tp.start_date,
			tp.end_date,
			tp.status,
			tp.description
		FROM trainingplan tp
		JOIN teams t ON t.team_id = tp.team_id
		WHERE tp.created_by = $1`

	rows, err := s.db.Query(query, coachId)
	if err != nil {
		return api.TrainingPlanListResponse{}, fmt.Errorf("ошибка получения тренировочных планов: %v", err)
	}
	defer rows.Close()

	var plans []api.TrainingPlan
	for rows.Next() {
		var plan api.TrainingPlan
		if err := rows.Scan(
			&plan.PlanId,
			&plan.Title,
			&plan.TeamName,
			&plan.StartDate,
			&plan.EndDate,
			&plan.Status,
			&plan.Description,
		); err != nil {
			return api.TrainingPlanListResponse{}, fmt.Errorf("ошибка сканирования тренировочного плана: %v", err)
		}
		plans = append(plans, plan)
	}

	return api.TrainingPlanListResponse{Plans: plans}, nil
}

func (s *Storage) GetCriticalAthletes(coachId int) (api.CriticalAthletesResponse, error) {
	query := `
		SELECT
		  a.athlete_id,
		  a.first_name,
		  a.middle_name,
		  a.last_name,
		  t.team_name,
		  pt.name AS period_type
		FROM athletes a
		JOIN teams t ON a.team_id = t.team_id
		JOIN trainingplan tp ON a.team_id = tp.team_id
		JOIN coachteamlink ctl ON ctl.coach_id = tp.created_by AND ctl.team_id = t.team_id
		JOIN periodtypes pt ON pt.period_type_id = tp.period_type_id
		WHERE ctl.coach_id = $1
	`
	rows, err := s.db.Query(query, coachId)
	if err != nil {
		return api.CriticalAthletesResponse{}, fmt.Errorf("ошибка получения критических атлетов: %v", err)
	}
	defer rows.Close()

	var result []api.CriticalAthleteProfileResponse

	for rows.Next() {
		var athlete api.CriticalAthleteProfileResponse
		if err := rows.Scan(
			&athlete.AthleteId,
			&athlete.FirstName,
			&athlete.MiddleName,
			&athlete.LastName,
			&athlete.TeamName,
			&athlete.PeriodType,
		); err != nil {
			return api.CriticalAthletesResponse{}, fmt.Errorf("ошибка сканирования критического атлета: %v", err)
		}

		// Подгружаем биометрию за последние 7 дней
		bioQuery := `
			SELECT
				date,
				hrv,
				morning_pulse,
				evening_pulse,
				weight
			FROM biometricdata
			WHERE athlete_id = $1
			  AND date >= CURRENT_DATE - INTERVAL '6 day'
			ORDER BY date DESC
		`
		bioRows, err := s.db.Query(bioQuery, athlete.AthleteId)
		if err != nil {
			return api.CriticalAthletesResponse{}, fmt.Errorf("ошибка получения биометрии: %v", err)
		}

		var indicators []api.AthleteIndicatorResponse
		for bioRows.Next() {
			var indicator api.AthleteIndicatorResponse
			if err := bioRows.Scan(
				&indicator.Date,
				&indicator.Hrv,
				&indicator.MorningPulse,
				&indicator.EveningPulse,
				&indicator.Weight,
			); err != nil {
				bioRows.Close()
				return api.CriticalAthletesResponse{}, fmt.Errorf("ошибка сканирования биометрии: %v", err)
			}
			indicators = append(indicators, indicator)
		}
		bioRows.Close()
		if len(indicators) == 0 {
			continue // Нет данных — пропускаем
		}

		// baseline (скользящее среднее за 7 дней)
		var sumHRV, sumPulse, sumWeight float64
		var baseHRV, basePulse, baseWeight float64
		var cntHRV, cntPulse, cntWeight int
		for _, ind := range indicators {
			if ind.Hrv > 0 {
				sumHRV += float64(ind.Hrv)
				cntHRV++
			}
			if ind.MorningPulse > 0 {
				sumPulse += float64(ind.MorningPulse)
				cntPulse++
			}
			if ind.Weight > 0 {
				sumWeight += float64(ind.Weight)
				cntWeight++
			}
		}
		if cntHRV > 0 {
			baseHRV = sumHRV / float64(cntHRV)
		}
		if cntPulse > 0 {
			basePulse = sumPulse / float64(cntPulse)
		}
		if cntWeight > 0 {
			baseWeight = sumWeight / float64(cntWeight)
		}

		// Берём последний день (он первый в indicators, так как ORDER BY DESC)
		today := indicators[0]

		// Определяем коэффициенты для периода
		sensitivityCoefHRV := getSensitivityCoefHRV(athlete.PeriodType)
		sensitivityCoefPulse := getSensitivityCoefPulse(athlete.PeriodType)
		sensitivityCoefWeight := getSensitivityCoefWeight(athlete.PeriodType)
		sensitivityCoefCombo1 := getSensitivityCoefCombo1(athlete.PeriodType)
		sensitivityCoefCombo2 := getSensitivityCoefCombo2(athlete.PeriodType)
		sensitivityCoefCombo3 := getSensitivityCoefCombo3(athlete.PeriodType)

		// Считаем дельты
		var (
			deltaHRV    float64
			deltaPulse  float64
			deltaWeight float64
		)
		if baseHRV > 0 {
			deltaHRV = (baseHRV - float64(today.Hrv)) / baseHRV * 100
		}
		if basePulse > 0 {
			deltaPulse = (float64(today.MorningPulse) - basePulse) / basePulse * 100
		}
		if baseWeight > 0 {
			deltaWeight = float64((baseWeight - float64(today.Weight)) / baseWeight * 100)
		}

		// Триггеры одиночные
		triggerHRV := deltaHRV >= (15.0 * sensitivityCoefHRV)
		triggerPulseAlert := deltaPulse >= (15.0 * sensitivityCoefPulse)
		triggerEveningPulse := today.EveningPulse-today.MorningPulse >= 10
		triggerWeight := deltaWeight >= (2.0 * sensitivityCoefWeight)

		// Комбо-триггеры
		triggerCombo1 := deltaHRV >= (15.0*sensitivityCoefCombo1) && deltaPulse >= (10.0*sensitivityCoefCombo1)   // HRV↓ + Pulse↑
		triggerCombo2 := deltaHRV >= (15.0*sensitivityCoefCombo2) && deltaWeight >= (2.0*sensitivityCoefCombo2)   // HRV↓ + Вес↓
		triggerCombo3 := deltaPulse >= (10.0*sensitivityCoefCombo3) && deltaWeight >= (2.0*sensitivityCoefCombo3) // Pulse↑ + Вес↓

		isCritical :=
			triggerHRV ||
				triggerPulseAlert ||
				triggerEveningPulse ||
				triggerWeight ||
				triggerCombo1 ||
				triggerCombo2 ||
				triggerCombo3

		limits := map[string]api.IndicatorLimit{
			"hrv": {
				Min: float32(math.Round(baseHRV*(1-0.15*sensitivityCoefHRV)*10) / 10),
				Max: float32(math.Round(baseHRV*(1+0.15*sensitivityCoefHRV)*10) / 10),
			},
			"morning_pulse": {
				Min: float32(math.Round(basePulse*(1-0.15*sensitivityCoefPulse)*10) / 10),
				Max: float32(math.Round(basePulse*(1+0.15*sensitivityCoefPulse)*10) / 10),
			},
			"evening_pulse": {
				Min: float32(math.Round(basePulse*10) / 10),
				Max: float32(math.Round((basePulse+10)*10) / 10),
			},
			"weight": {
				Min: float32(math.Round(baseWeight*(1-0.02*sensitivityCoefWeight)*10) / 10),
				Max: float32(math.Round(baseWeight*(1+0.02*sensitivityCoefWeight)*10) / 10),
			},
		}
		athlete.Limits = limits

		// Если критично — добавляем
		if isCritical {
			athlete.Indicators = indicators
			athlete.Limits = limits
			result = append(result, athlete)
		}
	}

	return api.CriticalAthletesResponse{Athletes: result}, nil
}

// Индивидуальные sensitivityCoef для разных показателей и комбо
func getSensitivityCoefHRV(period string) float64 {
	switch period {
	case "Восстановительный", "отпуск", "учеба":
		return 1.2
	case "Соревновательный":
		return 1.0
	case "Предсезонка", "сбор", "базовая подготовка":
		return 0.8
	case "Болезнь/травма":
		return 1.3
	default:
		return 1.0
	}
}
func getSensitivityCoefPulse(period string) float64 {
	switch period {
	case "Восстановительный", "отпуск", "учеба":
		return 1.3
	case "Соревновательный":
		return 0.9
	case "Тренировочный цикл":
		return 1.0
	case "Болезнь/травма":
		return 1.4
	default:
		return 1.0
	}
}
func getSensitivityCoefWeight(period string) float64 {
	switch period {
	case "Восстановительный", "отпуск", "учеба":
		return 1.2
	case "Соревновательный":
		return 0.8
	case "Болезнь/травма":
		return 1.5
	default:
		return 1.0
	}
}
func getSensitivityCoefCombo1(period string) float64 {
	switch period {
	case "Восстановительный", "отпуск", "учеба":
		return 1.4
	case "Соревновательный":
		return 1.0
	case "Предсезонка", "сбор", "базовая подготовка":
		return 1.2
	case "Болезнь/травма":
		return 1.6
	default:
		return 1.0
	}
}
func getSensitivityCoefCombo2(period string) float64 {
	switch period {
	case "Восстановительный", "отпуск", "учеба":
		return 1.3
	case "Соревновательный":
		return 1.0
	case "Предсезонка", "сбор", "базовая подготовка":
		return 1.2
	case "Болезнь/травма":
		return 1.5
	default:
		return 1.0
	}
}
func getSensitivityCoefCombo3(period string) float64 {
	// В формулировке не было точных коэффициентов, используем такие же, как для Pulse и Weight или свои
	switch period {
	case "Восстановительный", "отпуск", "учеба":
		return 1.2
	case "Соревновательный":
		return 1.0
	case "Болезнь/травма":
		return 1.5
	default:
		return 1.0
	}
}

func (s *Storage) ReferAthleteToMedicalstaff(athleteId, medicalStaffId, coachId int) error {
	query := "INSERT INTO medicalassignments (athlete_id, assigned_by, medical_staff_id) VALUES ($1, $2, $3)"

	_, err := s.db.Exec(query, athleteId, coachId, medicalStaffId)
	if err != nil {
		return fmt.Errorf("ошибка при назначении атлета на врача: %v", err)
	}

	return nil
}

func (s *Storage) GetMedicalstaff() (api.TeamGetMedicalstaffResponse, error) {
	query := "SELECT medical_staff_id, first_name, middle_name, last_name, specialization FROM medicalstaff"

	rows, err := s.db.Query(query)
	if err != nil {
		return api.TeamGetMedicalstaffResponse{}, fmt.Errorf("ошибка при получении списка врачей: %v", err)
	}
	defer rows.Close()

	var medicalstaff []api.MedicalstaffProfileResponse
	for rows.Next() {
		var ms api.MedicalstaffProfileResponse
		if err := rows.Scan(&ms.MedicalstaffId, &ms.FirstName, &ms.MiddleName, &ms.LastName, &ms.Specialization); err != nil {
			return api.TeamGetMedicalstaffResponse{}, fmt.Errorf("ошибка при чтении данных врача: %v", err)
		}
		medicalstaff = append(medicalstaff, ms)
	}

	return api.TeamGetMedicalstaffResponse{Medicalstaff: medicalstaff}, nil
}

func (s *Storage) GetMedicalAssigments(medicalId int) (api.TeamGetMedicalAssignmentsResponse, error) {
	query := `
		SELECT
			ma.athlete_id,
			ma.assigned_by,
			a.date_of_birth,
			st.name,
			ma.status
		FROM
			medicalassignments ma
		JOIN athletes a ON ma.athlete_id = a.athlete_id
		JOIN sporttypes st ON a.sport_type_id = st.sport_type_id
		WHERE
			ma.medical_staff_id = $1
	`

	rows, err := s.db.Query(query, medicalId)
	if err != nil {
		return api.TeamGetMedicalAssignmentsResponse{}, fmt.Errorf("ошибка при получении списка атлетов: %v", err)
	}
	defer rows.Close()

	var athletes []api.MedicalAssigment
	for rows.Next() {
		var a api.MedicalAssigment
		var dob string
		if err := rows.Scan(&a.AthleteId, &a.AssignedBy, &dob, &a.SportType, &a.Status); err != nil {
			return api.TeamGetMedicalAssignmentsResponse{}, fmt.Errorf("ошибка при чтении данных атлета: %v", err)
		}
		a.Age = int32(utils.CalcAge(dob))
		a.AthleteCode = utils.RandomAnonCode()
		athletes = append(athletes, a)
	}

	return api.TeamGetMedicalAssignmentsResponse{Athletes: athletes}, nil
}

func (s *Storage) GetCriticalAthlete(athleteId int) (api.CriticalAthleteResponse, error) {
	bioQuery := `
		SELECT
			pt.name AS period_type,
			bd.date,
			bd.hrv,
			bd.morning_pulse,
			bd.evening_pulse,
			bd.weight
		FROM biometricdata bd
		JOIN athletes a ON bd.athlete_id = a.athlete_id
		JOIN teams t ON a.team_id = t.team_id
		JOIN trainingplan tp ON t.team_id = tp.team_id
		JOIN periodtypes pt ON pt.period_type_id = tp.period_type_id
		WHERE bd.athlete_id = $1
			AND bd.date >= CURRENT_DATE - INTERVAL '6 day'
		ORDER BY bd.date DESC
	`
	bioRows, err := s.db.Query(bioQuery, athleteId)
	if err != nil {
		return api.CriticalAthleteResponse{}, fmt.Errorf("ошибка получения биометрии: %v", err)
	}
	defer bioRows.Close()

	var (
		periodType string
		indicators []api.AthleteIndicatorResponse
	)
	for bioRows.Next() {
		var indicator api.AthleteIndicatorResponse
		// Сначала получаем период, затем показатели
		if err := bioRows.Scan(
			&periodType,
			&indicator.Date,
			&indicator.Hrv,
			&indicator.MorningPulse,
			&indicator.EveningPulse,
			&indicator.Weight,
		); err != nil {
			return api.CriticalAthleteResponse{}, fmt.Errorf("ошибка сканирования биометрии: %v", err)
		}
		indicators = append(indicators, indicator)
	}
	if len(indicators) == 0 {
		return api.CriticalAthleteResponse{}, nil // Нет данных
	}

	// baseline (скользящее среднее за 7 дней)
	var sumHRV, sumPulse, sumWeight float64
	var baseHRV, basePulse, baseWeight float64
	var cntHRV, cntPulse, cntWeight int
	for _, ind := range indicators {
		if ind.Hrv > 0 {
			sumHRV += float64(ind.Hrv)
			cntHRV++
		}
		if ind.MorningPulse > 0 {
			sumPulse += float64(ind.MorningPulse)
			cntPulse++
		}
		if ind.Weight > 0 {
			sumWeight += float64(ind.Weight)
			cntWeight++
		}
	}
	if cntHRV > 0 {
		baseHRV = sumHRV / float64(cntHRV)
	}
	if cntPulse > 0 {
		basePulse = sumPulse / float64(cntPulse)
	}
	if cntWeight > 0 {
		baseWeight = sumWeight / float64(cntWeight)
	}

	// Берём последний день (он первый в indicators, так как ORDER BY DESC)
	today := indicators[0]

	// Определяем коэффициенты для периода
	sensitivityCoefHRV := getSensitivityCoefHRV(periodType)
	sensitivityCoefPulse := getSensitivityCoefPulse(periodType)
	sensitivityCoefWeight := getSensitivityCoefWeight(periodType)
	sensitivityCoefCombo1 := getSensitivityCoefCombo1(periodType)
	sensitivityCoefCombo2 := getSensitivityCoefCombo2(periodType)
	sensitivityCoefCombo3 := getSensitivityCoefCombo3(periodType)

	// Считаем дельты
	var (
		deltaHRV    float64
		deltaPulse  float64
		deltaWeight float64
	)
	if baseHRV > 0 {
		deltaHRV = (baseHRV - float64(today.Hrv)) / baseHRV * 100
	}
	if basePulse > 0 {
		deltaPulse = (float64(today.MorningPulse) - basePulse) / basePulse * 100
	}
	if baseWeight > 0 {
		deltaWeight = (baseWeight - float64(today.Weight)) / baseWeight * 100
	}

	// Триггеры одиночные
	triggerHRV := deltaHRV >= (15.0 * sensitivityCoefHRV)
	triggerPulseAlert := deltaPulse >= (15.0 * sensitivityCoefPulse)
	triggerEveningPulse := today.EveningPulse-today.MorningPulse >= 10
	triggerWeight := deltaWeight >= (2.0 * sensitivityCoefWeight)

	// Комбо-триггеры
	triggerCombo1 := deltaHRV >= (15.0*sensitivityCoefCombo1) && deltaPulse >= (10.0*sensitivityCoefCombo1)
	triggerCombo2 := deltaHRV >= (15.0*sensitivityCoefCombo2) && deltaWeight >= (2.0*sensitivityCoefCombo2)
	triggerCombo3 := deltaPulse >= (10.0*sensitivityCoefCombo3) && deltaWeight >= (2.0*sensitivityCoefCombo3)

	_ = triggerHRV || triggerPulseAlert || triggerEveningPulse ||
		triggerWeight || triggerCombo1 || triggerCombo2 || triggerCombo3

	limits := map[string]api.IndicatorLimit{
		"hrv": {
			Min: float32(math.Round(baseHRV*(1-0.15*sensitivityCoefHRV)*10) / 10),
			Max: float32(math.Round(baseHRV*(1+0.15*sensitivityCoefHRV)*10) / 10),
		},
		"morning_pulse": {
			Min: float32(math.Round(basePulse*(1-0.15*sensitivityCoefPulse)*10) / 10),
			Max: float32(math.Round(basePulse*(1+0.15*sensitivityCoefPulse)*10) / 10),
		},
		"evening_pulse": {
			Min: float32(math.Round(basePulse*10) / 10),
			Max: float32(math.Round((basePulse+10)*10) / 10),
		},
		"weight": {
			Min: float32(math.Round(baseWeight*(1-0.02*sensitivityCoefWeight)*10) / 10),
			Max: float32(math.Round(baseWeight*(1+0.02*sensitivityCoefWeight)*10) / 10),
		},
	}

	// Можно вернуть только если критично, либо всегда (тогда просто не пиши isCritical проверку)
	return api.CriticalAthleteResponse{
		Indicators: indicators,
		Limits:     limits,
		PeriodType: periodType,
	}, nil
}
