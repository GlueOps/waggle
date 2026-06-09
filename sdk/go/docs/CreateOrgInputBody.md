# CreateOrgInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Name** | **string** |  | 

## Methods

### NewCreateOrgInputBody

`func NewCreateOrgInputBody(name string, ) *CreateOrgInputBody`

NewCreateOrgInputBody instantiates a new CreateOrgInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateOrgInputBodyWithDefaults

`func NewCreateOrgInputBodyWithDefaults() *CreateOrgInputBody`

NewCreateOrgInputBodyWithDefaults instantiates a new CreateOrgInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *CreateOrgInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *CreateOrgInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *CreateOrgInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *CreateOrgInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetName

`func (o *CreateOrgInputBody) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateOrgInputBody) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateOrgInputBody) SetName(v string)`

SetName sets Name field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


