package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func CreatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CheckPasswordHash(password, passwordHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return fmt.Errorf("неверный пароль")
		}
		return fmt.Errorf("ошибка проверки пароля: %w", err)
	}
	return nil
}

func ParseDateYYYYMMDD(dateStr string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("неверный формат даты (ожидается YYYY-MM-DD): %v", err)
	}
	return date, nil
}

func GenerateInviteCode() string {
	return uuid.New().String()
}

func RandomAnonCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 5)
	_, err := rand.Read(b)
	if err != nil {
		return "ANON0"
	}
	for i := 0; i < 5; i++ {
		b[i] = chars[int(b[i])%len(chars)]
	}
	return string(b)
}

func CalcAge(dob string) int {
	t, err := time.Parse(time.RFC3339, dob)
	if err != nil {
		log.Println("Ошибка парсинга даты рождения:", err)
		return 0
	}
	now := time.Now()
	age := now.Year() - t.Year()
	if now.YearDay() < t.YearDay() {
		age--
	}
	return age
}
