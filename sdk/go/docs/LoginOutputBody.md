# LoginOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**AccountId** | Pointer to **string** |  | [optional] 
**Memberships** | Pointer to [**[]MembershipView**](MembershipView.md) |  | [optional] 
**Organization** | Pointer to [**OrgView**](OrgView.md) |  | [optional] 
**Tokens** | Pointer to [**AuthTokens**](AuthTokens.md) |  | [optional] 

## Methods

### NewLoginOutputBody

`func NewLoginOutputBody() *LoginOutputBody`

NewLoginOutputBody instantiates a new LoginOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLoginOutputBodyWithDefaults

`func NewLoginOutputBodyWithDefaults() *LoginOutputBody`

NewLoginOutputBodyWithDefaults instantiates a new LoginOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *LoginOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *LoginOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *LoginOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *LoginOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetAccountId

`func (o *LoginOutputBody) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *LoginOutputBody) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *LoginOutputBody) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.

### HasAccountId

`func (o *LoginOutputBody) HasAccountId() bool`

HasAccountId returns a boolean if a field has been set.

### GetMemberships

`func (o *LoginOutputBody) GetMemberships() []MembershipView`

GetMemberships returns the Memberships field if non-nil, zero value otherwise.

### GetMembershipsOk

`func (o *LoginOutputBody) GetMembershipsOk() (*[]MembershipView, bool)`

GetMembershipsOk returns a tuple with the Memberships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemberships

`func (o *LoginOutputBody) SetMemberships(v []MembershipView)`

SetMemberships sets Memberships field to given value.

### HasMemberships

`func (o *LoginOutputBody) HasMemberships() bool`

HasMemberships returns a boolean if a field has been set.

### SetMembershipsNil

`func (o *LoginOutputBody) SetMembershipsNil(b bool)`

 SetMembershipsNil sets the value for Memberships to be an explicit nil

### UnsetMemberships
`func (o *LoginOutputBody) UnsetMemberships()`

UnsetMemberships ensures that no value is present for Memberships, not even an explicit nil
### GetOrganization

`func (o *LoginOutputBody) GetOrganization() OrgView`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *LoginOutputBody) GetOrganizationOk() (*OrgView, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *LoginOutputBody) SetOrganization(v OrgView)`

SetOrganization sets Organization field to given value.

### HasOrganization

`func (o *LoginOutputBody) HasOrganization() bool`

HasOrganization returns a boolean if a field has been set.

### GetTokens

`func (o *LoginOutputBody) GetTokens() AuthTokens`

GetTokens returns the Tokens field if non-nil, zero value otherwise.

### GetTokensOk

`func (o *LoginOutputBody) GetTokensOk() (*AuthTokens, bool)`

GetTokensOk returns a tuple with the Tokens field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokens

`func (o *LoginOutputBody) SetTokens(v AuthTokens)`

SetTokens sets Tokens field to given value.

### HasTokens

`func (o *LoginOutputBody) HasTokens() bool`

HasTokens returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


