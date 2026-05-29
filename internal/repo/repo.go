package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netbill/restkit/pagi"
	"github.com/paranoideed/uni-products-svc/internal/domain"
	"github.com/paranoideed/uni-products-svc/internal/models"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func scanProduct(row pgx.Row) (r models.Product, err error) {
	err = row.Scan(
		&r.ID,
		&r.Name,
		&r.Price,
		&r.CreatedAt,
		&r.DeletedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Product{}, domain.ErrorProductNotFound.Raise(err)
	case r.DeletedAt != nil:
		return models.Product{}, domain.ErrorProductNotFound.Raise(fmt.Errorf("product is deleted"))
	default:
		return r, err
	}
}

func (r *Repo) CreateProduct(ctx context.Context, req domain.CreateProductRequest) (models.Product, error) {
	const q = `
		INSERT INTO products (id, name, price)
		VALUES ($1, $2, $3)
		RETURNING id, name, price, created_at, deleted_at
	`

	result, err := scanProduct(r.db.QueryRow(ctx, q, uuid.New(), req.Name, req.Price))
	if err != nil {
		return models.Product{}, err
	}

	return result, nil
}

func (r *Repo) DeleteProduct(ctx context.Context, ID uuid.UUID) error {
	const q = `UPDATE products SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	tag, err := r.db.Exec(ctx, q, time.Now().UTC(), ID)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrorProductNotFound.Raise(
			fmt.Errorf("product with id %v not found", ID),
		)
	}

	return nil
}

func validSortField(f domain.SortField) string {
	switch f {
	case domain.SortByPrice:
		return "price"
	default:
		return "created_at"
	}
}

func (r *Repo) GetProducts(ctx context.Context, opts domain.GetProductsOptions) (pagi.Page[[]models.Product], error) {
	const filterClause = `
		WHERE deleted_at IS NULL
		  AND ($1::text IS NULL OR name ILIKE $1)
		  AND ($2::numeric IS NULL OR price >= $2)
		  AND ($3::numeric IS NULL OR price <= $3)
		  AND ($4::timestamptz IS NULL OR created_at >= $4)
		  AND ($5::timestamptz IS NULL OR created_at <= $5)`

	var name, lowPrice, highPrice, startDate, endDate any
	if opts.Name != "" {
		name = "%" + opts.Name + "%"
	}
	if opts.LowPrice > 0 {
		lowPrice = opts.LowPrice
	}
	if opts.HighPrice > 0 {
		highPrice = opts.HighPrice
	}
	if !opts.StartDate.IsZero() {
		startDate = opts.StartDate
	}
	if !opts.EndDate.IsZero() {
		endDate = opts.EndDate
	}

	var total uint
	err := r.db.QueryRow(ctx,
		"SELECT COUNT(*) FROM products"+filterClause,
		name, lowPrice, highPrice, startDate, endDate,
	).Scan(&total)
	if err != nil {
		return pagi.Page[[]models.Product]{}, fmt.Errorf("count products: %w", err)
	}

	sortDir := "DESC"
	if opts.SortASC {
		sortDir = "ASC"
	}

	limit, page := opts.Limit, opts.Page
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	dataQuery := fmt.Sprintf(
		"SELECT id, name, price, created_at, deleted_at FROM products%s ORDER BY %s %s LIMIT $6 OFFSET $7",
		filterClause, validSortField(opts.SortBy), sortDir,
	)

	rows, err := r.db.Query(ctx, dataQuery, name, lowPrice, highPrice, startDate, endDate, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Product]{}, fmt.Errorf("get products: %w", err)
	}
	defer rows.Close()

	products := make([]models.Product, 0, opts.Limit)

	for rows.Next() {
		var p models.Product
		p, err = scanProduct(rows)
		if err != nil {
			return pagi.Page[[]models.Product]{}, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return pagi.Page[[]models.Product]{}, fmt.Errorf("rows error: %w", err)
	}

	return pagi.Page[[]models.Product]{
		Data:  products,
		Page:  uint(page),
		Size:  uint(limit),
		Total: total,
	}, nil
}
