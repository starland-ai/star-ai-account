package service

import (
	"starland-account/internal/service/account"
	"starland-account/internal/service/activity"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewService)

type Service struct {
	Account  *account.AccountService
	Activity *activity.ActivityService
}

func NewService(account *account.AccountService, activity *activity.ActivityService) *Service {
	return &Service{Account: account, Activity: activity}
}
