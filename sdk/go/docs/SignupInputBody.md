# SignupInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**DisplayName** | Pointer to **string** |  | [optional] 
**Email** | **string** |  | 
**OrganizationName** | **string** |  | 
**Password** | **string** |  | 

## Methods

### NewSignupInputBody

`func NewSignupInputBody(email string, organizationName string, password string, ) *SignupInputBody`

NewSignupInputBody instantiates a new SignupInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSignupInputBodyWithDefaults

`func NewSignupInputBodyWithDefaults() *SignupInputBody`

NewSignupInputBodyWithDefaults instantiates a new SignupInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *SignupInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *SignupInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *SignupInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *SignupInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetDisplayName

`func (o *SignupInputBody) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *SignupInputBody) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *SignupInputBody) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.

### HasDisplayName

`func (o *SignupInputBody) HasDisplayName() bool`

HasDisplayName returns a boolean if a field has been set.

### GetEmail

`func (o *SignupInputBody) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *SignupInputBody) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *SignupInputBody) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetOrganizationName

`func (o *SignupInputBody) GetOrganizationName() string`

GetOrganizationName returns the OrganizationName field if non-nil, zero value otherwise.

### GetOrganizationNameOk

`func (o *SignupInputBody) GetOrganizationNameOk() (*string, bool)`

GetOrganizationNameOk returns a tuple with the OrganizationName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationName

`func (o *SignupInputBody) SetOrganizationName(v string)`

SetOrganizationName sets OrganizationName field to given value.


### GetPassword

`func (o *SignupInputBody) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *SignupInputBody) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *SignupInputBody) SetPassword(v string)`

SetPassword sets Password field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


