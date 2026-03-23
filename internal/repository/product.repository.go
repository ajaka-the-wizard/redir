package repository

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePrivateKey(pool *pgxpool.Pool, cfg *configs.EnvData, productId int, hash string) (*models.Product, error) {
	ctx, cancel := utils.CreateContextWithStatedTime(cfg.CONTEXTTIMEOUT)
	defer cancel()
	var product models.Product
	query := `
	update client_keys
	SET private_key = $2
	WHERE product_id = $1
	RETURNING id,product_id,product_name,user_id,created_at,updated_at
	`
	err := pool.QueryRow(ctx, query, productId, hash).Scan(
		&product.ID,
		&product.ProductId,
		&product.ProductName,
		&product.UserId,
		&product.CreatedAt,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func GetProductById(pool *pgxpool.Pool, cfg *configs.EnvData, productId int) (*models.Product, error) {
	ctx, cancel := utils.CreateContextWithStatedTime(cfg.CONTEXTTIMEOUT)
	defer cancel()
	var product models.Product
	query := `
	SELECT id, product_id, user_id, private_key, created_at,updated_at
	FROM product
	WHERE product_id = $1
	`
	err := pool.QueryRow(ctx, query, productId).Scan(
		&product.ID,
		&product.ProductId,
		&product.UserId,
		&product.PrivateKey,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func CreateProduct(pool *pgxpool.Pool, cfg *configs.EnvData, data *domain.CreateProductDetails) (*models.Product, error) {
	ctx, cancel := utils.CreateContextWithStatedTime(cfg.CONTEXTTIMEOUT)
	defer cancel()
	var product models.Product

	query := `
	INSERT INTO product (product_name, user_id)
	VALUES($1, $2)
	RETURNING id,product_id,product_name,user_id,created_at,updated_at
	`
	err := pool.QueryRow(ctx, query, data.ProductName, data.UserId).Scan(
		&product.ID,
		&product.ProductId,
		&product.ProductName,
		&product.UserId,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
