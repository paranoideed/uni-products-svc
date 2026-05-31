package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createProductRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		} `json:"attributes"`
	} `json:"data"`
}

type productResponse struct {
	Data struct {
		Id         uuid.UUID `json:"id"`
		Type       string    `json:"type"`
		Attributes struct {
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		} `json:"attributes"`
	} `json:"data"`
}

func TestCreateProducts(t *testing.T) {
	n := testCfg.Load.Products
	base := testCfg.REST.BaseURL

	for i := range n {
		t.Run(fmt.Sprintf("product_%d", i+1), func(t *testing.T) {
			t.Parallel()

			price := float64(rand.IntN(99999)+1) + rand.Float64()
			product := createProduct(t, base, fmt.Sprintf("product-%d", i+1), price)

			assert.Equal(t, fmt.Sprintf("product-%d", i+1), product.Data.Attributes.Name)
			assert.InDelta(t, price, product.Data.Attributes.Price, 0.01)
			assert.NotEmpty(t, product.Data.Id)
		})
	}
}

func TestCreateAndDeleteProducts(t *testing.T) {
	n := testCfg.Load.Products
	base := testCfg.REST.BaseURL

	for i := range n {
		t.Run(fmt.Sprintf("product_%d", i+1), func(t *testing.T) {
			t.Parallel()

			price := float64(rand.IntN(99999) + 1)
			product := createProduct(t, base, fmt.Sprintf("product-%d", i+1), price)

			deleteProduct(t, base, product.Data.Id)
		})
	}
}

func createProduct(t *testing.T, base, name string, price float64) productResponse {
	t.Helper()

	var body createProductRequest
	body.Data.Type = "product"
	body.Data.Attributes.Name = name
	body.Data.Attributes.Price = price

	resp := do(t, http.MethodPost, base+"/products", body)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode,
		"create product failed for name=%s", name)

	var result productResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	return result
}

func deleteProduct(t *testing.T, base string, id uuid.UUID) {
	t.Helper()

	resp := do(t, http.MethodDelete, fmt.Sprintf("%s/products/%s", base, id), nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode,
		"delete product failed for id=%s", id)
}

func do(t *testing.T, method, url string, body any) *http.Response {
	t.Helper()

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := testClient.Do(req)
	require.NoError(t, err)

	return resp
}
