# AddMemberInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Email** | **string** |  | 
**Role** | Pointer to **string** |  | [optional] 

## Methods

### NewAddMemberInputBody

`func NewAddMemberInputBody(email string, ) *AddMemberInputBody`

NewAddMemberInputBody instantiates a new AddMemberInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAddMemberInputBodyWithDefaults

`func NewAddMemberInputBodyWithDefaults() *AddMemberInputBody`

NewAddMemberInputBodyWithDefaults instantiates a new AddMemberInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *AddMemberInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *AddMemberInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *AddMemberInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *AddMemberInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetEmail

`func (o *AddMemberInputBody) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *AddMemberInputBody) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *AddMemberInputBody) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetRole

`func (o *AddMemberInputBody) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *AddMemberInputBody) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *AddMemberInputBody) SetRole(v string)`

SetRole sets Role field to given value.

### HasRole

`func (o *AddMemberInputBody) HasRole() bool`

HasRole returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


