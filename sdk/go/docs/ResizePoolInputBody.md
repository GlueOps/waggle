# ResizePoolInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**DesiredCount** | **int64** |  | 

## Methods

### NewResizePoolInputBody

`func NewResizePoolInputBody(desiredCount int64, ) *ResizePoolInputBody`

NewResizePoolInputBody instantiates a new ResizePoolInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewResizePoolInputBodyWithDefaults

`func NewResizePoolInputBodyWithDefaults() *ResizePoolInputBody`

NewResizePoolInputBodyWithDefaults instantiates a new ResizePoolInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *ResizePoolInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *ResizePoolInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *ResizePoolInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *ResizePoolInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetDesiredCount

`func (o *ResizePoolInputBody) GetDesiredCount() int64`

GetDesiredCount returns the DesiredCount field if non-nil, zero value otherwise.

### GetDesiredCountOk

`func (o *ResizePoolInputBody) GetDesiredCountOk() (*int64, bool)`

GetDesiredCountOk returns a tuple with the DesiredCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDesiredCount

`func (o *ResizePoolInputBody) SetDesiredCount(v int64)`

SetDesiredCount sets DesiredCount field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


