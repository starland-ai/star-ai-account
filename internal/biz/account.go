package biz

import (
	"context"
	"fmt"
	"starland-account/internal/pkg/bizerr"

	"go.uber.org/zap"
)

type AccountRequest struct {
	AccountID  string
	Email      string
	Name       string
	Provider   string
	AvatarURL  string
	SolanaAddr string
	State      int
	ClaimCount int
}

type AccountResponse struct {
	AccountID  string
	Integral   int
	Received   int
	Email      string
	Name       string
	Provider   string
	AvatarURL  string
	SolanaAddr string
	ClaimCount int
}

type AccountRepo interface {
	SaveAccount(context.Context, *AccountRequest) error
	QueryAccount(context.Context, string, string, string) (*AccountResponse, error)
	UpdateAccountIntegral(context.Context, string, int) error
	UpdateClaimPoints(context.Context, string, int, int) error
	QueryAccounts(context.Context) ([]*AccountResponse, error)
	UpdateAddr(context.Context, string, string) error
}

type AccountUsecase struct {
	repo AccountRepo
}

func NewAccountUsecase(repo AccountRepo) *AccountUsecase {
	return &AccountUsecase{repo: repo}
}

func (uc *AccountUsecase) SaveAccount(ctx context.Context, req *AccountRequest) error {

	if err := uc.repo.SaveAccount(ctx, req); err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("SaveAccount: save(%+v) to db err: %w", *req, err))
	}
	return nil
}

func (uc *AccountUsecase) UpdateClaimPoints(ctx context.Context, accountID string, integral, received int) error {
	zap.S().Info("req:", integral, received)
	if err := uc.repo.UpdateClaimPoints(ctx, accountID, integral, received); err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("UpdateAccountIntegral: claimPoints err: %w", err))
	}
	return nil
}

func (uc *AccountUsecase) UpdateAccountIntegral(ctx context.Context, accountID string, integral int) error {
	if err := uc.repo.UpdateAccountIntegral(ctx, accountID, integral); err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("UpdateAccountIntegral: save(%s) to db err: %w", accountID, err))
	}
	return nil
}

func (uc *AccountUsecase) QueryAccount(ctx context.Context, accountID, email, provider string) (*AccountResponse, error) {
	if email == "" && provider == "" {
		provider = "Blockchain"
	}

	res, err := uc.repo.QueryAccount(ctx, accountID, email, provider)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryAccount: query accountInfo(%s) to db err: %w", accountID, err))
	}
	if res == nil {
		return nil, bizerr.ErrAccountNotExist
	}

	return res, nil
}

func (uc *AccountUsecase) QueryAccounts(ctx context.Context) ([]*AccountResponse, error) {
	res, err := uc.repo.QueryAccounts(ctx)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryAccounts: query accounts err: %w", err))
	}
	return res, nil
}

func (uc *AccountUsecase) UpdateAddr(ctx context.Context, account, addr string) error {
	if err := uc.repo.UpdateAddr(ctx, account, addr); err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("UpdateAddr: save(%s) to db err: %w", addr, err))
	}
	return nil
}
