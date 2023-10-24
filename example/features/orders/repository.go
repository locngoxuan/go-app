package orders

import (
	"context"
	"go-app/appbase/pkg/interfaces"
	"go-app/example/entity/model"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/dig"
)

type OrderRepository interface {
	Save(value interface{}) error
	SaveBatches(value interface{}, size int) error
	DeleteOrderById(id int64) (int64, error)
	GetAllOrders() ([]model.Order, error)
}

type orderRepositoryParams struct {
	dig.In
	Ctx        *cli.Context
	Repository interfaces.Repository
}

func NewOrderRepository(params orderRepositoryParams) OrderRepository {
	return &orderRepository{
		maxTimeout: params.Ctx.Duration(interfaces.FlagDbMaxExecTime),
		repository: params.Repository,
	}
}

type orderRepository struct {
	maxTimeout time.Duration
	repository interfaces.Repository
}

// GetAllOrders implements OrderRepository.
func (r *orderRepository) GetAllOrders() ([]model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.maxTimeout)
	defer cancel()
	var entities []model.Order
	tx := r.repository.Database().WithContext(ctx).Find(&entities)
	return entities, tx.Error
}

// DeleteOrderById implements OrderRepository.
func (r *orderRepository) DeleteOrderById(id int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.maxTimeout)
	defer cancel()
	tx := r.repository.Database().WithContext(ctx).Where("id = ?", id).Delete(model.Order{})
	return tx.RowsAffected, tx.Error
}

// Save implements OrderRepository.
func (r *orderRepository) Save(value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.maxTimeout)
	defer cancel()
	tx := r.repository.Database().WithContext(ctx).Create(value)
	return tx.Error
}

// SaveBatches implements OrderRepository.
func (r *orderRepository) SaveBatches(value interface{}, size int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.maxTimeout)
	defer cancel()
	tx := r.repository.Database().WithContext(ctx).CreateInBatches(value, size)
	return tx.Error
}
