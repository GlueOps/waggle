# MeOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**AccountId** | **string** |  | 
**CurrentOrganization** | Pointer to [**OrgView**](OrgView.md) |  | [optional] 
**DisplayName** | **string** |  | 
**Emails** | [**[]AccountEmailView**](AccountEmailView.md) |  | 
**LastLoginAt** | Pointer to **time.Time** |  | [optional] 
**Memberships** | [**[]MembershipView**](MembershipView.md) |  | 

## Methods

### NewMeOutputBody

`func NewMeOutputBody(accountId string, displayName string, emails []AccountEmailView, memberships []MembershipView, ) *MeOutputBody`

NewMeOutputBody instantiates a new MeOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMeOutputBodyWithDefaults

`func NewMeOutputBodyWithDefaults() *MeOutputBody`

NewMeOutputBodyWithDefaults instantiates a new MeOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *MeOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *MeOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *MeOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *MeOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetAccountId

`func (o *MeOutputBody) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *MeOutputBody) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *MeOutputBody) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetCurrentOrganization

`func (o *MeOutputBody) GetCurrentOrganization() OrgView`

GetCurrentOrganization returns the CurrentOrganization field if non-nil, zero value otherwise.

### GetCurrentOrganizationOk

`func (o *MeOutputBody) GetCurrentOrganizationOk() (*OrgView, bool)`

GetCurrentOrganizationOk returns a tuple with the CurrentOrganization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentOrganization

`func (o *MeOutputBody) SetCurrentOrganization(v OrgView)`

SetCurrentOrganization sets CurrentOrganization field to given value.

### HasCurrentOrganization

`func (o *MeOutputBody) HasCurrentOrganization() bool`

HasCurrentOrganization returns a boolean if a field has been set.

### GetDisplayName

`func (o *MeOutputBody) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *MeOutputBody) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *MeOutputBody) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.


### GetEmails

`func (o *MeOutputBody) GetEmails() []AccountEmailView`

GetEmails returns the Emails field if non-nil, zero value otherwise.

### GetEmailsOk

`func (o *MeOutputBody) GetEmailsOk() (*[]AccountEmailView, bool)`

GetEmailsOk returns a tuple with the Emails field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmails

`func (o *MeOutputBody) SetEmails(v []AccountEmailView)`

SetEmails sets Emails field to given value.


### SetEmailsNil

`func (o *MeOutputBody) SetEmailsNil(b bool)`

 SetEmailsNil sets the value for Emails to be an explicit nil

### UnsetEmails
`func (o *MeOutputBody) UnsetEmails()`

UnsetEmails ensures that no value is present for Emails, not even an explicit nil
### GetLastLoginAt

`func (o *MeOutputBody) GetLastLoginAt() time.Time`

GetLastLoginAt returns the LastLoginAt field if non-nil, zero value otherwise.

### GetLastLoginAtOk

`func (o *MeOutputBody) GetLastLoginAtOk() (*time.Time, bool)`

GetLastLoginAtOk returns a tuple with the LastLoginAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastLoginAt

`func (o *MeOutputBody) SetLastLoginAt(v time.Time)`

SetLastLoginAt sets LastLoginAt field to given value.

### HasLastLoginAt

`func (o *MeOutputBody) HasLastLoginAt() bool`

HasLastLoginAt returns a boolean if a field has been set.

### GetMemberships

`func (o *MeOutputBody) GetMemberships() []MembershipView`

GetMemberships returns the Memberships field if non-nil, zero value otherwise.

### GetMembershipsOk

`func (o *MeOutputBody) GetMembershipsOk() (*[]MembershipView, bool)`

GetMembershipsOk returns a tuple with the Memberships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemberships

`func (o *MeOutputBody) SetMemberships(v []MembershipView)`

SetMemberships sets Memberships field to given value.


### SetMembershipsNil

`func (o *MeOutputBody) SetMembershipsNil(b bool)`

 SetMembershipsNil sets the value for Memberships to be an explicit nil

### UnsetMemberships
`func (o *MeOutputBody) UnsetMemberships()`

UnsetMemberships ensures that no value is present for Memberships, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


