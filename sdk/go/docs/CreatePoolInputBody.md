# CreatePoolInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**DatacenterId** | **string** |  | 
**DesiredCount** | **int64** |  | 
**Metadata** | Pointer to **interface{}** |  | [optional] 
**Name** | **string** |  | 
**SlotId** | **string** |  | 

## Methods

### NewCreatePoolInputBody

`func NewCreatePoolInputBody(datacenterId string, desiredCount int64, name string, slotId string, ) *CreatePoolInputBody`

NewCreatePoolInputBody instantiates a new CreatePoolInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreatePoolInputBodyWithDefaults

`func NewCreatePoolInputBodyWithDefaults() *CreatePoolInputBody`

NewCreatePoolInputBodyWithDefaults instantiates a new CreatePoolInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *CreatePoolInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *CreatePoolInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *CreatePoolInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *CreatePoolInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetDatacenterId

`func (o *CreatePoolInputBody) GetDatacenterId() string`

GetDatacenterId returns the DatacenterId field if non-nil, zero value otherwise.

### GetDatacenterIdOk

`func (o *CreatePoolInputBody) GetDatacenterIdOk() (*string, bool)`

GetDatacenterIdOk returns a tuple with the DatacenterId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDatacenterId

`func (o *CreatePoolInputBody) SetDatacenterId(v string)`

SetDatacenterId sets DatacenterId field to given value.


### GetDesiredCount

`func (o *CreatePoolInputBody) GetDesiredCount() int64`

GetDesiredCount returns the DesiredCount field if non-nil, zero value otherwise.

### GetDesiredCountOk

`func (o *CreatePoolInputBody) GetDesiredCountOk() (*int64, bool)`

GetDesiredCountOk returns a tuple with the DesiredCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDesiredCount

`func (o *CreatePoolInputBody) SetDesiredCount(v int64)`

SetDesiredCount sets DesiredCount field to given value.


### GetMetadata

`func (o *CreatePoolInputBody) GetMetadata() interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *CreatePoolInputBody) GetMetadataOk() (*interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *CreatePoolInputBody) SetMetadata(v interface{})`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *CreatePoolInputBody) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### SetMetadataNil

`func (o *CreatePoolInputBody) SetMetadataNil(b bool)`

 SetMetadataNil sets the value for Metadata to be an explicit nil

### UnsetMetadata
`func (o *CreatePoolInputBody) UnsetMetadata()`

UnsetMetadata ensures that no value is present for Metadata, not even an explicit nil
### GetName

`func (o *CreatePoolInputBody) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreatePoolInputBody) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreatePoolInputBody) SetName(v string)`

SetName sets Name field to given value.


### GetSlotId

`func (o *CreatePoolInputBody) GetSlotId() string`

GetSlotId returns the SlotId field if non-nil, zero value otherwise.

### GetSlotIdOk

`func (o *CreatePoolInputBody) GetSlotIdOk() (*string, bool)`

GetSlotIdOk returns a tuple with the SlotId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSlotId

`func (o *CreatePoolInputBody) SetSlotId(v string)`

SetSlotId sets SlotId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


