package account

import (
	"starland-account/configs"
	"starland-account/internal/biz"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewAccountService)

type AccountService struct {
	cfg     *configs.Config
	account *biz.AccountUsecase
}

func NewAccountService(cfg *configs.Config, account *biz.AccountUsecase) *AccountService {
	s := &AccountService{cfg: cfg, account: account}
	go s.solanaChainDataCheckTask()
	return s
}

type AccountResponse struct {
	AccountID  string `json:"account_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Provider   string `json:"provider"`
	AvatarURL  string `json:"avatar_url"`
	Integral   int    `json:"integral"`
	Received   int    `json:"received"`
	SolanaAddr string `json:"solana_addr"`
	ClaimCount int    `json:"claim_count"`
}

type AccountRequest struct {
	AccountID string
	Email     string
	Name      string
	Provider  string
	AvatarURL string
}

type ClaimPointsRequest struct {
	AccountID string
	Points    int
	IsOK      bool
}
