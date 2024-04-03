package data

import (
	"context"
	"errors"
	"starland-account/configs"
	"starland-account/internal/biz"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type Activity struct {
	gorm.Model
	UUID         string `json:"uuid" gorm:"primary_key;size:255"`
	ActivityCode int
	ActivityName string
	Integral     int
	Limit        int
}

type activityRepo struct {
	cfg  *configs.Config
	data *Data
}

func NewActivityRepo(c *configs.Config, data *Data) biz.ActivityRepo {
	return &activityRepo{
		cfg:  c,
		data: data,
	}
}

func (r *activityRepo) QueryActivity(ctx context.Context) ([]*biz.ActivityResponse, error) {
	var (
		acts []*Activity
	)
	err := r.data.db.WithContext(ctx).Model(&Activity{}).Find(&acts).Error
	if err != nil {
		return nil, err
	}
	return makeActivityToBizResponse(acts), nil
}

func (r *activityRepo) ConsumeActivityLimit(ctx context.Context, key string, n int, timeOut time.Duration) error {
	_, err := r.data.rdb.WithContext(ctx).Set(key, n, timeOut).Result()
	return err
}

func (r *activityRepo) QueryActivityExpend(ctx context.Context, key string) (int, error) {
	value, err := r.data.rdb.WithContext(ctx).Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	res, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return res, err
}

func makeActivityToBizResponse(acts []*Activity) []*biz.ActivityResponse {
	res := make([]*biz.ActivityResponse, len(acts))
	for i := range acts {
		res[i] = &biz.ActivityResponse{
			ActivityCode: acts[i].ActivityCode,
			ActivityName: acts[i].ActivityName,
			Integral:     acts[i].Integral,
			Limit:        acts[i].Limit,
		}
	}
	return res
}
