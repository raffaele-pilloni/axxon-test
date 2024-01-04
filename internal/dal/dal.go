package dal

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type DAL struct {
	gormDB       *gorm.DB
	queryTimeout time.Duration
}

func NewDAL(
	gormDB *gorm.DB,
	queryTimeout time.Duration,
) *DAL {
	return &DAL{
		gormDB:       gormDB,
		queryTimeout: queryTimeout,
	}
}

func (d DAL) FindById(ctx context.Context, entity interface{}, ID int) error {
	ctx, cancelCtx := context.WithTimeout(ctx, d.queryTimeout*time.Second)
	defer cancelCtx()

	query := d.gormDB.WithContext(ctx).Find(entity, ID)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return query.Error
}

func (d DAL) Save(ctx context.Context, entity interface{}) error {
	ctx, cancelCtx := context.WithTimeout(ctx, d.queryTimeout*time.Second)
	defer cancelCtx()

	insert := d.gormDB.WithContext(ctx).Save(entity)

	return insert.Error
}
