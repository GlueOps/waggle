# PoolResultBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Placements** | [**[]PlacementView**](PlacementView.md) |  | 
**Pool** | [**PoolView**](PoolView.md) |  | 

## Methods

### NewPoolResultBody

`func NewPoolResultBody(placements []PlacementView, pool PoolView, ) *PoolResultBody`

NewPoolResultBody instantiates a new PoolResultBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPoolResultBodyWithDefaults

`func NewPoolResultBodyWithDefaults() *PoolResultBody`

NewPoolResultBodyWithDefaults instantiates a new PoolResultBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *PoolResultBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *PoolResultBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *PoolResultBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *PoolResultBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetPlacements

`func (o *PoolResultBody) GetPlacements() []PlacementView`

GetPlacements returns the Placements field if non-nil, zero value otherwise.

### GetPlacementsOk

`func (o *PoolResultBody) GetPlacementsOk() (*[]PlacementView, bool)`

GetPlacementsOk returns a tuple with the Placements field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlacements

`func (o *PoolResultBody) SetPlacements(v []PlacementView)`

SetPlacements sets Placements field to given value.


### SetPlacementsNil

`func (o *PoolResultBody) SetPlacementsNil(b bool)`

 SetPlacementsNil sets the value for Placements to be an explicit nil

### UnsetPlacements
`func (o *PoolResultBody) UnsetPlacements()`

UnsetPlacements ensures that no value is present for Placements, not even an explicit nil
### GetPool

`func (o *PoolResultBody) GetPool() PoolView`

GetPool returns the Pool field if non-nil, zero value otherwise.

### GetPoolOk

`func (o *PoolResultBody) GetPoolOk() (*PoolView, bool)`

GetPoolOk returns a tuple with the Pool field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPool

`func (o *PoolResultBody) SetPool(v PoolView)`

SetPool sets Pool field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


