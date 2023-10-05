package main

import (
	"fmt"
	"os"
	"time"

	"github.com/helloferdie/golib/libecho"
	"github.com/helloferdie/golib/libecho/libmiddleware"
	"github.com/helloferdie/pusher/handler"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	godotenv.Load()
}

func main() {
	fmt.Println("hello")

	e := echo.New()
	libecho.Initialize(e)

	rateLimiterConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 5, ExpiresIn: 2 * time.Minute},
		),
	}
	e.Use(middleware.RateLimiterWithConfig(rateLimiterConfig))

	headerSecretConfig := libmiddleware.VerifyHeaderSecretConfig{
		Field: os.Getenv("header_secret_field"),
		Value: os.Getenv("header_secret_value"),
	}

	e.POST("/deploy", handler.Deploy, libmiddleware.VerifyHeaderSecret(headerSecretConfig))
	libecho.StartHTTP(e)
}
