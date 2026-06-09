# SwitchOrgInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**OrganizationId** | **string** |  | 

## Methods

### NewSwitchOrgInputBody

`func NewSwitchOrgInputBody(organizationId string, ) *SwitchOrgInputBody`

NewSwitchOrgInputBody instantiates a new SwitchOrgInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwitchOrgInputBodyWithDefaults

`func NewSwitchOrgInputBodyWithDefaults() *SwitchOrgInputBody`

NewSwitchOrgInputBodyWithDefaults instantiates a new SwitchOrgInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *SwitchOrgInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *SwitchOrgInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *SwitchOrgInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *SwitchOrgInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetOrganizationId

`func (o *SwitchOrgInputBody) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *SwitchOrgInputBody) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *SwitchOrgInputBody) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


