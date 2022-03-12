package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"rpg/internal/server/controller"
	"rpg/internal/server/service"
	"rpg/pkg/hubber"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: false,
		FullTimestamp:    true,
	})
	logger.SetOutput(os.Stdout)
	srv := service.NewService(logger)
	ctrl := controller.NewController(logger, srv)
	app := hubber.NewServer(logger, "3000", ctrl)
	err := app.Run()
	if err != nil {
		logger.WithField("err", err).Error("failed to run app")
	}
}
