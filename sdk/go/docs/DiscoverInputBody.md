# DiscoverInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Async** | Pointer to **bool** |  | [optional] 

## Methods

### NewDiscoverInputBody

`func NewDiscoverInputBody() *DiscoverInputBody`

NewDiscoverInputBody instantiates a new DiscoverInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDiscoverInputBodyWithDefaults

`func NewDiscoverInputBodyWithDefaults() *DiscoverInputBody`

NewDiscoverInputBodyWithDefaults instantiates a new DiscoverInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *DiscoverInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *DiscoverInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *DiscoverInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *DiscoverInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetAsync

`func (o *DiscoverInputBody) GetAsync() bool`

GetAsync returns the Async field if non-nil, zero value otherwise.

### GetAsyncOk

`func (o *DiscoverInputBody) GetAsyncOk() (*bool, bool)`

GetAsyncOk returns a tuple with the Async field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsync

`func (o *DiscoverInputBody) SetAsync(v bool)`

SetAsync sets Async field to given value.

### HasAsync

`func (o *DiscoverInputBody) HasAsync() bool`

HasAsync returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


