# FleetPlacementView

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreatedAt** | **time.Time** |  | 
**DiskGb** | **int64** |  | 
**HypervisorId** | **string** |  | 
**HypervisorName** | **string** |  | 
**Id** | **string** |  | 
**PoolId** | **string** |  | 
**PoolName** | **string** |  | 
**RamGb** | **int64** |  | 
**SlotName** | **string** |  | 
**Vcpu** | **int64** |  | 
**Vmid** | Pointer to **int64** |  | [optional] 

## Methods

### NewFleetPlacementView

`func NewFleetPlacementView(createdAt time.Time, diskGb int64, hypervisorId string, hypervisorName string, id string, poolId string, poolName string, ramGb int64, slotName string, vcpu int64, ) *FleetPlacementView`

NewFleetPlacementView instantiates a new FleetPlacementView object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewFleetPlacementViewWithDefaults

`func NewFleetPlacementViewWithDefaults() *FleetPlacementView`

NewFleetPlacementViewWithDefaults instantiates a new FleetPlacementView object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedAt

`func (o *FleetPlacementView) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *FleetPlacementView) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *FleetPlacementView) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetDiskGb

`func (o *FleetPlacementView) GetDiskGb() int64`

GetDiskGb returns the DiskGb field if non-nil, zero value otherwise.

### GetDiskGbOk

`func (o *FleetPlacementView) GetDiskGbOk() (*int64, bool)`

GetDiskGbOk returns a tuple with the DiskGb field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiskGb

`func (o *FleetPlacementView) SetDiskGb(v int64)`

SetDiskGb sets DiskGb field to given value.


### GetHypervisorId

`func (o *FleetPlacementView) GetHypervisorId() string`

GetHypervisorId returns the HypervisorId field if non-nil, zero value otherwise.

### GetHypervisorIdOk

`func (o *FleetPlacementView) GetHypervisorIdOk() (*string, bool)`

GetHypervisorIdOk returns a tuple with the HypervisorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHypervisorId

`func (o *FleetPlacementView) SetHypervisorId(v string)`

SetHypervisorId sets HypervisorId field to given value.


### GetHypervisorName

`func (o *FleetPlacementView) GetHypervisorName() string`

GetHypervisorName returns the HypervisorName field if non-nil, zero value otherwise.

### GetHypervisorNameOk

`func (o *FleetPlacementView) GetHypervisorNameOk() (*string, bool)`

GetHypervisorNameOk returns a tuple with the HypervisorName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHypervisorName

`func (o *FleetPlacementView) SetHypervisorName(v string)`

SetHypervisorName sets HypervisorName field to given value.


### GetId

`func (o *FleetPlacementView) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *FleetPlacementView) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *FleetPlacementView) SetId(v string)`

SetId sets Id field to given value.


### GetPoolId

`func (o *FleetPlacementView) GetPoolId() string`

GetPoolId returns the PoolId field if non-nil, zero value otherwise.

### GetPoolIdOk

`func (o *FleetPlacementView) GetPoolIdOk() (*string, bool)`

GetPoolIdOk returns a tuple with the PoolId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolId

`func (o *FleetPlacementView) SetPoolId(v string)`

SetPoolId sets PoolId field to given value.


### GetPoolName

`func (o *FleetPlacementView) GetPoolName() string`

GetPoolName returns the PoolName field if non-nil, zero value otherwise.

### GetPoolNameOk

`func (o *FleetPlacementView) GetPoolNameOk() (*string, bool)`

GetPoolNameOk returns a tuple with the PoolName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolName

`func (o *FleetPlacementView) SetPoolName(v string)`

SetPoolName sets PoolName field to given value.


### GetRamGb

`func (o *FleetPlacementView) GetRamGb() int64`

GetRamGb returns the RamGb field if non-nil, zero value otherwise.

### GetRamGbOk

`func (o *FleetPlacementView) GetRamGbOk() (*int64, bool)`

GetRamGbOk returns a tuple with the RamGb field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRamGb

`func (o *FleetPlacementView) SetRamGb(v int64)`

SetRamGb sets RamGb field to given value.


### GetSlotName

`func (o *FleetPlacementView) GetSlotName() string`

GetSlotName returns the SlotName field if non-nil, zero value otherwise.

### GetSlotNameOk

`func (o *FleetPlacementView) GetSlotNameOk() (*string, bool)`

GetSlotNameOk returns a tuple with the SlotName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSlotName

`func (o *FleetPlacementView) SetSlotName(v string)`

SetSlotName sets SlotName field to given value.


### GetVcpu

`func (o *FleetPlacementView) GetVcpu() int64`

GetVcpu returns the Vcpu field if non-nil, zero value otherwise.

### GetVcpuOk

`func (o *FleetPlacementView) GetVcpuOk() (*int64, bool)`

GetVcpuOk returns a tuple with the Vcpu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVcpu

`func (o *FleetPlacementView) SetVcpu(v int64)`

SetVcpu sets Vcpu field to given value.


### GetVmid

`func (o *FleetPlacementView) GetVmid() int64`

GetVmid returns the Vmid field if non-nil, zero value otherwise.

### GetVmidOk

`func (o *FleetPlacementView) GetVmidOk() (*int64, bool)`

GetVmidOk returns a tuple with the Vmid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVmid

`func (o *FleetPlacementView) SetVmid(v int64)`

SetVmid sets Vmid field to given value.

### HasVmid

`func (o *FleetPlacementView) HasVmid() bool`

HasVmid returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


