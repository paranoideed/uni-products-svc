# ProductsGet400Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Errors** | [**[]ProductsGet400ResponseErrorsInner**](ProductsGet400ResponseErrorsInner.md) | Non empty array of errors occurred during request processing | 

## Methods

### NewProductsGet400Response

`func NewProductsGet400Response(errors []ProductsGet400ResponseErrorsInner, ) *ProductsGet400Response`

NewProductsGet400Response instantiates a new ProductsGet400Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProductsGet400ResponseWithDefaults

`func NewProductsGet400ResponseWithDefaults() *ProductsGet400Response`

NewProductsGet400ResponseWithDefaults instantiates a new ProductsGet400Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrors

`func (o *ProductsGet400Response) GetErrors() []ProductsGet400ResponseErrorsInner`

GetErrors returns the Errors field if non-nil, zero value otherwise.

### GetErrorsOk

`func (o *ProductsGet400Response) GetErrorsOk() (*[]ProductsGet400ResponseErrorsInner, bool)`

GetErrorsOk returns a tuple with the Errors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrors

`func (o *ProductsGet400Response) SetErrors(v []ProductsGet400ResponseErrorsInner)`

SetErrors sets Errors field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


