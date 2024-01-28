package node

import (
	"context"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core/box/service"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/libp2p/go-libp2p/core/host"

	"go.uber.org/fx"
)

func NewService(lc fx.Lifecycle, p2pHost host.Host, cfg *config.Config) (*service.HttpServer, error) {
	httpserver := service.NewHttpServer(cfg, p2pHost)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			httpserver.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			httpserver.Shutdown()
			return nil
		},
	})
	return httpserver, nil
}
