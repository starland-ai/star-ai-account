package activity

import (
	"starland-account/configs"
	"starland-account/internal/biz"
	"sync"
	"time"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewActivityService)

type ActivityService struct {
	cfg        *configs.Config
	activity   *biz.ActivityUsecase
	account    *biz.AccountUsecase
	actMap     map[int]*biz.ActivityResponse
	actMaplock sync.RWMutex
}

func NewActivityService(cfg *configs.Config,
	act *biz.ActivityUsecase, ac *biz.AccountUsecase) *ActivityService {
	s := &ActivityService{cfg: cfg,
		activity: act,
		account:  ac,
		actMap:   make(map[int]*biz.ActivityResponse)}
	go s.refreshTask()
	return s
}

type ActivityLogResponse struct {
	CreateAt     time.Time `json:"create_at"`
	Account      string    `json:"account"`
	ActivityName string    `json:"activity_name"`
	Integral     int       `json:"integral"`
}

type ActivityResponse struct {
	ActivityName string `json:"activity_name"`
	ActivityCode int    `json:"activity_code"`
	Integral     int    `json:"integral"`
}
