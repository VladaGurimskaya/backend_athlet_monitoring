package storage

import (
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
