# MemberJSONView

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**AccountId** | **string** |  | 
**CreatedAt** | **time.Time** |  | 
**DisplayName** | **string** |  | 
**Email** | **string** |  | 
**IsActive** | **bool** |  | 
**LastLoginAt** | Pointer to **time.Time** |  | [optional] 
**Pending** | **bool** | Invited but hasn&#39;t accepted (no password set yet). | 
**Role** | **string** |  | 
**UserId** | **string** |  | 

## Methods

### NewMemberJSONView

`func NewMemberJSONView(accountId string, createdAt time.Time, displayName string, email string, isActive bool, pending bool, role string, userId string, ) *MemberJSONView`

NewMemberJSONView instantiates a new MemberJSONView object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMemberJSONViewWithDefaults

`func NewMemberJSONViewWithDefaults() *MemberJSONView`

NewMemberJSONViewWithDefaults instantiates a new MemberJSONView object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *MemberJSONView) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *MemberJSONView) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *MemberJSONView) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *MemberJSONView) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetAccountId

`func (o *MemberJSONView) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *MemberJSONView) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *MemberJSONView) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetCreatedAt

`func (o *MemberJSONView) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *MemberJSONView) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *MemberJSONView) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetDisplayName

`func (o *MemberJSONView) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *MemberJSONView) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *MemberJSONView) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.


### GetEmail

`func (o *MemberJSONView) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *MemberJSONView) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *MemberJSONView) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetIsActive

`func (o *MemberJSONView) GetIsActive() bool`

GetIsActive returns the IsActive field if non-nil, zero value otherwise.

### GetIsActiveOk

`func (o *MemberJSONView) GetIsActiveOk() (*bool, bool)`

GetIsActiveOk returns a tuple with the IsActive field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsActive

`func (o *MemberJSONView) SetIsActive(v bool)`

SetIsActive sets IsActive field to given value.


### GetLastLoginAt

`func (o *MemberJSONView) GetLastLoginAt() time.Time`

GetLastLoginAt returns the LastLoginAt field if non-nil, zero value otherwise.

### GetLastLoginAtOk

`func (o *MemberJSONView) GetLastLoginAtOk() (*time.Time, bool)`

GetLastLoginAtOk returns a tuple with the LastLoginAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastLoginAt

`func (o *MemberJSONView) SetLastLoginAt(v time.Time)`

SetLastLoginAt sets LastLoginAt field to given value.

### HasLastLoginAt

`func (o *MemberJSONView) HasLastLoginAt() bool`

HasLastLoginAt returns a boolean if a field has been set.

### GetPending

`func (o *MemberJSONView) GetPending() bool`

GetPending returns the Pending field if non-nil, zero value otherwise.

### GetPendingOk

`func (o *MemberJSONView) GetPendingOk() (*bool, bool)`

GetPendingOk returns a tuple with the Pending field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPending

`func (o *MemberJSONView) SetPending(v bool)`

SetPending sets Pending field to given value.


### GetRole

`func (o *MemberJSONView) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *MemberJSONView) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *MemberJSONView) SetRole(v string)`

SetRole sets Role field to given value.


### GetUserId

`func (o *MemberJSONView) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *MemberJSONView) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *MemberJSONView) SetUserId(v string)`

SetUserId sets UserId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


