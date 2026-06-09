# AcceptInviteInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**DisplayName** | Pointer to **string** |  | [optional] 
**Password** | Pointer to **string** |  | [optional] 
**Token** | **string** |  | 

## Methods

### NewAcceptInviteInputBody

`func NewAcceptInviteInputBody(token string, ) *AcceptInviteInputBody`

NewAcceptInviteInputBody instantiates a new AcceptInviteInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAcceptInviteInputBodyWithDefaults

`func NewAcceptInviteInputBodyWithDefaults() *AcceptInviteInputBody`

NewAcceptInviteInputBodyWithDefaults instantiates a new AcceptInviteInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *AcceptInviteInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *AcceptInviteInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *AcceptInviteInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *AcceptInviteInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetDisplayName

`func (o *AcceptInviteInputBody) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *AcceptInviteInputBody) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *AcceptInviteInputBody) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.

### HasDisplayName

`func (o *AcceptInviteInputBody) HasDisplayName() bool`

HasDisplayName returns a boolean if a field has been set.

### GetPassword

`func (o *AcceptInviteInputBody) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *AcceptInviteInputBody) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *AcceptInviteInputBody) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *AcceptInviteInputBody) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetToken

`func (o *AcceptInviteInputBody) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *AcceptInviteInputBody) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *AcceptInviteInputBody) SetToken(v string)`

SetToken sets Token field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


