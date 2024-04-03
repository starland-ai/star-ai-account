package biz

import (
	"context"
	"fmt"
	"starland-account/internal/pkg/bizerr"
	"time"

	"go.uber.org/zap"
)

type ActivityLogRequest struct {
	AccountID    string
	ActivityCode int
	ActivityName string
	Integral     int
}

type ActivityLogResponse struct {
	AccountID    string
	ActivityCode int
	ActivityName string
	Integral     int
	CreateAt     time.Time
}

type ActivityResponse struct {
	ActivityCode int
	ActivityName string
	Integral     int
	Limit        int
}

type ActivityRepo interface {
	QueryActivity(context.Context) ([]*ActivityResponse, error)
	ConsumeActivityLimit(context.Context, string, int, time.Duration) error
	QueryActivityExpend(context.Context, string) (int, error)
}

type ActivityLogRepo interface {
	AddActivityLog(context.Context, *ActivityLogRequest) error
	QueryActivityLog(context.Context, string, int, int) ([]*ActivityLogResponse, int64, error)
}

type ActivityUsecase struct {
	activity    ActivityRepo
	activityLog ActivityLogRepo
}

func NewActivityUsecase(act ActivityRepo, actLog ActivityLogRepo) *ActivityUsecase {
	return &ActivityUsecase{activity: act, activityLog: actLog}
}

func (uc *ActivityUsecase) QueryActivity(ctx context.Context) (map[int]*ActivityResponse, error) {
	acts, err := uc.activity.QueryActivity(ctx)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivity: query activity err: %w", err))
	}

	res := make(map[int]*ActivityResponse, len(acts))

	for i := range acts {
		zap.S().Info(*acts[i])
		res[acts[i].ActivityCode] = acts[i]
	}
	return res, nil
}

func (uc *ActivityUsecase) AddActivityLog(ctx context.Context, req *ActivityLogRequest) error {
	err := uc.activityLog.AddActivityLog(ctx, req)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivity: query activity err: %w", err))
	}
	return nil
}

func (uc *ActivityUsecase) QueryActivityLog(ctx context.Context, addr string, page, limit int) ([]*ActivityLogResponse, int64, error) {
	res, count, err := uc.activityLog.QueryActivityLog(ctx, addr, page, limit)
	if err != nil {
		return nil, count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivity: query activity err: %w", err))
	}
	return res, count, nil
}

func (uc *ActivityUsecase) ConsumeActivityLimit(ctx context.Context, key string, n int, t time.Duration) error {
	err := uc.activity.ConsumeActivityLimit(ctx, key, n, t)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("ConsumeActivityLimit: err: %w", err))
	}
	return nil
}

func (uc *ActivityUsecase) QueryActivityExpend(ctx context.Context, key string) (int, error) {
	res, err := uc.activity.QueryActivityExpend(ctx, key)
	if err != nil {
		return 0, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityExpend: err: %w", err))
	}
	return res, nil
}
