# HypervisorBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**CpuReserved** | **int64** |  | 
**CpuTotal** | **int64** |  | 
**DatacenterId** | **string** |  | 
**DiskGbReserved** | **int64** |  | 
**DiskGbTotal** | **int64** |  | 
**Name** | **string** |  | 
**RamGbReserved** | **int64** |  | 
**RamGbTotal** | **int64** |  | 
**Schedulable** | Pointer to **bool** |  | [optional] 

## Methods

### NewHypervisorBody

`func NewHypervisorBody(cpuReserved int64, cpuTotal int64, datacenterId string, diskGbReserved int64, diskGbTotal int64, name string, ramGbReserved int64, ramGbTotal int64, ) *HypervisorBody`

NewHypervisorBody instantiates a new HypervisorBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHypervisorBodyWithDefaults

`func NewHypervisorBodyWithDefaults() *HypervisorBody`

NewHypervisorBodyWithDefaults instantiates a new HypervisorBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *HypervisorBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *HypervisorBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *HypervisorBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *HypervisorBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetCpuReserved

`func (o *HypervisorBody) GetCpuReserved() int64`

GetCpuReserved returns the CpuReserved field if non-nil, zero value otherwise.

### GetCpuReservedOk

`func (o *HypervisorBody) GetCpuReservedOk() (*int64, bool)`

GetCpuReservedOk returns a tuple with the CpuReserved field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpuReserved

`func (o *HypervisorBody) SetCpuReserved(v int64)`

SetCpuReserved sets CpuReserved field to given value.


### GetCpuTotal

`func (o *HypervisorBody) GetCpuTotal() int64`

GetCpuTotal returns the CpuTotal field if non-nil, zero value otherwise.

### GetCpuTotalOk

`func (o *HypervisorBody) GetCpuTotalOk() (*int64, bool)`

GetCpuTotalOk returns a tuple with the CpuTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpuTotal

`func (o *HypervisorBody) SetCpuTotal(v int64)`

SetCpuTotal sets CpuTotal field to given value.


### GetDatacenterId

`func (o *HypervisorBody) GetDatacenterId() string`

GetDatacenterId returns the DatacenterId field if non-nil, zero value otherwise.

### GetDatacenterIdOk

`func (o *HypervisorBody) GetDatacenterIdOk() (*string, bool)`

GetDatacenterIdOk returns a tuple with the DatacenterId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDatacenterId

`func (o *HypervisorBody) SetDatacenterId(v string)`

SetDatacenterId sets DatacenterId field to given value.


### GetDiskGbReserved

`func (o *HypervisorBody) GetDiskGbReserved() int64`

GetDiskGbReserved returns the DiskGbReserved field if non-nil, zero value otherwise.

### GetDiskGbReservedOk

`func (o *HypervisorBody) GetDiskGbReservedOk() (*int64, bool)`

GetDiskGbReservedOk returns a tuple with the DiskGbReserved field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiskGbReserved

`func (o *HypervisorBody) SetDiskGbReserved(v int64)`

SetDiskGbReserved sets DiskGbReserved field to given value.


### GetDiskGbTotal

`func (o *HypervisorBody) GetDiskGbTotal() int64`

GetDiskGbTotal returns the DiskGbTotal field if non-nil, zero value otherwise.

### GetDiskGbTotalOk

`func (o *HypervisorBody) GetDiskGbTotalOk() (*int64, bool)`

GetDiskGbTotalOk returns a tuple with the DiskGbTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiskGbTotal

`func (o *HypervisorBody) SetDiskGbTotal(v int64)`

SetDiskGbTotal sets DiskGbTotal field to given value.


### GetName

`func (o *HypervisorBody) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *HypervisorBody) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *HypervisorBody) SetName(v string)`

SetName sets Name field to given value.


### GetRamGbReserved

`func (o *HypervisorBody) GetRamGbReserved() int64`

GetRamGbReserved returns the RamGbReserved field if non-nil, zero value otherwise.

### GetRamGbReservedOk

`func (o *HypervisorBody) GetRamGbReservedOk() (*int64, bool)`

GetRamGbReservedOk returns a tuple with the RamGbReserved field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRamGbReserved

`func (o *HypervisorBody) SetRamGbReserved(v int64)`

SetRamGbReserved sets RamGbReserved field to given value.


### GetRamGbTotal

`func (o *HypervisorBody) GetRamGbTotal() int64`

GetRamGbTotal returns the RamGbTotal field if non-nil, zero value otherwise.

### GetRamGbTotalOk

`func (o *HypervisorBody) GetRamGbTotalOk() (*int64, bool)`

GetRamGbTotalOk returns a tuple with the RamGbTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRamGbTotal

`func (o *HypervisorBody) SetRamGbTotal(v int64)`

SetRamGbTotal sets RamGbTotal field to given value.


### GetSchedulable

`func (o *HypervisorBody) GetSchedulable() bool`

GetSchedulable returns the Schedulable field if non-nil, zero value otherwise.

### GetSchedulableOk

`func (o *HypervisorBody) GetSchedulableOk() (*bool, bool)`

GetSchedulableOk returns a tuple with the Schedulable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchedulable

`func (o *HypervisorBody) SetSchedulable(v bool)`

SetSchedulable sets Schedulable field to given value.

### HasSchedulable

`func (o *HypervisorBody) HasSchedulable() bool`

HasSchedulable returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


