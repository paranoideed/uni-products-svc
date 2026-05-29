# ProductsGet200ResponseDataInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | product ID | 
**Type** | **string** |  | 
**Attributes** | [**ProductsGet200ResponseDataInnerAttributes**](ProductsGet200ResponseDataInnerAttributes.md) |  | 

## Methods

### NewProductsGet200ResponseDataInner

`func NewProductsGet200ResponseDataInner(id uuid.UUID, type_ string, attributes ProductsGet200ResponseDataInnerAttributes, ) *ProductsGet200ResponseDataInner`

NewProductsGet200ResponseDataInner instantiates a new ProductsGet200ResponseDataInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProductsGet200ResponseDataInnerWithDefaults

`func NewProductsGet200ResponseDataInnerWithDefaults() *ProductsGet200ResponseDataInner`

NewProductsGet200ResponseDataInnerWithDefaults instantiates a new ProductsGet200ResponseDataInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ProductsGet200ResponseDataInner) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ProductsGet200ResponseDataInner) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ProductsGet200ResponseDataInner) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *ProductsGet200ResponseDataInner) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ProductsGet200ResponseDataInner) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ProductsGet200ResponseDataInner) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *ProductsGet200ResponseDataInner) GetAttributes() ProductsGet200ResponseDataInnerAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *ProductsGet200ResponseDataInner) GetAttributesOk() (*ProductsGet200ResponseDataInnerAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *ProductsGet200ResponseDataInner) SetAttributes(v ProductsGet200ResponseDataInnerAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


