# HypervisorView

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**CpuBookable** | **int64** |  | 
**CpuReserved** | **int64** |  | 
**CpuTotal** | **int64** |  | 
**CpuUsed** | **int64** | vCPU allocated to existing guests (from discovery). | 
**CreatedAt** | **time.Time** |  | 
**DatacenterId** | **string** |  | 
**DiskGbBookable** | **int64** |  | 
**DiskGbReserved** | **int64** |  | 
**DiskGbTotal** | **int64** |  | 
**DiskGbUsed** | **int64** | Disk (GB) allocated to existing guests (from discovery). | 
**Id** | **string** |  | 
**LastSyncedAt** | Pointer to **time.Time** |  | [optional] 
**Name** | **string** |  | 
**RamGbBookable** | **int64** |  | 
**RamGbReserved** | **int64** |  | 
**RamGbTotal** | **int64** |  | 
**RamGbUsed** | **int64** | RAM (GB) allocated to existing guests (from discovery). | 
**Schedulable** | **bool** | When false, placement excludes this hypervisor. | 
**UpdatedAt** | **time.Time** |  | 

## Methods

### NewHypervisorView

`func NewHypervisorView(cpuBookable int64, cpuReserved int64, cpuTotal int64, cpuUsed int64, createdAt time.Time, datacenterId string, diskGbBookable int64, diskGbReserved int64, diskGbTotal int64, diskGbUsed int64, id string, name string, ramGbBookable int64, ramGbReserved int64, ramGbTotal int64, ramGbUsed int64, schedulable bool, updatedAt time.Time, ) *HypervisorView`

NewHypervisorView instantiates a new HypervisorView object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHypervisorViewWithDefaults

`func NewHypervisorViewWithDefaults() *HypervisorView`

NewHypervisorViewWithDefaults instantiates a new HypervisorView object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *HypervisorView) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *HypervisorView) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *HypervisorView) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *HypervisorView) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetCpuBookable

`func (o *HypervisorView) GetCpuBookable() int64`

GetCpuBookable returns the CpuBookable field if non-nil, zero value otherwise.

### GetCpuBookableOk

`func (o *HypervisorView) GetCpuBookableOk() (*int64, bool)`

GetCpuBookableOk returns a tuple with the CpuBookable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpuBookable

`func (o *HypervisorView) SetCpuBookable(v int64)`

SetCpuBookable sets CpuBookable field to given value.


### GetCpuReserved

`func (o *HypervisorView) GetCpuReserved() int64`

GetCpuReserved returns the CpuReserved field if non-nil, zero value otherwise.

### GetCpuReservedOk

`func (o *HypervisorView) GetCpuReservedOk() (*int64, bool)`

GetCpuReservedOk returns a tuple with the CpuReserved field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpuReserved

`func (o *HypervisorView) SetCpuReserved(v int64)`

SetCpuReserved sets CpuReserved field to given value.


### GetCpuTotal

`func (o *HypervisorView) GetCpuTotal() int64`

GetCpuTotal returns the CpuTotal field if non-nil, zero value otherwise.

### GetCpuTotalOk

`func (o *HypervisorView) GetCpuTotalOk() (*int64, bool)`

GetCpuTotalOk returns a tuple with the CpuTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpuTotal

`func (o *HypervisorView) SetCpuTotal(v int64)`

SetCpuTotal sets CpuTotal field to given value.


### GetCpuUsed

`func (o *HypervisorView) GetCpuUsed() int64`

GetCpuUsed returns the CpuUsed field if non-nil, zero value otherwise.

### GetCpuUsedOk

`func (o *HypervisorView) GetCpuUsedOk() (*int64, bool)`

GetCpuUsedOk returns a tuple with the CpuUsed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpuUsed

`func (o *HypervisorView) SetCpuUsed(v int64)`

SetCpuUsed sets CpuUsed field to given value.


### GetCreatedAt

`func (o *HypervisorView) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *HypervisorView) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *HypervisorView) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetDatacenterId

`func (o *HypervisorView) GetDatacenterId() string`

GetDatacenterId returns the DatacenterId field if non-nil, zero value otherwise.

### GetDatacenterIdOk

`func (o *HypervisorView) GetDatacenterIdOk() (*string, bool)`

GetDatacenterIdOk returns a tuple with the DatacenterId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDatacenterId

`func (o *HypervisorView) SetDatacenterId(v string)`

SetDatacenterId sets DatacenterId field to given value.


### GetDiskGbBookable

`func (o *HypervisorView) GetDiskGbBookable() int64`

GetDiskGbBookable returns the DiskGbBookable field if non-nil, zero value otherwise.

### GetDiskGbBookableOk

`func (o *HypervisorView) GetDiskGbBookableOk() (*int64, bool)`

GetDiskGbBookableOk returns a tuple with the DiskGbBookable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiskGbBookable

`func (o *HypervisorView) SetDiskGbBookable(v int64)`

SetDiskGbBookable sets DiskGbBookable field to given value.


### GetDiskGbReserved

`func (o *HypervisorView) GetDiskGbReserved() int64`

GetDiskGbReserved returns the DiskGbReserved field if non-nil, zero value otherwise.

### GetDiskGbReservedOk

`func (o *HypervisorView) GetDiskGbReservedOk() (*int64, bool)`

GetDiskGbReservedOk returns a tuple with the DiskGbReserved field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiskGbReserved

`func (o *HypervisorView) SetDiskGbReserved(v int64)`

SetDiskGbReserved sets DiskGbReserved field to given value.


### GetDiskGbTotal

`func (o *HypervisorView) GetDiskGbTotal() int64`

GetDiskGbTotal returns the DiskGbTotal field if non-nil, zero value otherwise.

### GetDiskGbTotalOk

`func (o *HypervisorView) GetDiskGbTotalOk() (*int64, bool)`

GetDiskGbTotalOk returns a tuple with the DiskGbTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiskGbTotal

`func (o *HypervisorView) SetDiskGbTotal(v int64)`

SetDiskGbTotal sets DiskGbTotal field to given value.


### GetDiskGbUsed

`func (o *HypervisorView) GetDiskGbUsed() int64`

GetDiskGbUsed returns the DiskGbUsed field if non-nil, zero value otherwise.

### GetDiskGbUsedOk

`func (o *HypervisorView) GetDiskGbUsedOk() (*int64, bool)`

GetDiskGbUsedOk returns a tuple with the DiskGbUsed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiskGbUsed

`func (o *HypervisorView) SetDiskGbUsed(v int64)`

SetDiskGbUsed sets DiskGbUsed field to given value.


### GetId

`func (o *HypervisorView) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *HypervisorView) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *HypervisorView) SetId(v string)`

SetId sets Id field to given value.


### GetLastSyncedAt

`func (o *HypervisorView) GetLastSyncedAt() time.Time`

GetLastSyncedAt returns the LastSyncedAt field if non-nil, zero value otherwise.

### GetLastSyncedAtOk

`func (o *HypervisorView) GetLastSyncedAtOk() (*time.Time, bool)`

GetLastSyncedAtOk returns a tuple with the LastSyncedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastSyncedAt

`func (o *HypervisorView) SetLastSyncedAt(v time.Time)`

SetLastSyncedAt sets LastSyncedAt field to given value.

### HasLastSyncedAt

`func (o *HypervisorView) HasLastSyncedAt() bool`

HasLastSyncedAt returns a boolean if a field has been set.

### GetName

`func (o *HypervisorView) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *HypervisorView) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *HypervisorView) SetName(v string)`

SetName sets Name field to given value.


### GetRamGbBookable

`func (o *HypervisorView) GetRamGbBookable() int64`

GetRamGbBookable returns the RamGbBookable field if non-nil, zero value otherwise.

### GetRamGbBookableOk

`func (o *HypervisorView) GetRamGbBookableOk() (*int64, bool)`

GetRamGbBookableOk returns a tuple with the RamGbBookable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRamGbBookable

`func (o *HypervisorView) SetRamGbBookable(v int64)`

SetRamGbBookable sets RamGbBookable field to given value.


### GetRamGbReserved

`func (o *HypervisorView) GetRamGbReserved() int64`

GetRamGbReserved returns the RamGbReserved field if non-nil, zero value otherwise.

### GetRamGbReservedOk

`func (o *HypervisorView) GetRamGbReservedOk() (*int64, bool)`

GetRamGbReservedOk returns a tuple with the RamGbReserved field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRamGbReserved

`func (o *HypervisorView) SetRamGbReserved(v int64)`

SetRamGbReserved sets RamGbReserved field to given value.


### GetRamGbTotal

`func (o *HypervisorView) GetRamGbTotal() int64`

GetRamGbTotal returns the RamGbTotal field if non-nil, zero value otherwise.

### GetRamGbTotalOk

`func (o *HypervisorView) GetRamGbTotalOk() (*int64, bool)`

GetRamGbTotalOk returns a tuple with the RamGbTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRamGbTotal

`func (o *HypervisorView) SetRamGbTotal(v int64)`

SetRamGbTotal sets RamGbTotal field to given value.


### GetRamGbUsed

`func (o *HypervisorView) GetRamGbUsed() int64`

GetRamGbUsed returns the RamGbUsed field if non-nil, zero value otherwise.

### GetRamGbUsedOk

`func (o *HypervisorView) GetRamGbUsedOk() (*int64, bool)`

GetRamGbUsedOk returns a tuple with the RamGbUsed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRamGbUsed

`func (o *HypervisorView) SetRamGbUsed(v int64)`

SetRamGbUsed sets RamGbUsed field to given value.


### GetSchedulable

`func (o *HypervisorView) GetSchedulable() bool`

GetSchedulable returns the Schedulable field if non-nil, zero value otherwise.

### GetSchedulableOk

`func (o *HypervisorView) GetSchedulableOk() (*bool, bool)`

GetSchedulableOk returns a tuple with the Schedulable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchedulable

`func (o *HypervisorView) SetSchedulable(v bool)`

SetSchedulable sets Schedulable field to given value.


### GetUpdatedAt

`func (o *HypervisorView) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *HypervisorView) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *HypervisorView) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


