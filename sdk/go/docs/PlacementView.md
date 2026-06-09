# PlacementView

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**CreatedAt** | **time.Time** |  | 
**HypervisorId** | **string** |  | 
**HypervisorName** | **string** |  | 
**Id** | **string** |  | 
**Vmid** | Pointer to **int64** |  | [optional] 

## Methods

### NewPlacementView

`func NewPlacementView(createdAt time.Time, hypervisorId string, hypervisorName string, id string, ) *PlacementView`

NewPlacementView instantiates a new PlacementView object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlacementViewWithDefaults

`func NewPlacementViewWithDefaults() *PlacementView`

NewPlacementViewWithDefaults instantiates a new PlacementView object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *PlacementView) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *PlacementView) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *PlacementView) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *PlacementView) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetCreatedAt

`func (o *PlacementView) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *PlacementView) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *PlacementView) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetHypervisorId

`func (o *PlacementView) GetHypervisorId() string`

GetHypervisorId returns the HypervisorId field if non-nil, zero value otherwise.

### GetHypervisorIdOk

`func (o *PlacementView) GetHypervisorIdOk() (*string, bool)`

GetHypervisorIdOk returns a tuple with the HypervisorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHypervisorId

`func (o *PlacementView) SetHypervisorId(v string)`

SetHypervisorId sets HypervisorId field to given value.


### GetHypervisorName

`func (o *PlacementView) GetHypervisorName() string`

GetHypervisorName returns the HypervisorName field if non-nil, zero value otherwise.

### GetHypervisorNameOk

`func (o *PlacementView) GetHypervisorNameOk() (*string, bool)`

GetHypervisorNameOk returns a tuple with the HypervisorName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHypervisorName

`func (o *PlacementView) SetHypervisorName(v string)`

SetHypervisorName sets HypervisorName field to given value.


### GetId

`func (o *PlacementView) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PlacementView) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PlacementView) SetId(v string)`

SetId sets Id field to given value.


### GetVmid

`func (o *PlacementView) GetVmid() int64`

GetVmid returns the Vmid field if non-nil, zero value otherwise.

### GetVmidOk

`func (o *PlacementView) GetVmidOk() (*int64, bool)`

GetVmidOk returns a tuple with the Vmid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVmid

`func (o *PlacementView) SetVmid(v int64)`

SetVmid sets Vmid field to given value.

### HasVmid

`func (o *PlacementView) HasVmid() bool`

HasVmid returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


