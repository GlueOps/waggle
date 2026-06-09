# SignupOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**AccessExpiresAt** | **time.Time** |  | 
**AccessToken** | **string** |  | 
**AccountId** | **string** |  | 
**Organization** | [**OrgView**](OrgView.md) |  | 
**RefreshExpiresAt** | **time.Time** |  | 
**RefreshToken** | **string** |  | 

## Methods

### NewSignupOutputBody

`func NewSignupOutputBody(accessExpiresAt time.Time, accessToken string, accountId string, organization OrgView, refreshExpiresAt time.Time, refreshToken string, ) *SignupOutputBody`

NewSignupOutputBody instantiates a new SignupOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSignupOutputBodyWithDefaults

`func NewSignupOutputBodyWithDefaults() *SignupOutputBody`

NewSignupOutputBodyWithDefaults instantiates a new SignupOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *SignupOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *SignupOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *SignupOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *SignupOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetAccessExpiresAt

`func (o *SignupOutputBody) GetAccessExpiresAt() time.Time`

GetAccessExpiresAt returns the AccessExpiresAt field if non-nil, zero value otherwise.

### GetAccessExpiresAtOk

`func (o *SignupOutputBody) GetAccessExpiresAtOk() (*time.Time, bool)`

GetAccessExpiresAtOk returns a tuple with the AccessExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessExpiresAt

`func (o *SignupOutputBody) SetAccessExpiresAt(v time.Time)`

SetAccessExpiresAt sets AccessExpiresAt field to given value.


### GetAccessToken

`func (o *SignupOutputBody) GetAccessToken() string`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *SignupOutputBody) GetAccessTokenOk() (*string, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *SignupOutputBody) SetAccessToken(v string)`

SetAccessToken sets AccessToken field to given value.


### GetAccountId

`func (o *SignupOutputBody) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *SignupOutputBody) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *SignupOutputBody) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetOrganization

`func (o *SignupOutputBody) GetOrganization() OrgView`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *SignupOutputBody) GetOrganizationOk() (*OrgView, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *SignupOutputBody) SetOrganization(v OrgView)`

SetOrganization sets Organization field to given value.


### GetRefreshExpiresAt

`func (o *SignupOutputBody) GetRefreshExpiresAt() time.Time`

GetRefreshExpiresAt returns the RefreshExpiresAt field if non-nil, zero value otherwise.

### GetRefreshExpiresAtOk

`func (o *SignupOutputBody) GetRefreshExpiresAtOk() (*time.Time, bool)`

GetRefreshExpiresAtOk returns a tuple with the RefreshExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshExpiresAt

`func (o *SignupOutputBody) SetRefreshExpiresAt(v time.Time)`

SetRefreshExpiresAt sets RefreshExpiresAt field to given value.


### GetRefreshToken

`func (o *SignupOutputBody) GetRefreshToken() string`

GetRefreshToken returns the RefreshToken field if non-nil, zero value otherwise.

### GetRefreshTokenOk

`func (o *SignupOutputBody) GetRefreshTokenOk() (*string, bool)`

GetRefreshTokenOk returns a tuple with the RefreshToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshToken

`func (o *SignupOutputBody) SetRefreshToken(v string)`

SetRefreshToken sets RefreshToken field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


