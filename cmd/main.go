package main

import (
	"github.com/MTUCIBOY/MedodsTest/pkg/app"
	"github.com/MTUCIBOY/MedodsTest/pkg/config"
	"github.com/MTUCIBOY/MedodsTest/pkg/logger"
	"github.com/MTUCIBOY/MedodsTest/pkg/storage/psql"
)

// @title Auth Service API
// @version 1.0
// @description Сервис авторизации, который использует два JWT-токена:
//
// - **Access Token**: Токен с информацией о пользователе в заголовке "Access-Token: <token>"
// - **Refresh Token**: Токен для связки БД и Access-token в заголовке "Refres-Token: <token>"

// @host localhost:8888
// @BasePath /

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	application := app.New(cfg, log, psql.New(log))

	application.Run()
}
