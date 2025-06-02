package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AccessToken struct {
	UserID int
	Email  string
	Role   string
	jwt.RegisteredClaims
}

type RefreshToken struct {
	UserID int
	Email  string
	Role   string
	jwt.RegisteredClaims
}

type User struct {
	ID           int
	Email        string
	PasswordHash string
	Role         string
	IsActive     bool
	CreatedAt    time.Time
}

type Athlete struct {
	AthleteID    int
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	MiddleName   string
	Gender       string
	DateOfBirth  time.Time
	Phone        string
}

type Coach struct {
	CoachID         int
	Email           string
	Phone           string
	FirstName       string
	LastName        string
	MiddleName      string
	PasswordHash    string
	LicenseNumber   string
	ExperienceLevel string
	SportTypeID     int
}

type MedicalStaff struct {
	MedicalStaffID int
	Email          string
	FirstName      string
	LastName       string
	MiddleName     string
	PasswordHash   string
	Specialization string
	LicenseNumber  string
	Verified       bool
	OrganizationID int
}

type CreateInvite struct {
	Email          string
	Role           string
	FullName       string
	LicenseNumber  string
	OrganizationID int
}

type DetailsInvite struct {
	InviteCode     string
	Email          string
	Role           string
	LicenseNumber  string
	OrganizationID int
}

type Team struct {
	TeamID   int
	TeamName string
}
