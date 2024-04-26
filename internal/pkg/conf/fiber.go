package conf

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

var StartTime = time.Now()

func NewFiber() *fiber.App {
	socketSharding, _ := strconv.ParseBool(os.Getenv("SOCKET_SHARDING"))
	bodyLimit, _ := strconv.Atoi(os.Getenv("BODY_LIMIT"))
	fiber := fiber.New(fiber.Config{
		AppName:           os.Getenv("APP_NAME"),
		Prefork:           socketSharding,
		BodyLimit:         bodyLimit * 1024 * 1024,
		DisableKeepalive:  true,
		StrictRouting:     true,
		CaseSensitive:     true,
		EnablePrintRoutes: true,
		JSONEncoder:       jsoniter.Marshal,
		JSONDecoder:       jsoniter.Unmarshal,
	})
	return fiber
}
