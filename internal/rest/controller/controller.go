package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
	"github.com/paranoideed/uni-products-svc/internal/domain"
	"github.com/paranoideed/uni-products-svc/internal/models"
	"github.com/paranoideed/uni-products-svc/internal/rest/requests"
	"github.com/paranoideed/uni-products-svc/internal/rest/responses"
	"github.com/paranoideed/uni-products-svc/internal/rest/scope"
)

type core interface {
	CreateProduct(ctx context.Context, req domain.CreateProductRequest) (models.Product, error)
	DeleteProduct(ctx context.Context, ID uuid.UUID) error
	GetProducts(ctx context.Context, opts ...domain.GetProductsOption) (pagi.Page[[]models.Product], error)
}

type Controller struct {
	core core
}

func New(core core) *Controller {
	return &Controller{
		core: core,
	}
}

const operationCreateProduct = "create_product"

func (s *Controller) CreateProduct(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).With("operation", operationCreateProduct)

	req, err := requests.CreateProduct(r)
	if err != nil {
		log.Error("invalid request", "error", err)
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	product, err := s.core.CreateProduct(r.Context(), domain.CreateProductRequest{
		Name:  req.Data.Attributes.Name,
		Price: req.Data.Attributes.Price,
	})
	switch {
	case errors.Is(err, domain.ErrorNotValidInput):
		log.Error("invalid request", "error", err)
		render.ResponseError(w, problems.BadRequest(err)...)
	case err != nil:
		log.Error("failed to create product", "error", err)
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("successfully created product", "product_id", product.ID)
		render.Response(w, http.StatusOK, responses.Product(product))
	}
}

const operationDeleteProduct = "delete_product"

func (s *Controller) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).With("operation", operationDeleteProduct)

	productID, err := uuid.Parse(chi.URLParam(r, "product_id"))
	if err != nil {
		log.Error("invalid request", "error", err)
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path": fmt.Errorf("invalid product id format: %s, must be uuid", chi.URLParam(r, "product_id")),
		})...)
		return
	}

	err = s.core.DeleteProduct(r.Context(), productID)
	switch {
	case errors.Is(err, domain.ErrorProductNotFound):
		log.Error("product not found", "product_id", productID)
		render.ResponseError(w, problems.NotFound(
			fmt.Sprintf("product with id %s not found", productID.String())),
		)
	case err != nil:
		log.Error("failed to delete product", "error", err)
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("successfully deleted product", "product_id", productID)
		render.Response(w, http.StatusNoContent, nil)
	}
}

const operationGetProducts = "get_products"

func (s *Controller) GetProducts(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).With("operation", operationGetProducts)

	params := r.URL.Query()
	var opts []domain.GetProductsOption

	// Pagination
	limit, offset := pagi.GetPagination(r)
	page := uint(1)
	if limit > 0 {
		page = offset/limit + 1
	}
	opts = append(opts, domain.WithPage(int(page)), domain.WithLimit(int(limit)))

	// Sort
	if sort := pagi.GetSort(r); sort != nil {
		opts = append(opts, domain.WithSort(domain.SortField(sort.Field), sort.Ascend))
	}

	// filter[name]
	if name := params.Get("filter[name]"); name != "" {
		opts = append(opts, domain.WithName(name))
	}

	// filter[price][gte] / filter[price][lte]
	var lowPrice, highPrice float64
	if v, err := strconv.ParseFloat(params.Get("filter[price][gte]"), 64); err == nil {
		lowPrice = v
	}
	if v, err := strconv.ParseFloat(params.Get("filter[price][lte]"), 64); err == nil {
		highPrice = v
	}
	if lowPrice > 0 || highPrice > 0 {
		opts = append(opts, domain.WithPriceRange(lowPrice, highPrice))
	}

	// filter[created_at][gte] / filter[created_at][lte]
	var startDate, endDate time.Time
	if t, err := time.Parse(time.RFC3339, params.Get("filter[created_at][gte]")); err == nil {
		startDate = t
	}
	if t, err := time.Parse(time.RFC3339, params.Get("filter[created_at][lte]")); err == nil {
		endDate = t
	}
	if !startDate.IsZero() || !endDate.IsZero() {
		opts = append(opts, domain.WithTimeRange(startDate, endDate))
	}

	result, err := s.core.GetProducts(r.Context(), opts...)
	switch {
	case errors.Is(err, domain.ErrorNotValidInput):
		log.Error("invalid request", "error", err)
		render.ResponseError(w, problems.BadRequest(err)...)
	case err != nil:
		log.Error("failed to get products", "error", err)
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("successfully retrieved products", "total", result.Total)
		render.Response(w, http.StatusOK, responses.ProductsCollection(r, result))
	}
}
