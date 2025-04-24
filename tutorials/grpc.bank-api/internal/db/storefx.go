package db

import "go.uber.org/fx"

var Module = fx.Module("store", fx.Provide(NewDatabase))
