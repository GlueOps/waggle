# SlotBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**DiskGb** | **int64** |  | 
**Name** | **string** |  | 
**RamGb** | **int64** |  | 
**Vcpu** | **int64** |  | 

## Methods

### NewSlotBody

`func NewSlotBody(diskGb int64, name string, ramGb int64, vcpu int64, ) *SlotBody`

NewSlotBody instantiates a new SlotBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSlotBodyWithDefaults

`func NewSlotBodyWithDefaults() *SlotBody`

NewSlotBodyWithDefaults instantiates a new SlotBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *SlotBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *SlotBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *SlotBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *SlotBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetDiskGb

`func (o *SlotBody) GetDiskGb() int64`

GetDiskGb returns the DiskGb field if non-nil, zero value otherwise.

### GetDiskGbOk

`func (o *SlotBody) GetDiskGbOk() (*int64, bool)`

GetDiskGbOk returns a tuple with the DiskGb field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiskGb

`func (o *SlotBody) SetDiskGb(v int64)`

SetDiskGb sets DiskGb field to given value.


### GetName

`func (o *SlotBody) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SlotBody) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SlotBody) SetName(v string)`

SetName sets Name field to given value.


### GetRamGb

`func (o *SlotBody) GetRamGb() int64`

GetRamGb returns the RamGb field if non-nil, zero value otherwise.

### GetRamGbOk

`func (o *SlotBody) GetRamGbOk() (*int64, bool)`

GetRamGbOk returns a tuple with the RamGb field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRamGb

`func (o *SlotBody) SetRamGb(v int64)`

SetRamGb sets RamGb field to given value.


### GetVcpu

`func (o *SlotBody) GetVcpu() int64`

GetVcpu returns the Vcpu field if non-nil, zero value otherwise.

### GetVcpuOk

`func (o *SlotBody) GetVcpuOk() (*int64, bool)`

GetVcpuOk returns a tuple with the Vcpu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVcpu

`func (o *SlotBody) SetVcpu(v int64)`

SetVcpu sets Vcpu field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


