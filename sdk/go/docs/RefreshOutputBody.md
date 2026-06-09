# RefreshOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**AccessExpiresAt** | **time.Time** |  | 
**AccessToken** | **string** |  | 
**RefreshExpiresAt** | **time.Time** |  | 
**RefreshToken** | **string** |  | 

## Methods

### NewRefreshOutputBody

`func NewRefreshOutputBody(accessExpiresAt time.Time, accessToken string, refreshExpiresAt time.Time, refreshToken string, ) *RefreshOutputBody`

NewRefreshOutputBody instantiates a new RefreshOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRefreshOutputBodyWithDefaults

`func NewRefreshOutputBodyWithDefaults() *RefreshOutputBody`

NewRefreshOutputBodyWithDefaults instantiates a new RefreshOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *RefreshOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *RefreshOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *RefreshOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *RefreshOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetAccessExpiresAt

`func (o *RefreshOutputBody) GetAccessExpiresAt() time.Time`

GetAccessExpiresAt returns the AccessExpiresAt field if non-nil, zero value otherwise.

### GetAccessExpiresAtOk

`func (o *RefreshOutputBody) GetAccessExpiresAtOk() (*time.Time, bool)`

GetAccessExpiresAtOk returns a tuple with the AccessExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessExpiresAt

`func (o *RefreshOutputBody) SetAccessExpiresAt(v time.Time)`

SetAccessExpiresAt sets AccessExpiresAt field to given value.


### GetAccessToken

`func (o *RefreshOutputBody) GetAccessToken() string`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *RefreshOutputBody) GetAccessTokenOk() (*string, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *RefreshOutputBody) SetAccessToken(v string)`

SetAccessToken sets AccessToken field to given value.


### GetRefreshExpiresAt

`func (o *RefreshOutputBody) GetRefreshExpiresAt() time.Time`

GetRefreshExpiresAt returns the RefreshExpiresAt field if non-nil, zero value otherwise.

### GetRefreshExpiresAtOk

`func (o *RefreshOutputBody) GetRefreshExpiresAtOk() (*time.Time, bool)`

GetRefreshExpiresAtOk returns a tuple with the RefreshExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshExpiresAt

`func (o *RefreshOutputBody) SetRefreshExpiresAt(v time.Time)`

SetRefreshExpiresAt sets RefreshExpiresAt field to given value.


### GetRefreshToken

`func (o *RefreshOutputBody) GetRefreshToken() string`

GetRefreshToken returns the RefreshToken field if non-nil, zero value otherwise.

### GetRefreshTokenOk

`func (o *RefreshOutputBody) GetRefreshTokenOk() (*string, bool)`

GetRefreshTokenOk returns a tuple with the RefreshToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshToken

`func (o *RefreshOutputBody) SetRefreshToken(v string)`

SetRefreshToken sets RefreshToken field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


