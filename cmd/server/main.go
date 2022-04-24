package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"rpg/internal/server/matchmaker"
	"rpg/internal/server/sync_controller"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(
		&logrus.TextFormatter{
			ForceColors:      true,
			DisableTimestamp: false,
			FullTimestamp:    true,
		},
	)
	logger.SetOutput(os.Stdout)
	services := matchmaker.NewService(logger)
	// ctrl := async_controller.NewController(logger, services)
	// app := hubber.NewServer(logger, "3000", ctrl)
	// err := app.Run()
	// if err != nil {
	// 	logger.WithField("err", err).Error("failed to run app")
	// }

	syncCtrl := sync_controller.NewController(logger, services)

	cfg := GetConfig()

	fmt.Printf("Starting server on %s...\n", cfg.Port)
	panic(http.ListenAndServe(":"+cfg.Port, syncCtrl))
}
