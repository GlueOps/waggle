# AuthTokens

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessExpiresAt** | **time.Time** |  | 
**AccessToken** | **string** |  | 
**RefreshExpiresAt** | **time.Time** |  | 
**RefreshToken** | **string** |  | 

## Methods

### NewAuthTokens

`func NewAuthTokens(accessExpiresAt time.Time, accessToken string, refreshExpiresAt time.Time, refreshToken string, ) *AuthTokens`

NewAuthTokens instantiates a new AuthTokens object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthTokensWithDefaults

`func NewAuthTokensWithDefaults() *AuthTokens`

NewAuthTokensWithDefaults instantiates a new AuthTokens object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccessExpiresAt

`func (o *AuthTokens) GetAccessExpiresAt() time.Time`

GetAccessExpiresAt returns the AccessExpiresAt field if non-nil, zero value otherwise.

### GetAccessExpiresAtOk

`func (o *AuthTokens) GetAccessExpiresAtOk() (*time.Time, bool)`

GetAccessExpiresAtOk returns a tuple with the AccessExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessExpiresAt

`func (o *AuthTokens) SetAccessExpiresAt(v time.Time)`

SetAccessExpiresAt sets AccessExpiresAt field to given value.


### GetAccessToken

`func (o *AuthTokens) GetAccessToken() string`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *AuthTokens) GetAccessTokenOk() (*string, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *AuthTokens) SetAccessToken(v string)`

SetAccessToken sets AccessToken field to given value.


### GetRefreshExpiresAt

`func (o *AuthTokens) GetRefreshExpiresAt() time.Time`

GetRefreshExpiresAt returns the RefreshExpiresAt field if non-nil, zero value otherwise.

### GetRefreshExpiresAtOk

`func (o *AuthTokens) GetRefreshExpiresAtOk() (*time.Time, bool)`

GetRefreshExpiresAtOk returns a tuple with the RefreshExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshExpiresAt

`func (o *AuthTokens) SetRefreshExpiresAt(v time.Time)`

SetRefreshExpiresAt sets RefreshExpiresAt field to given value.


### GetRefreshToken

`func (o *AuthTokens) GetRefreshToken() string`

GetRefreshToken returns the RefreshToken field if non-nil, zero value otherwise.

### GetRefreshTokenOk

`func (o *AuthTokens) GetRefreshTokenOk() (*string, bool)`

GetRefreshTokenOk returns a tuple with the RefreshToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshToken

`func (o *AuthTokens) SetRefreshToken(v string)`

SetRefreshToken sets RefreshToken field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


