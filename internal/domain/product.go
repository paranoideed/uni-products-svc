package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/restkit/pagi"
)

type Product struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Price     float32    `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

//go:generate mockery --name=Repository --inpackage --filename=mock_repository_test.go
type Repository interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (Product, error)
	DeleteProduct(ctx context.Context, ID uuid.UUID) error
	GetProducts(ctx context.Context, opts GetProductsOptions) (pagi.Page[[]Product], error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateProductRequest struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func (s Service) CreateProduct(ctx context.Context, req CreateProductRequest) (Product, error) {
	if strings.TrimSpace(req.Name) == "" {
		return Product{}, ErrorNotValidInput.Raise(errors.New("name cannot be empty"))
	}
	if req.Price <= 0 || req.Price > 100000 {
		return Product{}, ErrorNotValidInput.Raise(errors.New("price must be between 0 and 100000"))
	}

	return s.repo.CreateProduct(ctx, req)
}

func (s Service) DeleteProduct(ctx context.Context, ID uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, ID)
}

type SortField string

const (
	SortByPrice     SortField = "price"
	SortByCreatedAt SortField = "created_at"
)

type GetProductsOptions struct {
	Page  int
	Limit int

	Name      string
	LowPrice  float64
	HighPrice float64

	StartDate time.Time
	EndDate   time.Time

	SortBy  SortField
	SortASC bool
}

type GetProductsOption func(*GetProductsOptions)

func ApplyGetProductsOptions(opt []GetProductsOption) GetProductsOptions {
	res := GetProductsOptions{
		Page:  1,
		Limit: 20,
	}

	for _, fn := range opt {
		fn(&res)
	}

	return res
}

func WithPage(page int) GetProductsOption {
	return func(o *GetProductsOptions) {
		o.Page = page
	}
}

func WithLimit(limit int) GetProductsOption {
	return func(o *GetProductsOptions) {
		o.Limit = limit
	}
}

func WithName(name string) GetProductsOption {
	return func(o *GetProductsOptions) {
		o.Name = name
	}
}

func WithPriceRange(low, high float64) GetProductsOption {
	return func(o *GetProductsOptions) {
		o.LowPrice = low
		o.HighPrice = high
	}
}

func WithTimeRange(start, end time.Time) GetProductsOption {
	return func(o *GetProductsOptions) {
		o.StartDate = start
		o.EndDate = end
	}
}

func WithSort(field SortField, asc bool) GetProductsOption {
	return func(o *GetProductsOptions) {
		o.SortBy = field
		o.SortASC = asc
	}
}

func (s Service) GetProducts(
	ctx context.Context,
	opts ...GetProductsOption,
) (pagi.Page[[]Product], error) {
	o := ApplyGetProductsOptions(opts)

	if o.LowPrice > 0 && o.HighPrice > 0 && o.LowPrice > o.HighPrice {
		return pagi.Page[[]Product]{}, ErrorNotValidInput.Raise(
			errors.New("low price cannot be greater than high price"),
		)
	}
	if !o.StartDate.IsZero() && !o.EndDate.IsZero() && o.StartDate.After(o.EndDate) {
		return pagi.Page[[]Product]{}, ErrorNotValidInput.Raise(
			errors.New("start date cannot be after end date"),
		)
	}

	return s.repo.GetProducts(ctx, o)
}
