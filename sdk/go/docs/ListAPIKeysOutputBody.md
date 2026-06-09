# ListAPIKeysOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Items** | [**[]ApiKeyView**](ApiKeyView.md) |  | 

## Methods

### NewListAPIKeysOutputBody

`func NewListAPIKeysOutputBody(items []ApiKeyView, ) *ListAPIKeysOutputBody`

NewListAPIKeysOutputBody instantiates a new ListAPIKeysOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListAPIKeysOutputBodyWithDefaults

`func NewListAPIKeysOutputBodyWithDefaults() *ListAPIKeysOutputBody`

NewListAPIKeysOutputBodyWithDefaults instantiates a new ListAPIKeysOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *ListAPIKeysOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *ListAPIKeysOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *ListAPIKeysOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *ListAPIKeysOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetItems

`func (o *ListAPIKeysOutputBody) GetItems() []ApiKeyView`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ListAPIKeysOutputBody) GetItemsOk() (*[]ApiKeyView, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ListAPIKeysOutputBody) SetItems(v []ApiKeyView)`

SetItems sets Items field to given value.


### SetItemsNil

`func (o *ListAPIKeysOutputBody) SetItemsNil(b bool)`

 SetItemsNil sets the value for Items to be an explicit nil

### UnsetItems
`func (o *ListAPIKeysOutputBody) UnsetItems()`

UnsetItems ensures that no value is present for Items, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


