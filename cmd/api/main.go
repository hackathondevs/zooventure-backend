package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mirzahilmi/hackathon/internal/delivery/middleware"
	"github.com/mirzahilmi/hackathon/internal/delivery/rest"
	"github.com/mirzahilmi/hackathon/internal/pkg/conf"
	"github.com/mirzahilmi/hackathon/internal/pkg/email"
	"github.com/mirzahilmi/hackathon/internal/repository"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Fail to load .env file: %v", err)
	}

	log := conf.NewLogger()
	app := conf.NewFiber(log)
	app.Use(middleware.CORS())
	db := conf.NewDatabase(log)
	supabaseImg := conf.NewSupabaseImage(log)
	mailer := email.NewMailer()
	validator := conf.NewValidator()
	api := app.Group("/api")
	rest.RegisterUtilsHandler(api)

	userRepo := repository.NewUserRepository(db)
	verifRepo := repository.NewVerificationRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo, verifRepo, mailer, log)
	rest.RegisterAuthHandler(authUsecase, validator, api)
	userUsecase := usecase.NewUserUsecase(userRepo, supabaseImg, log)
	rest.RegisterUserHandler(userUsecase, validator, api)

	animalUsecase := usecase.NewAnimalUsecase()
	rest.RegisterAnimalHandler(animalUsecase, api)

	if err := app.Listen(":" + os.Getenv("APP_PORT")); err != nil {
		log.Fatalf("Fail to start server: %v", err)
	}
}
