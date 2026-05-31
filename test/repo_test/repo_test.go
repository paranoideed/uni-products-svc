package repo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paranoideed/uni-products-svc/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createProduct — хелпер: создаёт продукт и регистрирует cleanup для удаления.
func createProduct(t *testing.T, name string, price float64) domain.Product {
	t.Helper()

	p, err := testRepo.CreateProduct(context.Background(), domain.CreateProductRequest{
		Name:  name,
		Price: price,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		// игнорируем ошибку — продукт мог быть удалён в самом тесте
		_ = testRepo.DeleteProduct(context.Background(), p.ID)
	})

	return p
}

func TestCreateProduct(t *testing.T) {
	p := createProduct(t, "apple", 9.99)

	assert.NotEqual(t, uuid.Nil, p.ID)
	assert.Equal(t, "apple", p.Name)
	assert.InDelta(t, 9.99, p.Price, 0.001)
	assert.False(t, p.CreatedAt.IsZero())
	assert.Nil(t, p.DeletedAt)
}

func TestDeleteProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		p := createProduct(t, "to-delete", 1.00)

		err := testRepo.DeleteProduct(context.Background(), p.ID)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		err := testRepo.DeleteProduct(context.Background(), uuid.New())

		assert.True(t, errors.Is(err, domain.ErrorProductNotFound))
	})

	t.Run("already deleted", func(t *testing.T) {
		p := createProduct(t, "double-delete", 1.00)

		require.NoError(t, testRepo.DeleteProduct(context.Background(), p.ID))

		err := testRepo.DeleteProduct(context.Background(), p.ID)
		assert.True(t, errors.Is(err, domain.ErrorProductNotFound))
	})
}

func TestGetProducts(t *testing.T) {
	// создаём изолированный набор продуктов с уникальным префиксом
	prefix := "test-get-" + uuid.New().String()[:8] + "-"

	p1 := createProduct(t, prefix+"apple", 10.00)
	p2 := createProduct(t, prefix+"banana", 20.00)
	p3 := createProduct(t, prefix+"cherry", 30.00)

	t.Run("no filters returns all", func(t *testing.T) {
		page, err := testRepo.GetProducts(context.Background(), domain.GetProductsOptions{
			Page: 1, Limit: 100,
			Name: prefix,
		})
		require.NoError(t, err)
		assert.Equal(t, uint(3), page.Total)
	})

	t.Run("filter by name", func(t *testing.T) {
		page, err := testRepo.GetProducts(context.Background(), domain.GetProductsOptions{
			Page: 1, Limit: 10,
			Name: prefix + "ban",
		})
		require.NoError(t, err)
		assert.Equal(t, uint(1), page.Total)
		assert.Equal(t, p2.ID, page.Data[0].ID)
	})

	t.Run("filter by price range", func(t *testing.T) {
		page, err := testRepo.GetProducts(context.Background(), domain.GetProductsOptions{
			Page: 1, Limit: 10,
			Name:      prefix,
			LowPrice:  15.0,
			HighPrice: 25.0,
		})
		require.NoError(t, err)
		assert.Equal(t, uint(1), page.Total)
		assert.Equal(t, p2.ID, page.Data[0].ID)
	})

	t.Run("filter by date range", func(t *testing.T) {
		page, err := testRepo.GetProducts(context.Background(), domain.GetProductsOptions{
			Page:      1,
			Limit:     10,
			Name:      prefix,
			StartDate: p1.CreatedAt.Add(-time.Second),
			EndDate:   p3.CreatedAt.Add(time.Second),
		})
		require.NoError(t, err)
		assert.Equal(t, uint(3), page.Total)
	})

	t.Run("pagination", func(t *testing.T) {
		page, err := testRepo.GetProducts(context.Background(), domain.GetProductsOptions{
			Page: 1, Limit: 2,
			Name: prefix,
		})
		require.NoError(t, err)
		assert.Equal(t, uint(3), page.Total)
		assert.Len(t, page.Data, 2)

		page2, err := testRepo.GetProducts(context.Background(), domain.GetProductsOptions{
			Page: 2, Limit: 2,
			Name: prefix,
		})
		require.NoError(t, err)
		assert.Len(t, page2.Data, 1)
	})

	t.Run("sort by price asc", func(t *testing.T) {
		page, err := testRepo.GetProducts(context.Background(), domain.GetProductsOptions{
			Page: 1, Limit: 10,
			Name:    prefix,
			SortBy:  domain.SortByPrice,
			SortASC: true,
		})
		require.NoError(t, err)
		require.Len(t, page.Data, 3)
		assert.Equal(t, p1.ID, page.Data[0].ID)
		assert.Equal(t, p2.ID, page.Data[1].ID)
		assert.Equal(t, p3.ID, page.Data[2].ID)
	})

	t.Run("sort by price desc", func(t *testing.T) {
		page, err := testRepo.GetProducts(context.Background(), domain.GetProductsOptions{
			Page: 1, Limit: 10,
			Name:    prefix,
			SortBy:  domain.SortByPrice,
			SortASC: false,
		})
		require.NoError(t, err)
		require.Len(t, page.Data, 3)
		assert.Equal(t, p3.ID, page.Data[0].ID)
	})

	t.Run("deleted product not returned", func(t *testing.T) {
		p := createProduct(t, prefix+"deleted", 5.00)
		require.NoError(t, testRepo.DeleteProduct(context.Background(), p.ID))

		page, err := testRepo.GetProducts(context.Background(), domain.GetProductsOptions{
			Page: 1, Limit: 100,
			Name: prefix,
		})
		require.NoError(t, err)
		assert.Equal(t, uint(3), page.Total)

		for _, product := range page.Data {
			assert.NotEqual(t, p.ID, product.ID)
		}
	})
}
