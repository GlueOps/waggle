# AddMemberOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Invited** | **bool** | An invite email was sent (pending account). | 
**Member** | [**MemberJSONView**](MemberJSONView.md) |  | 

## Methods

### NewAddMemberOutputBody

`func NewAddMemberOutputBody(invited bool, member MemberJSONView, ) *AddMemberOutputBody`

NewAddMemberOutputBody instantiates a new AddMemberOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAddMemberOutputBodyWithDefaults

`func NewAddMemberOutputBodyWithDefaults() *AddMemberOutputBody`

NewAddMemberOutputBodyWithDefaults instantiates a new AddMemberOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *AddMemberOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *AddMemberOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *AddMemberOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *AddMemberOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetInvited

`func (o *AddMemberOutputBody) GetInvited() bool`

GetInvited returns the Invited field if non-nil, zero value otherwise.

### GetInvitedOk

`func (o *AddMemberOutputBody) GetInvitedOk() (*bool, bool)`

GetInvitedOk returns a tuple with the Invited field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInvited

`func (o *AddMemberOutputBody) SetInvited(v bool)`

SetInvited sets Invited field to given value.


### GetMember

`func (o *AddMemberOutputBody) GetMember() MemberJSONView`

GetMember returns the Member field if non-nil, zero value otherwise.

### GetMemberOk

`func (o *AddMemberOutputBody) GetMemberOk() (*MemberJSONView, bool)`

GetMemberOk returns a tuple with the Member field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMember

`func (o *AddMemberOutputBody) SetMember(v MemberJSONView)`

SetMember sets Member field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


