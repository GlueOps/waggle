# BackfillVMIDInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Vmid** | **int64** |  | 

## Methods

### NewBackfillVMIDInputBody

`func NewBackfillVMIDInputBody(vmid int64, ) *BackfillVMIDInputBody`

NewBackfillVMIDInputBody instantiates a new BackfillVMIDInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBackfillVMIDInputBodyWithDefaults

`func NewBackfillVMIDInputBodyWithDefaults() *BackfillVMIDInputBody`

NewBackfillVMIDInputBodyWithDefaults instantiates a new BackfillVMIDInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *BackfillVMIDInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *BackfillVMIDInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *BackfillVMIDInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *BackfillVMIDInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetVmid

`func (o *BackfillVMIDInputBody) GetVmid() int64`

GetVmid returns the Vmid field if non-nil, zero value otherwise.

### GetVmidOk

`func (o *BackfillVMIDInputBody) GetVmidOk() (*int64, bool)`

GetVmidOk returns a tuple with the Vmid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVmid

`func (o *BackfillVMIDInputBody) SetVmid(v int64)`

SetVmid sets Vmid field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


