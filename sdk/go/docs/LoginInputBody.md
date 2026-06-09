# LoginInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Email** | **string** |  | 
**OrganizationId** | Pointer to **string** |  | [optional] 
**Password** | **string** |  | 

## Methods

### NewLoginInputBody

`func NewLoginInputBody(email string, password string, ) *LoginInputBody`

NewLoginInputBody instantiates a new LoginInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLoginInputBodyWithDefaults

`func NewLoginInputBodyWithDefaults() *LoginInputBody`

NewLoginInputBodyWithDefaults instantiates a new LoginInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *LoginInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *LoginInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *LoginInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *LoginInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetEmail

`func (o *LoginInputBody) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *LoginInputBody) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *LoginInputBody) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetOrganizationId

`func (o *LoginInputBody) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *LoginInputBody) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *LoginInputBody) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.

### HasOrganizationId

`func (o *LoginInputBody) HasOrganizationId() bool`

HasOrganizationId returns a boolean if a field has been set.

### GetPassword

`func (o *LoginInputBody) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *LoginInputBody) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *LoginInputBody) SetPassword(v string)`

SetPassword sets Password field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


