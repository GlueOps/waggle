# CreateAPIKeyInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**ExpiresInDays** | Pointer to **int64** | Days until the key expires; omit or 0 for no expiry. | [optional] 
**Name** | **string** | Human label for the key. | 

## Methods

### NewCreateAPIKeyInputBody

`func NewCreateAPIKeyInputBody(name string, ) *CreateAPIKeyInputBody`

NewCreateAPIKeyInputBody instantiates a new CreateAPIKeyInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateAPIKeyInputBodyWithDefaults

`func NewCreateAPIKeyInputBodyWithDefaults() *CreateAPIKeyInputBody`

NewCreateAPIKeyInputBodyWithDefaults instantiates a new CreateAPIKeyInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *CreateAPIKeyInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *CreateAPIKeyInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *CreateAPIKeyInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *CreateAPIKeyInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetExpiresInDays

`func (o *CreateAPIKeyInputBody) GetExpiresInDays() int64`

GetExpiresInDays returns the ExpiresInDays field if non-nil, zero value otherwise.

### GetExpiresInDaysOk

`func (o *CreateAPIKeyInputBody) GetExpiresInDaysOk() (*int64, bool)`

GetExpiresInDaysOk returns a tuple with the ExpiresInDays field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresInDays

`func (o *CreateAPIKeyInputBody) SetExpiresInDays(v int64)`

SetExpiresInDays sets ExpiresInDays field to given value.

### HasExpiresInDays

`func (o *CreateAPIKeyInputBody) HasExpiresInDays() bool`

HasExpiresInDays returns a boolean if a field has been set.

### GetName

`func (o *CreateAPIKeyInputBody) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateAPIKeyInputBody) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateAPIKeyInputBody) SetName(v string)`

SetName sets Name field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


