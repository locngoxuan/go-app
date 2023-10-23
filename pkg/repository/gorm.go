package repository

import "context"

func (r GormRepository) Save(value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.MaxExecTime)
	defer cancel()
	tx := r.db.WithContext(ctx).Create(value)
	return tx.Error
}

func (r GormRepository) SaveBatches(value interface{}, size int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.MaxExecTime)
	defer cancel()
	tx := r.db.WithContext(ctx).CreateInBatches(value, size)
	return tx.Error
}
