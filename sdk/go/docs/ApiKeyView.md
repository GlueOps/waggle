# ApiKeyView

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreatedAt** | **time.Time** |  | 
**ExpiresAt** | Pointer to **time.Time** |  | [optional] 
**Id** | **string** |  | 
**LastUsedAt** | Pointer to **time.Time** |  | [optional] 
**Name** | **string** |  | 
**Prefix** | **string** | Leading characters of the key, for identification. | 
**RevokedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewApiKeyView

`func NewApiKeyView(createdAt time.Time, id string, name string, prefix string, ) *ApiKeyView`

NewApiKeyView instantiates a new ApiKeyView object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiKeyViewWithDefaults

`func NewApiKeyViewWithDefaults() *ApiKeyView`

NewApiKeyViewWithDefaults instantiates a new ApiKeyView object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedAt

`func (o *ApiKeyView) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ApiKeyView) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ApiKeyView) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetExpiresAt

`func (o *ApiKeyView) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *ApiKeyView) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *ApiKeyView) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *ApiKeyView) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### GetId

`func (o *ApiKeyView) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ApiKeyView) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ApiKeyView) SetId(v string)`

SetId sets Id field to given value.


### GetLastUsedAt

`func (o *ApiKeyView) GetLastUsedAt() time.Time`

GetLastUsedAt returns the LastUsedAt field if non-nil, zero value otherwise.

### GetLastUsedAtOk

`func (o *ApiKeyView) GetLastUsedAtOk() (*time.Time, bool)`

GetLastUsedAtOk returns a tuple with the LastUsedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastUsedAt

`func (o *ApiKeyView) SetLastUsedAt(v time.Time)`

SetLastUsedAt sets LastUsedAt field to given value.

### HasLastUsedAt

`func (o *ApiKeyView) HasLastUsedAt() bool`

HasLastUsedAt returns a boolean if a field has been set.

### GetName

`func (o *ApiKeyView) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ApiKeyView) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ApiKeyView) SetName(v string)`

SetName sets Name field to given value.


### GetPrefix

`func (o *ApiKeyView) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *ApiKeyView) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *ApiKeyView) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.


### GetRevokedAt

`func (o *ApiKeyView) GetRevokedAt() time.Time`

GetRevokedAt returns the RevokedAt field if non-nil, zero value otherwise.

### GetRevokedAtOk

`func (o *ApiKeyView) GetRevokedAtOk() (*time.Time, bool)`

GetRevokedAtOk returns a tuple with the RevokedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevokedAt

`func (o *ApiKeyView) SetRevokedAt(v time.Time)`

SetRevokedAt sets RevokedAt field to given value.

### HasRevokedAt

`func (o *ApiKeyView) HasRevokedAt() bool`

HasRevokedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


