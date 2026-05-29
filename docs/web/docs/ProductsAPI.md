# \ProductsAPI

All URIs are relative to *http://localhost:8001*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ProductsGet**](ProductsAPI.md#ProductsGet) | **Get** /products | Get products
[**ProductsPost**](ProductsAPI.md#ProductsPost) | **Post** /products | Create a new product
[**ProductsProductIdDelete**](ProductsAPI.md#ProductsProductIdDelete) | **Delete** /products/{product_id} | Delete a product



## ProductsGet

> ProductsCollection ProductsGet(ctx).PageLimit(pageLimit).PageOffset(pageOffset).FilterName(filterName).FilterPriceGte(filterPriceGte).FilterPriceLte(filterPriceLte).FilterCreatedAtGte(filterCreatedAtGte).FilterCreatedAtLte(filterCreatedAtLte).Sort(sort).Execute()

Get products



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	pageLimit := int32(56) // int32 | Max number of items to return (optional) (default to 20)
	pageOffset := int32(56) // int32 | Number of items to skip (optional) (default to 0)
	filterName := "filterName_example" // string | Filter by product name (partial match) (optional)
	filterPriceGte := float32(3.4) // float32 | Filter products with price greater than or equal to this value (optional)
	filterPriceLte := float32(3.4) // float32 | Filter products with price less than or equal to this value (optional)
	filterCreatedAtGte := time.Now() // time.Time | Filter products created at or after this timestamp (RFC3339) (optional)
	filterCreatedAtLte := time.Now() // time.Time | Filter products created at or before this timestamp (RFC3339) (optional)
	sort := "sort_example" // string | Sort field and direction. Prefix `-` means descending order. Example: `sort=price` (cheapest first), `sort=-price` (most expensive first)  (optional) (default to "-created_at")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProductsAPI.ProductsGet(context.Background()).PageLimit(pageLimit).PageOffset(pageOffset).FilterName(filterName).FilterPriceGte(filterPriceGte).FilterPriceLte(filterPriceLte).FilterCreatedAtGte(filterCreatedAtGte).FilterCreatedAtLte(filterCreatedAtLte).Sort(sort).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProductsAPI.ProductsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProductsGet`: ProductsCollection
	fmt.Fprintf(os.Stdout, "Response from `ProductsAPI.ProductsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiProductsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageLimit** | **int32** | Max number of items to return | [default to 20]
 **pageOffset** | **int32** | Number of items to skip | [default to 0]
 **filterName** | **string** | Filter by product name (partial match) | 
 **filterPriceGte** | **float32** | Filter products with price greater than or equal to this value | 
 **filterPriceLte** | **float32** | Filter products with price less than or equal to this value | 
 **filterCreatedAtGte** | **time.Time** | Filter products created at or after this timestamp (RFC3339) | 
 **filterCreatedAtLte** | **time.Time** | Filter products created at or before this timestamp (RFC3339) | 
 **sort** | **string** | Sort field and direction. Prefix &#x60;-&#x60; means descending order. Example: &#x60;sort&#x3D;price&#x60; (cheapest first), &#x60;sort&#x3D;-price&#x60; (most expensive first)  | [default to &quot;-created_at&quot;]

### Return type

[**ProductsCollection**](ProductsCollection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProductsPost

> Product ProductsPost(ctx).CreateProduct(createProduct).Execute()

Create a new product



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	createProduct := *openapiclient.NewCreateProduct(*openapiclient.NewCreateProductData("Type_example", *openapiclient.NewCreateProductDataAttributes("Name_example", float32(123)))) // CreateProduct | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProductsAPI.ProductsPost(context.Background()).CreateProduct(createProduct).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProductsAPI.ProductsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProductsPost`: Product
	fmt.Fprintf(os.Stdout, "Response from `ProductsAPI.ProductsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiProductsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createProduct** | [**CreateProduct**](CreateProduct.md) |  | 

### Return type

[**Product**](Product.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProductsProductIdDelete

> ProductsProductIdDelete(ctx, productId).Execute()

Delete a product



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	productId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | product ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ProductsAPI.ProductsProductIdDelete(context.Background(), productId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProductsAPI.ProductsProductIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**productId** | **uuid.UUID** | product ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiProductsProductIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

