package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mirzahilmi/hackathon/internal/delivery/middleware"
	"github.com/mirzahilmi/hackathon/internal/delivery/rest"
	"github.com/mirzahilmi/hackathon/internal/pkg/conf"
	"github.com/mirzahilmi/hackathon/internal/pkg/email"
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

	if err := app.Listen(":" + os.Getenv("APP_PORT")); err != nil {
		log.Fatalf("Fail to start server: %v", err)
	}
}
