package data

import (
	"context"
	"starland-account/configs"
	"starland-account/internal/biz"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityLog struct {
	gorm.Model
	UUID         string `json:"uuid" gorm:"primary_key;size:255"`
	AccountID    string
	ActivityCode int
	ActivityName string
	Integral     int
}

type activityLogRepo struct {
	cfg  *configs.Config
	data *Data
}

func NewActivityLogRepo(c *configs.Config, data *Data) biz.ActivityLogRepo {
	return &activityLogRepo{
		cfg:  c,
		data: data,
	}
}

func (r *activityLogRepo) AddActivityLog(ctx context.Context, req *biz.ActivityLogRequest) error {
	actlog := &ActivityLog{
		UUID:         uuid.NewString(),
		AccountID:    req.AccountID,
		ActivityCode: req.ActivityCode,
		ActivityName: req.ActivityName,
		Integral:     req.Integral,
	}

	return r.data.db.WithContext(ctx).Model(&ActivityLog{}).Create(&actlog).Error
}

func (r *activityLogRepo) QueryActivityLog(ctx context.Context, account string, page, limit int) ([]*biz.ActivityLogResponse, int64, error) {
	var (
		actlogs []*ActivityLog
		count   int64
	)
	err := r.data.db.WithContext(ctx).Model(&ActivityLog{}).Where("account_id = ?", account).Offset((page - 1) * limit).Limit(limit).Order("updated_at desc ").Find(&actlogs).Error
	if err != nil {
		return nil, count, err
	}
	err = r.data.db.WithContext(ctx).Model(&ActivityLog{}).Where("account_id = ?", account).Count(&count).Error
	if err != nil {
		return nil, count, err
	}

	return makeActivityLogsToBizRes(actlogs), count, nil
}

func makeActivityLogsToBizRes(actlogs []*ActivityLog) []*biz.ActivityLogResponse {
	res := make([]*biz.ActivityLogResponse, len(actlogs))
	for i := range actlogs {
		res[i] = &biz.ActivityLogResponse{
			AccountID:    actlogs[i].AccountID,
			ActivityCode: actlogs[i].ActivityCode,
			ActivityName: actlogs[i].ActivityName,
			Integral:     actlogs[i].Integral,
			CreateAt:     actlogs[i].CreatedAt,
		}
	}
	return res
}
