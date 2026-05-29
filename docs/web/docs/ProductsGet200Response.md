# ProductsGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]ProductsGet200ResponseDataInner**](ProductsGet200ResponseDataInner.md) |  | 
**Links** | [**ProductsGet200ResponseLinks**](ProductsGet200ResponseLinks.md) |  | 

## Methods

### NewProductsGet200Response

`func NewProductsGet200Response(data []ProductsGet200ResponseDataInner, links ProductsGet200ResponseLinks, ) *ProductsGet200Response`

NewProductsGet200Response instantiates a new ProductsGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProductsGet200ResponseWithDefaults

`func NewProductsGet200ResponseWithDefaults() *ProductsGet200Response`

NewProductsGet200ResponseWithDefaults instantiates a new ProductsGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *ProductsGet200Response) GetData() []ProductsGet200ResponseDataInner`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ProductsGet200Response) GetDataOk() (*[]ProductsGet200ResponseDataInner, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ProductsGet200Response) SetData(v []ProductsGet200ResponseDataInner)`

SetData sets Data field to given value.


### GetLinks

`func (o *ProductsGet200Response) GetLinks() ProductsGet200ResponseLinks`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *ProductsGet200Response) GetLinksOk() (*ProductsGet200ResponseLinks, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *ProductsGet200Response) SetLinks(v ProductsGet200ResponseLinks)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


