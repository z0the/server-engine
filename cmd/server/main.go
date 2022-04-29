package main

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"rpg/internal/server/async_controller"
	"rpg/internal/server/auth"
	"rpg/internal/server/eventbus"
	"rpg/internal/server/matchmaker"
	"rpg/internal/server/sync_controller"
	"rpg/pkg/hubber"
)

func main() {
	cfg := GetConfig()
	var logger *zap.Logger
	var err error

	if cfg.IsDev {
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	}
	sugarLg := logger.Sugar()

	bus := eventbus.NewBus(sugarLg)

	authSvc := auth.NewAuthService(sugarLg)

	matchMakerSvc := matchmaker.NewMatchMakerService(sugarLg, bus)
	syncCtrl := sync_controller.NewController(sugarLg, authSvc, matchMakerSvc)

	ctrl := async_controller.NewController(sugarLg, bus, authSvc, matchMakerSvc)
	server := hubber.NewServer(sugarLg, "3000", ctrl)

	go func() {
		err = server.Run()
		if err != nil {
			sugarLg.Errorw("failed to run server", "err", err)
		}
	}()

	fmt.Printf("Starting server on %s...\n", cfg.Port)
	panic(http.ListenAndServe(":"+cfg.Port, syncCtrl))
}
