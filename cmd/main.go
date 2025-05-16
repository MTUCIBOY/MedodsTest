package main

import (
	"github.com/MTUCIBOY/MedodsTest/pkg/app"
	"github.com/MTUCIBOY/MedodsTest/pkg/config"
	"github.com/MTUCIBOY/MedodsTest/pkg/logger"
	"github.com/MTUCIBOY/MedodsTest/pkg/storage/psql"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	application := app.New(cfg, log, psql.New(log))

	application.Run()
}
