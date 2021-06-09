package main

import (
	"os"
	"rpg/internal/handler"
	"rpg/internal/service"
	"rpg/pkg/hubber"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: false,
		FullTimestamp:    true,
	})
	logger.SetOutput(os.Stdout)
	messagePipe := make(chan hubber.IResponse, 100)
	srv := service.NewService(logger, messagePipe)
	hdl := handler.NewHandler(logger, messagePipe, srv)
	app := hubber.NewServer("3000", logger)
	app.Run(hdl)
}
