package data

import (
	"context"
	"errors"
	"starland-account/configs"
	"starland-account/internal/biz"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	AccountID  string `json:"account_id" gorm:"primary_key;size:255"`
	Integral   int
	Received   int
	Email      string `gorm:"index:idx_member"`
	Name       string
	Provider   string `gorm:"index:idx_member"`
	AvatarURL  string
	State      int
	SolanaAddr string
	ClaimCount int
}

type accountRepo struct {
	cfg  *configs.Config
	data *Data
}

func NewAccountRepo(c *configs.Config, data *Data) biz.AccountRepo {
	return &accountRepo{
		cfg:  c,
		data: data,
	}
}

func (r *accountRepo) SaveAccount(ctx context.Context, req *biz.AccountRequest) error {

	var a *Account
	if err := r.data.db.WithContext(ctx).Model(&Account{}).Where("account_id = ?", req.AccountID).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a = &Account{
				AccountID:  req.AccountID,
				Email:      req.Email,
				Name:       req.Name,
				AvatarURL:  req.AvatarURL,
				Provider:   req.Provider,
				State:      req.State,
				ClaimCount: req.ClaimCount,
				SolanaAddr: req.SolanaAddr,
			}
		} else {
			return err
		}
	}
	a.Email = req.Email
	a.Email = req.Email
	a.Name = req.Name
	a.AvatarURL = req.AvatarURL
	a.Provider = req.Provider
	a.State = req.State
	if a.ClaimCount != 0 {
		a.ClaimCount = req.ClaimCount
	}

	if req.SolanaAddr != "" {
		a.SolanaAddr = req.SolanaAddr
	}
	return r.data.db.WithContext(ctx).Model(&Account{}).Where("account_id = ?", req.AccountID).Save(&a).Error
}

func (r *accountRepo) UpdateAccountIntegral(ctx context.Context, accountID string, integral int) error {
	if err := r.data.db.WithContext(ctx).Model(&Account{}).Where("account_id = ?", accountID).
		Update("integral", gorm.Expr("integral+ ?", integral)).Error; err != nil {
		return err
	}
	return nil
}

func (r *accountRepo) UpdateClaimPoints(ctx context.Context, accountID string, integral, received int) error {
	if err := r.data.db.WithContext(ctx).Model(&Account{}).Where("account_id = ?", accountID).
		Updates(Account{Integral: integral, Received: received}).Error; err != nil {
		return err
	}
	return nil

}

func (r *accountRepo) UpdateAddr(ctx context.Context, accountID string, addr string) error {
	if err := r.data.db.WithContext(ctx).Model(&Account{}).Where("account_id = ?", accountID).
		Updates(Account{SolanaAddr: addr}).Error; err != nil {
		return err
	}
	return nil

}

func (r *accountRepo) QueryAccount(ctx context.Context, accountID, email, provider string) (*biz.AccountResponse, error) {
	var a *Account
	if accountID == "" {
		if err := r.data.db.WithContext(ctx).Model(&Account{}).Where("email = ? and provider = ?", email, provider).First(&a).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}
			return nil, err
		}
	} else {
		if err := r.data.db.WithContext(ctx).Model(&Account{}).Where("account_id = ?", accountID).First(&a).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}
			return nil, err
		}
	}
	return makeAccountResponse(a), nil
}

func (r *accountRepo) QueryAccounts(ctx context.Context) ([]*biz.AccountResponse, error) {
	var a []*Account
	if err := r.data.db.WithContext(ctx).Model(&Account{}).Where("state <> ?", -1).Find(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return makeAccountResponses(a), nil
}

func makeAccountResponse(a *Account) *biz.AccountResponse {
	return &biz.AccountResponse{
		AccountID:  a.AccountID,
		Integral:   a.Integral,
		Received:   a.Received,
		AvatarURL:  a.AvatarURL,
		Name:       a.Name,
		SolanaAddr: a.SolanaAddr,
		ClaimCount: a.ClaimCount,
	}
}

func makeAccountResponses(req []*Account) []*biz.AccountResponse {
	res := make([]*biz.AccountResponse, len(req))
	for i := range req {
		res[i] = makeAccountResponse(req[i])
	}
	return res
}
