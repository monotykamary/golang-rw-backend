package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/monotykamary/golang-rw-backend/config"
	"github.com/monotykamary/golang-rw-backend/handler"
	"github.com/monotykamary/golang-rw-backend/log"
	"github.com/monotykamary/golang-rw-backend/repo"
	"github.com/monotykamary/golang-rw-backend/repo/pg"
	"github.com/monotykamary/golang-rw-backend/routes"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

func main() {
	//init config
	cls := config.DefaultConfigLoaders()
	cfg := config.LoadConfig(cls)

	// init logger
	undo := log.New()
	defer zap.L().Sync()
	defer undo()

	s, close := pg.NewPostgresStore(&cfg)
	defer close()

	router := setupRouter(cfg, s)

	router.Logger.Fatal(router.Start(":8080"))
}

func setupRouter(cfg config.Config, s repo.DBRepo) *echo.Echo {
	r := echo.New()
	r.Use(log.ZapLogger(zap.L()))
	r.Use(middleware.Recover())
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders: []string{"Origin", "Host",
			"Content-Type", "Content-Length",
			"Accept-Encoding", "Accept-Language", "Accept",
			"X-CSRF-Token", "Authorization", "X-Requested-With", "X-Access-Token", "Signature", "Challenge"},
		ExposeHeaders:    []string{"MeAllowMethodsntent-Length"},
		AllowCredentials: true,
	}))
	h := handler.NewHandler(cfg, s)
	url := echoSwagger.URL(cfg.BaseURL + "/swagger/doc.json")

	r.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))
	r.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("ok. Git version: %s", cfg.GitVersion))
	})

	routes.NewRoutes(r, h, cfg, s)
	return r
}
