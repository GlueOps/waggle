# HealthOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Status** | **string** |  | 

## Methods

### NewHealthOutputBody

`func NewHealthOutputBody(status string, ) *HealthOutputBody`

NewHealthOutputBody instantiates a new HealthOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHealthOutputBodyWithDefaults

`func NewHealthOutputBodyWithDefaults() *HealthOutputBody`

NewHealthOutputBodyWithDefaults instantiates a new HealthOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *HealthOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *HealthOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *HealthOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *HealthOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetStatus

`func (o *HealthOutputBody) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *HealthOutputBody) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *HealthOutputBody) SetStatus(v string)`

SetStatus sets Status field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


