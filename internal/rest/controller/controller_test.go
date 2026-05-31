package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/restkit/pagi"
	"github.com/paranoideed/uni-products-svc/internal/domain"
	"github.com/paranoideed/uni-products-svc/internal/rest/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, core *mockCore) *httptest.Server {
	t.Helper()

	m := new(mockMetrics)
	m.On("RecordCreated", mock.Anything, mock.Anything).Maybe()
	m.On("RecordDeleted", mock.Anything, mock.Anything).Maybe()

	c := New(core, m)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(scope.CtxLog(r.Context(), log)))
		})
	})
	r.Post("/products", c.CreateProduct)
	r.Get("/products", c.GetProducts)
	r.Delete("/products/{product_id}", c.DeleteProduct)

	return httptest.NewServer(r)
}

func doJSON(t *testing.T, method, url string, body any) *http.Response {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}

	req, err := http.NewRequest(method, url, &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func decodeJSON(t *testing.T, r *http.Response, dst any) {
	t.Helper()
	defer r.Body.Close()
	require.NoError(t, json.NewDecoder(r.Body).Decode(dst))
}

// TestCreateProduct_Success - simple create product
func TestCreateProduct_Success(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	product := domain.Product{
		ID:        uuid.New(),
		Name:      "apple",
		Price:     9.99,
		CreatedAt: time.Now(),
	}
	core.On("CreateProduct", mock.Anything, domain.CreateProductRequest{
		Name: "apple", Price: 9.99,
	}).Return(product, nil)

	body := map[string]any{
		"data": map[string]any{
			"type": "product",
			"attributes": map[string]any{
				"name": "apple", "price": 9.99,
			},
		},
	}

	resp := doJSON(t, http.MethodPost, srv.URL+"/products", body)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]any
	decodeJSON(t, resp, &result)
	data := result["data"].(map[string]any)
	assert.Equal(t, product.ID.String(), data["id"])
	assert.Equal(t, "product", data["type"])

	core.AssertExpectations(t)
}

// TestCreateProduct_InvalidBody - invalid body
func TestCreateProduct_InvalidBody(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/products", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	core.AssertNotCalled(t, "CreateProduct")
}

// TestCreateProduct_WrongType - wrong type bad requests
func TestCreateProduct_WrongType(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	body := map[string]any{
		"data": map[string]any{
			"type":       "wrong",
			"attributes": map[string]any{"name": "apple", "price": 9.99},
		},
	}

	resp := doJSON(t, http.MethodPost, srv.URL+"/products", body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	core.AssertNotCalled(t, "CreateProduct")
}

// TestCreateProduct_DomainValidationError - invalid input
func TestCreateProduct_DomainValidationError(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	core.On("CreateProduct", mock.Anything, mock.Anything).
		Return(domain.Product{}, domain.ErrorNotValidInput.Raise(errors.New("name cannot be empty")))

	body := map[string]any{
		"data": map[string]any{
			"type":       "product",
			"attributes": map[string]any{"name": "", "price": 9.99},
		},
	}

	resp := doJSON(t, http.MethodPost, srv.URL+"/products", body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestCreateProduct_InternalError - 500
func TestCreateProduct_InternalError(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	core.On("CreateProduct", mock.Anything, mock.Anything).
		Return(domain.Product{}, errors.New("db is down"))

	body := map[string]any{
		"data": map[string]any{
			"type":       "product",
			"attributes": map[string]any{"name": "apple", "price": 9.99},
		},
	}

	resp := doJSON(t, http.MethodPost, srv.URL+"/products", body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// TestDeleteProduct_Success - delete product
func TestDeleteProduct_Success(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	id := uuid.New()
	core.On("DeleteProduct", mock.Anything, id).Return(nil)

	resp := doJSON(t, http.MethodDelete, fmt.Sprintf("%s/products/%s", srv.URL, id), nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	core.AssertExpectations(t)
}

// TestDeleteProduct_InvalidUUID - invalid data
func TestDeleteProduct_InvalidUUID(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	resp := doJSON(t, http.MethodDelete, srv.URL+"/products/not-a-uuid", nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	core.AssertNotCalled(t, "DeleteProduct")
}

// TestDeleteProduct_NotFound - not found
func TestDeleteProduct_NotFound(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	id := uuid.New()
	core.On("DeleteProduct", mock.Anything, id).
		Return(domain.ErrorProductNotFound.Raise(errors.New("not found")))

	resp := doJSON(t, http.MethodDelete, fmt.Sprintf("%s/products/%s", srv.URL, id), nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestDeleteProduct_InternalError - 500
func TestDeleteProduct_InternalError(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	id := uuid.New()
	core.On("DeleteProduct", mock.Anything, id).Return(errors.New("db is down"))

	resp := doJSON(t, http.MethodDelete, fmt.Sprintf("%s/products/%s", srv.URL, id), nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// TestGetProducts_Success - get success
func TestGetProducts_Success(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	expected := pagi.Page[[]domain.Product]{
		Data:  []domain.Product{{ID: uuid.New(), Name: "apple", Price: 9.99}},
		Page:  1,
		Size:  20,
		Total: 1,
	}
	core.On("GetProducts", mock.Anything, mock.Anything).Return(expected, nil)

	resp := doJSON(t, http.MethodGet, srv.URL+"/products", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	decodeJSON(t, resp, &result)
	data := result["data"].([]any)
	assert.Len(t, data, 1)
}

// TestGetProducts_EmptyList - empty list
func TestGetProducts_EmptyList(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	core.On("GetProducts", mock.Anything, mock.Anything).
		Return(pagi.Page[[]domain.Product]{Data: []domain.Product{}, Page: 1, Size: 20}, nil)

	resp := doJSON(t, http.MethodGet, srv.URL+"/products", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	decodeJSON(t, resp, &result)
	data := result["data"].([]any)
	assert.Empty(t, data)
}

// TestGetProducts_InternalError - 500
func TestGetProducts_InternalError(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	core.On("GetProducts", mock.Anything, mock.Anything).
		Return(pagi.Page[[]domain.Product]{}, errors.New("db is down"))

	resp := doJSON(t, http.MethodGet, srv.URL+"/products", nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// TestGetProducts_ValidationError - invalid query params
func TestGetProducts_ValidationError(t *testing.T) {
	core := new(mockCore)
	srv := newTestServer(t, core)
	defer srv.Close()

	core.On("GetProducts", mock.Anything, mock.Anything).
		Return(pagi.Page[[]domain.Product]{}, domain.ErrorNotValidInput.Raise(errors.New("invalid price range")))

	resp := doJSON(t, http.MethodGet, srv.URL+"/products?filter[price][gte]=100&filter[price][lte]=10", nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
