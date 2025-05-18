package storage

import (
	api "backend_athlet_monitoring/.api_athlet_monitoring/go"
	"backend_athlet_monitoring/internal/database"
	"backend_athlet_monitoring/internal/models"
	"backend_athlet_monitoring/internal/utils"
	"fmt"
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

func (s *Storage) RemoveAthleteFromTeam(athleteId int) error {
	query := `UPDATE athletes SET team_id = NULL WHERE athlete_id = $1`

	_, err := s.db.Exec(query, athleteId)
	if err != nil {
		return fmt.Errorf("ошибка удаления атлета из команды: %v", err)
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
