package activity

import (
	"context"
	"errors"
	"fmt"
	"starland-account/internal/biz"
	"starland-account/internal/pkg/bizerr"
	"time"

	"go.uber.org/zap"
)

const (
	Chat         = 0
	ChatIntegral = 1
)

var (
	ActivityKey = "starland-account:%d_%s"
	LimitError  = errors.New("You've reached the limit")
)

func (s *ActivityService) refreshTask() {
	defer func() {
		if p := recover(); p != nil {
			zap.S().Errorf("refreshTask: panic: %v", p)
		}
		s.refreshTask()
	}()
	s.refreshActMap()
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			s.refreshActMap()
		}
	}
}

func (s *ActivityService) refreshActMap() {
	ctx := context.Background()
	s.actMaplock.Lock()
	defer s.actMaplock.Unlock()
	res, err := s.activity.QueryActivity(ctx)
	if err != nil {
		zap.S().Errorf("Play: query activity to map err: %w", err)
		return
	}
	s.actMap = res
}

func (s *ActivityService) QueryActivityLogs(ctx context.Context, account string, page, limit int) ([]*ActivityLogResponse, int64, error) {
	res, count, err := s.activity.QueryActivityLog(ctx, account, page, limit)
	if err != nil {
		return nil, count, fmt.Errorf("QueryActivityLog: query err: %w", err)
	}
	return makeActivityLogs(res), count, nil
}

func (s *ActivityService) QueryActivitys(ctx context.Context) ([]*ActivityResponse, error) {
	res, err := s.activity.QueryActivity(ctx)
	if err != nil {
		return nil, fmt.Errorf("QueryActivityLog: query err: %w", err)
	}
	return makeActivitys(res), nil
}

func (s *ActivityService) Play(ctx context.Context, activityCode int, account string) error {
	key := fmt.Sprintf(ActivityKey, activityCode, account)
	if v, ok := s.actMap[activityCode]; ok {
		expend, err := s.activity.QueryActivityExpend(ctx, key)
		if err != nil {
			return fmt.Errorf("Play: query left activity count err: %w", err)
		}
		if expend >= v.Limit {
			zap.S().Infof("Play: Activity Count[account: %s activity: %s count: %d]", account, v.ActivityName, expend)
			return LimitError
		}
		err = s.account.UpdateAccountIntegral(ctx, account, v.Integral)
		if err != nil {
			return fmt.Errorf("Play: [%s] update account integral err: %w", v.ActivityName, err)
		}
		log := &biz.ActivityLogRequest{
			AccountID:    account,
			ActivityCode: v.ActivityCode,
			ActivityName: v.ActivityName,
			Integral:     v.Integral,
		}
		err = s.activity.AddActivityLog(ctx, log)
		if err != nil {
			return fmt.Errorf("Play: [%+v] add activity log err: %w", *log, err)
		}

		err = s.activity.ConsumeActivityLimit(ctx, key, expend+1, 24*time.Hour)
		if err != nil {
			return fmt.Errorf("Play: ConsumeActivityLimit err: %w", err)
		}
	} else {
		zap.S().Infof("Play: req activityCode:%d activity:(%+v)", activityCode, s.actMap)
		return bizerr.ErrActivityNotExist
	}

	return nil
}
func (s *ActivityService) QueryIsLimit(ctx context.Context, activityCode int, account string) (bool, error) {
	key := fmt.Sprintf(ActivityKey, activityCode, account)
	if v, ok := s.actMap[activityCode]; ok {
		expend, err := s.activity.QueryActivityExpend(ctx, key)
		if err != nil {
			zap.S().Errorf("Play: query left activity count err: %w", err)
			return false, nil
		}
		zap.S().Info(expend,v.Limit)
		if expend >= v.Limit {
			zap.S().Errorf("Play: Activity Count[account: %s activity: %s count: %d]", account, v.ActivityName, expend)
			return true, nil
		}
	}
	return false, nil
}

func makeActivityLogs(acts []*biz.ActivityLogResponse) []*ActivityLogResponse {
	res := make([]*ActivityLogResponse, len(acts))

	for i := range acts {
		res[i] = &ActivityLogResponse{
			Account:      acts[i].AccountID,
			CreateAt:     acts[i].CreateAt,
			ActivityName: acts[i].ActivityName,
			Integral:     acts[i].Integral,
		}
	}
	return res
}

func makeActivitys(acts map[int]*biz.ActivityResponse) []*ActivityResponse {
	res := make([]*ActivityResponse, len(acts))
	i := 0
	for k := range acts {
		res[i] = &ActivityResponse{
			ActivityCode: acts[k].ActivityCode,
			ActivityName: acts[k].ActivityName,
			Integral:     acts[k].Integral,
		}
		i++
	}
	return res
}
