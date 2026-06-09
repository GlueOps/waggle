# DatacenterView

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**CreatedAt** | **time.Time** |  | 
**HasToken** | **bool** | Whether a Proxmox API token is configured (the token itself is never returned). | 
**Id** | **string** |  | 
**InsecureSkipVerify** | **bool** | Whether TLS verification is disabled for this cluster (self-signed certs). | 
**Name** | **string** |  | 
**UpdatedAt** | **time.Time** |  | 
**Url** | **string** |  | 

## Methods

### NewDatacenterView

`func NewDatacenterView(createdAt time.Time, hasToken bool, id string, insecureSkipVerify bool, name string, updatedAt time.Time, url string, ) *DatacenterView`

NewDatacenterView instantiates a new DatacenterView object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDatacenterViewWithDefaults

`func NewDatacenterViewWithDefaults() *DatacenterView`

NewDatacenterViewWithDefaults instantiates a new DatacenterView object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *DatacenterView) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *DatacenterView) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *DatacenterView) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *DatacenterView) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetCreatedAt

`func (o *DatacenterView) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *DatacenterView) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *DatacenterView) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetHasToken

`func (o *DatacenterView) GetHasToken() bool`

GetHasToken returns the HasToken field if non-nil, zero value otherwise.

### GetHasTokenOk

`func (o *DatacenterView) GetHasTokenOk() (*bool, bool)`

GetHasTokenOk returns a tuple with the HasToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasToken

`func (o *DatacenterView) SetHasToken(v bool)`

SetHasToken sets HasToken field to given value.


### GetId

`func (o *DatacenterView) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DatacenterView) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DatacenterView) SetId(v string)`

SetId sets Id field to given value.


### GetInsecureSkipVerify

`func (o *DatacenterView) GetInsecureSkipVerify() bool`

GetInsecureSkipVerify returns the InsecureSkipVerify field if non-nil, zero value otherwise.

### GetInsecureSkipVerifyOk

`func (o *DatacenterView) GetInsecureSkipVerifyOk() (*bool, bool)`

GetInsecureSkipVerifyOk returns a tuple with the InsecureSkipVerify field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInsecureSkipVerify

`func (o *DatacenterView) SetInsecureSkipVerify(v bool)`

SetInsecureSkipVerify sets InsecureSkipVerify field to given value.


### GetName

`func (o *DatacenterView) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DatacenterView) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DatacenterView) SetName(v string)`

SetName sets Name field to given value.


### GetUpdatedAt

`func (o *DatacenterView) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *DatacenterView) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *DatacenterView) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetUrl

`func (o *DatacenterView) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *DatacenterView) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *DatacenterView) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


