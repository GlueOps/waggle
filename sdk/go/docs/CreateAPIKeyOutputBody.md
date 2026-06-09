# CreateAPIKeyOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Key** | [**ApiKeyView**](ApiKeyView.md) |  | 
**Token** | **string** | Plaintext API key — shown once, store it securely. | 

## Methods

### NewCreateAPIKeyOutputBody

`func NewCreateAPIKeyOutputBody(key ApiKeyView, token string, ) *CreateAPIKeyOutputBody`

NewCreateAPIKeyOutputBody instantiates a new CreateAPIKeyOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateAPIKeyOutputBodyWithDefaults

`func NewCreateAPIKeyOutputBodyWithDefaults() *CreateAPIKeyOutputBody`

NewCreateAPIKeyOutputBodyWithDefaults instantiates a new CreateAPIKeyOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *CreateAPIKeyOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *CreateAPIKeyOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *CreateAPIKeyOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *CreateAPIKeyOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetKey

`func (o *CreateAPIKeyOutputBody) GetKey() ApiKeyView`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *CreateAPIKeyOutputBody) GetKeyOk() (*ApiKeyView, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *CreateAPIKeyOutputBody) SetKey(v ApiKeyView)`

SetKey sets Key field to given value.


### GetToken

`func (o *CreateAPIKeyOutputBody) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *CreateAPIKeyOutputBody) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *CreateAPIKeyOutputBody) SetToken(v string)`

SetToken sets Token field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


