# PoolListOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Items** | [**[]PoolView**](PoolView.md) |  | 

## Methods

### NewPoolListOutputBody

`func NewPoolListOutputBody(items []PoolView, ) *PoolListOutputBody`

NewPoolListOutputBody instantiates a new PoolListOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPoolListOutputBodyWithDefaults

`func NewPoolListOutputBodyWithDefaults() *PoolListOutputBody`

NewPoolListOutputBodyWithDefaults instantiates a new PoolListOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *PoolListOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *PoolListOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *PoolListOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *PoolListOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetItems

`func (o *PoolListOutputBody) GetItems() []PoolView`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *PoolListOutputBody) GetItemsOk() (*[]PoolView, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *PoolListOutputBody) SetItems(v []PoolView)`

SetItems sets Items field to given value.


### SetItemsNil

`func (o *PoolListOutputBody) SetItemsNil(b bool)`

 SetItemsNil sets the value for Items to be an explicit nil

### UnsetItems
`func (o *PoolListOutputBody) UnsetItems()`

UnsetItems ensures that no value is present for Items, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


