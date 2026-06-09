# AccountEmailView

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** |  | 
**Id** | **string** |  | 
**IsPrimary** | **bool** |  | 
**VerifiedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewAccountEmailView

`func NewAccountEmailView(email string, id string, isPrimary bool, ) *AccountEmailView`

NewAccountEmailView instantiates a new AccountEmailView object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAccountEmailViewWithDefaults

`func NewAccountEmailViewWithDefaults() *AccountEmailView`

NewAccountEmailViewWithDefaults instantiates a new AccountEmailView object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *AccountEmailView) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *AccountEmailView) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *AccountEmailView) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetId

`func (o *AccountEmailView) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AccountEmailView) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AccountEmailView) SetId(v string)`

SetId sets Id field to given value.


### GetIsPrimary

`func (o *AccountEmailView) GetIsPrimary() bool`

GetIsPrimary returns the IsPrimary field if non-nil, zero value otherwise.

### GetIsPrimaryOk

`func (o *AccountEmailView) GetIsPrimaryOk() (*bool, bool)`

GetIsPrimaryOk returns a tuple with the IsPrimary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPrimary

`func (o *AccountEmailView) SetIsPrimary(v bool)`

SetIsPrimary sets IsPrimary field to given value.


### GetVerifiedAt

`func (o *AccountEmailView) GetVerifiedAt() time.Time`

GetVerifiedAt returns the VerifiedAt field if non-nil, zero value otherwise.

### GetVerifiedAtOk

`func (o *AccountEmailView) GetVerifiedAtOk() (*time.Time, bool)`

GetVerifiedAtOk returns a tuple with the VerifiedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVerifiedAt

`func (o *AccountEmailView) SetVerifiedAt(v time.Time)`

SetVerifiedAt sets VerifiedAt field to given value.

### HasVerifiedAt

`func (o *AccountEmailView) HasVerifiedAt() bool`

HasVerifiedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


