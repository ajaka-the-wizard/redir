package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/models"
)

func (r *Sredis) GetProduct(ctx context.Context, productId int) (*models.Product, error) {
	var product models.Product
	key := fmt.Sprintf("product:%d", productId)
	err := r.rdb.HGetAll(ctx, key).Scan(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Sredis) SetProduct(ctx context.Context, product models.Product) error {
	key := fmt.Sprintf("product:%d", product.ProductId)
	exp := 20 * time.Minute
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, key, product)
	pipe.Expire(ctx, key, exp)
	_, err := pipe.Exec(ctx)
	return err
}
