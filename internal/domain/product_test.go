package domain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/restkit/pagi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newService() (*Service, *MockRepository) {
	repo := new(MockRepository)
	return NewService(repo), repo
}

// TestCreateProduct_Success - simple create product
func TestCreateProduct_Success(t *testing.T) {
	svc, repo := newService()

	req := CreateProductRequest{Name: "apple", Price: 1.5}
	expected := Product{
		ID:    uuid.New(),
		Name:  "apple",
		Price: 1.5,
	}

	repo.On("CreateProduct", mock.Anything, req).Return(expected, nil)

	result, err := svc.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// TestCreateProduct_InvalidReq - invalid request data or price invalid
func TestCreateProduct_InvalidReq(t *testing.T) {
	svc, _ := newService()

	_, err := svc.CreateProduct(context.Background(), CreateProductRequest{
		Name:  "",
		Price: 1.5,
	})

	assert.True(t, errors.Is(err, ErrorNotValidInput))

	_, err = svc.CreateProduct(context.Background(), CreateProductRequest{
		Name:  "apple",
		Price: -1,
	})

	assert.True(t, errors.Is(err, ErrorNotValidInput))

	_, err = svc.CreateProduct(context.Background(), CreateProductRequest{
		Name:  "apple",
		Price: 100001,
	})

	assert.True(t, errors.Is(err, ErrorNotValidInput))
}

// TestCreateProduct_EdgeCases - edge cases for create product
func TestCreateProduct_EdgeCases(t *testing.T) {
	svc, _ := newService()

	_, err := svc.CreateProduct(context.Background(), CreateProductRequest{
		Name:  "   ",
		Price: 1.5,
	})
	assert.True(t, errors.Is(err, ErrorNotValidInput))

	_, err = svc.CreateProduct(context.Background(), CreateProductRequest{
		Name:  "apple",
		Price: 0,
	})
	assert.True(t, errors.Is(err, ErrorNotValidInput))
}

// TestGetProducts_EdgeCases - edge cases for get products
func TestGetProducts_EdgeCases(t *testing.T) {
	svc, repo := newService()

	repo.On("GetProducts", mock.Anything, mock.AnythingOfType("GetProductsOptions")).Return(pagi.Page[[]Product]{}, nil)

	_, err := svc.GetProducts(context.Background(), WithPriceRange(5, 5))
	assert.NoError(t, err)

	_, err = svc.GetProducts(context.Background(), WithPriceRange(0, 0))
	assert.NoError(t, err)

	now := time.Now()
	_, err = svc.GetProducts(context.Background(), WithTimeRange(now, now))
	assert.NoError(t, err)
}

// TestDeleteProduct_Success - success delete product
func TestDeleteProduct_Success(t *testing.T) {
	svc, repo := newService()

	id := uuid.New()
	repo.On("DeleteProduct", mock.Anything, id).Return(nil)

	err := svc.DeleteProduct(context.Background(), id)

	assert.NoError(t, err)
}

// TestDeleteProduct_NotFound - delete product not found
func TestDeleteProduct_NotFound(t *testing.T) {
	svc, repo := newService()

	id := uuid.New()

	repo.On("DeleteProduct", mock.Anything, id).Return(
		ErrorProductNotFound.Raise(errors.New("not found")),
	)

	err := svc.DeleteProduct(context.Background(), id)

	assert.True(t, errors.Is(err, ErrorProductNotFound))
}

// TestGetProducts_Success - success get products with default options
func TestGetProducts_Success(t *testing.T) {
	svc, repo := newService()

	expected := pagi.Page[[]Product]{
		Data:  []Product{{ID: uuid.New(), Name: "apple", Price: 1.5}},
		Page:  1,
		Size:  20,
		Total: 1,
	}

	repo.On("GetProducts", mock.Anything, mock.AnythingOfType("GetProductsOptions")).Return(expected, nil)

	result, err := svc.GetProducts(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// TestGetProducts_InvalidParams - some params are invalid
func TestGetProducts_InvalidParams(t *testing.T) {
	//invalid price range
	svc, _ := newService()

	_, err := svc.GetProducts(context.Background(), WithPriceRange(10, 5))

	assert.True(t, errors.Is(err, ErrorNotValidInput))

	//invalid time range
	start := time.Now()
	end := start.Add(-time.Hour)

	_, err = svc.GetProducts(context.Background(), WithTimeRange(start, end))

	assert.True(t, errors.Is(err, ErrorNotValidInput))
}
