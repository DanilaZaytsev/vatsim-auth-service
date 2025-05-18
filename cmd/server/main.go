package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"vatsim-auth-service/internal/db"
	"vatsim-auth-service/internal/handler"
	"vatsim-auth-service/internal/middleware"
)

func main() {
	// Загружаем переменные окружения из .env
	_ = godotenv.Load()

	ctx := context.Background()
	if err := db.InitYDB(ctx); err != nil {
		log.Panicf("Failed to connect to YDB: %v", err)
	}

	// Корневой роутер
	r := mux.NewRouter()

	// Глобальный логгер
	r.Use(middleware.LoggerMiddleware)

	// Публичные ручки
	r.HandleFunc("/auth/vatsim/login", handler.VatsimLoginHandler).Methods("GET")
	r.HandleFunc("/auth/vatsim/callback", handler.VatsimCallbackHandler).Methods("GET")
	r.HandleFunc("/health", handler.HealthHandler).Methods("GET")
	r.HandleFunc("/ready", handler.ReadyHandler).Methods("GET")

	//Защищённые ручки (JWT или cookie)
	private := r.PathPrefix("/").Subrouter()
	private.Use(middleware.AuthMiddleware)
	private.HandleFunc("/me", handler.MeHandler).Methods("GET")
	private.HandleFunc("/token", handler.TokenHandler).Methods("GET")
	admin := private.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RequireRolesMiddleware("admin"))
	admin.HandleFunc("/update-role", handler.UpdateUserRoleHandler).Methods("POST")
	admin.Handle("/monitoring", promhttp.Handler()).Methods("GET")

	log.Println("Server started :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
