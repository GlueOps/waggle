# OrgFullView

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**CreatedAt** | **time.Time** |  | 
**Domain** | Pointer to **string** |  | [optional] 
**Id** | **string** |  | 
**Name** | **string** |  | 
**Role** | **string** | The calling account&#39;s role in this organization. | 
**Slug** | **string** |  | 
**Status** | **string** |  | 

## Methods

### NewOrgFullView

`func NewOrgFullView(createdAt time.Time, id string, name string, role string, slug string, status string, ) *OrgFullView`

NewOrgFullView instantiates a new OrgFullView object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrgFullViewWithDefaults

`func NewOrgFullViewWithDefaults() *OrgFullView`

NewOrgFullViewWithDefaults instantiates a new OrgFullView object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *OrgFullView) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *OrgFullView) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *OrgFullView) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *OrgFullView) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetCreatedAt

`func (o *OrgFullView) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *OrgFullView) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *OrgFullView) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetDomain

`func (o *OrgFullView) GetDomain() string`

GetDomain returns the Domain field if non-nil, zero value otherwise.

### GetDomainOk

`func (o *OrgFullView) GetDomainOk() (*string, bool)`

GetDomainOk returns a tuple with the Domain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomain

`func (o *OrgFullView) SetDomain(v string)`

SetDomain sets Domain field to given value.

### HasDomain

`func (o *OrgFullView) HasDomain() bool`

HasDomain returns a boolean if a field has been set.

### GetId

`func (o *OrgFullView) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *OrgFullView) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *OrgFullView) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *OrgFullView) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *OrgFullView) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *OrgFullView) SetName(v string)`

SetName sets Name field to given value.


### GetRole

`func (o *OrgFullView) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *OrgFullView) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *OrgFullView) SetRole(v string)`

SetRole sets Role field to given value.


### GetSlug

`func (o *OrgFullView) GetSlug() string`

GetSlug returns the Slug field if non-nil, zero value otherwise.

### GetSlugOk

`func (o *OrgFullView) GetSlugOk() (*string, bool)`

GetSlugOk returns a tuple with the Slug field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSlug

`func (o *OrgFullView) SetSlug(v string)`

SetSlug sets Slug field to given value.


### GetStatus

`func (o *OrgFullView) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *OrgFullView) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *OrgFullView) SetStatus(v string)`

SetStatus sets Status field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


