package main

import (
	api "backend_athlet_monitoring/.api_athlet_monitoring/go"
	"backend_athlet_monitoring/internal/database"
	"backend_athlet_monitoring/internal/service"
	"backend_athlet_monitoring/internal/storage"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func runMigrations() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL не установлен")
	}

	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("ошибка создания экземпляра migrate: %v", err)
	}
	defer m.Close()

	_, _, err = m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("Ошибка получения версии миграции: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка применения миграций: %v", err)
	}

	return nil
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequestContextInjector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "http-request", r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	if err := runMigrations(); err != nil {
		log.Printf("Ошибка при применении миграций: %v", err)
	} else {
		log.Println("Миграции успешно применены")
	}

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных: ", err)
	}
	defer db.Close()

	store := storage.NewStorage(db)

	if err := store.CreateAdmin(); err != nil {
		log.Printf("Ошибка создания администратора: %v", err)
	} else {
		log.Println("Администратор успешно создан")
	}

	serviceAuth, err := service.NewServiceAuth(store)
	if err != nil {
		log.Fatal("Ошибка сервиса авторизации: ", err)
	}

	serviceBiometric, err := service.NewServiceBiometric(store)
	if err != nil {
		log.Fatal("Ошибка сервиса биометрии: ", err)
	}

	serviceTeams, err := service.NewServiceTeams(store)
	if err != nil {
		log.Fatal("Ошибка сервиса команд: ", err)
	}

	serviceTraining, err := service.NewServiceTraining(store)
	if err != nil {
		log.Fatal("Ошибка сервиса тренировок: ", err)
	}

	serviceAuthController := api.NewAuthAPIController(serviceAuth)
	serviceBiometricController := api.NewBiometricsAPIController(serviceBiometric)
	serviceTeamsController := api.NewTeamsAPIController(serviceTeams)
	serviceTrainingController := api.NewTrainingAPIController(serviceTraining)

	mainRouter := chi.NewRouter()
	mainRouter.Use(CORS)
	mainRouter.Use(RequestContextInjector)

	router := api.NewRouter(serviceAuthController, serviceBiometricController, serviceTeamsController, serviceTrainingController)
	mainRouter.Mount("/", router)

	log.Println("Сервер запущен на порту :8000")
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}
