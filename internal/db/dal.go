package db

import (
	"context"
	applicationerror "github.com/raffaele-pilloni/axxon-test/internal/error"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type DALInterface interface {
	FindByID(ctx context.Context, entity interface{}, ID int) error
	Save(ctx context.Context, entity interface{}) error
}

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

func (d DAL) FindByID(ctx context.Context, entity interface{}, ID int) error {
	ctx, cancelCtx := context.WithTimeout(ctx, d.queryTimeout*time.Second)
	defer cancelCtx()

	query := d.gormDB.WithContext(ctx).First(entity, ID)
	if query.Error == gorm.ErrRecordNotFound {
		return applicationerror.NewEntityNotFoundError(reflect.TypeOf(entity).Elem().Name(), ID)
	}

	return query.Error
}

func (d DAL) Save(ctx context.Context, entity interface{}) error {
	ctx, cancelCtx := context.WithTimeout(ctx, d.queryTimeout*time.Second)
	defer cancelCtx()

	insert := d.gormDB.WithContext(ctx).Save(entity)

	return insert.Error
}
