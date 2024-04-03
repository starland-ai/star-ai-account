//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"starland-account/configs"
	"starland-account/internal/biz"
	"starland-account/internal/data"
	"starland-account/internal/service"
	account_service "starland-account/internal/service/account"
	activity_service "starland-account/internal/service/activity"

	"github.com/google/wire"
)

// initApp
func initApp(cfg *configs.Config) (*service.Service, error) {
	panic(wire.Build(data.ProviderSet,
		biz.ProviderSet,
		account_service.ProviderSet,
		activity_service.ProviderSet,
		service.ProviderSet))
}
