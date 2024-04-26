package conf

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	storage_go "github.com/supabase-community/storage-go"
)

func NewSupabaseImage(log *logrus.Logger) *storage_go.Client {
	return storage_go.NewClient(
		fmt.Sprintf("https://%s.supabase.co/storage/v1", os.Getenv("SUPABASE_PROJECT_REFERENCE_ID")),
		os.Getenv("SUPABASE_SECRET_API_KEY"),
		map[string]string{fiber.HeaderContentType: "image/jpeg"},
	)
}

func NewSupabaseAttachment(log *logrus.Logger) *storage_go.Client {
	return storage_go.NewClient(
		fmt.Sprintf("https://%s.supabase.co/storage/v1", os.Getenv("SUPABASE_PROJECT_REFERENCE_ID")),
		os.Getenv("SUPABASE_SECRET_API_KEY"),
		map[string]string{
			fiber.HeaderContentType:        "application/octet-stream",
			fiber.HeaderContentDisposition: "attachment",
		},
	)
}
