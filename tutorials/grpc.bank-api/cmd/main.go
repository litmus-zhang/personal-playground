package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"grpc.bank-api/internal/api"
	"grpc.bank-api/internal/config"
	"grpc.bank-api/internal/db"
)

func main() {
	fx.New(
		config.Module,
		db.Module,
		api.Module,
		fx.Provide(zap.NewDevelopment),
		fx.Invoke(func(cfg *config.Config, server *api.Server) error {
			return server.Start()
		}),
	).Run()
}
