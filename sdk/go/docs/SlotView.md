# SlotView

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**CreatedAt** | **time.Time** |  | 
**DiskGb** | **int64** |  | 
**Id** | **string** |  | 
**Name** | **string** |  | 
**RamGb** | **int64** |  | 
**UpdatedAt** | **time.Time** |  | 
**Vcpu** | **int64** |  | 

## Methods

### NewSlotView

`func NewSlotView(createdAt time.Time, diskGb int64, id string, name string, ramGb int64, updatedAt time.Time, vcpu int64, ) *SlotView`

NewSlotView instantiates a new SlotView object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSlotViewWithDefaults

`func NewSlotViewWithDefaults() *SlotView`

NewSlotViewWithDefaults instantiates a new SlotView object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *SlotView) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *SlotView) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *SlotView) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *SlotView) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetCreatedAt

`func (o *SlotView) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *SlotView) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *SlotView) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetDiskGb

`func (o *SlotView) GetDiskGb() int64`

GetDiskGb returns the DiskGb field if non-nil, zero value otherwise.

### GetDiskGbOk

`func (o *SlotView) GetDiskGbOk() (*int64, bool)`

GetDiskGbOk returns a tuple with the DiskGb field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiskGb

`func (o *SlotView) SetDiskGb(v int64)`

SetDiskGb sets DiskGb field to given value.


### GetId

`func (o *SlotView) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *SlotView) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *SlotView) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *SlotView) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SlotView) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SlotView) SetName(v string)`

SetName sets Name field to given value.


### GetRamGb

`func (o *SlotView) GetRamGb() int64`

GetRamGb returns the RamGb field if non-nil, zero value otherwise.

### GetRamGbOk

`func (o *SlotView) GetRamGbOk() (*int64, bool)`

GetRamGbOk returns a tuple with the RamGb field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRamGb

`func (o *SlotView) SetRamGb(v int64)`

SetRamGb sets RamGb field to given value.


### GetUpdatedAt

`func (o *SlotView) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *SlotView) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *SlotView) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetVcpu

`func (o *SlotView) GetVcpu() int64`

GetVcpu returns the Vcpu field if non-nil, zero value otherwise.

### GetVcpuOk

`func (o *SlotView) GetVcpuOk() (*int64, bool)`

GetVcpuOk returns a tuple with the Vcpu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVcpu

`func (o *SlotView) SetVcpu(v int64)`

SetVcpu sets Vcpu field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


