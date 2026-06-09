# OrgListOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Items** | [**[]OrgFullView**](OrgFullView.md) |  | 

## Methods

### NewOrgListOutputBody

`func NewOrgListOutputBody(items []OrgFullView, ) *OrgListOutputBody`

NewOrgListOutputBody instantiates a new OrgListOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrgListOutputBodyWithDefaults

`func NewOrgListOutputBodyWithDefaults() *OrgListOutputBody`

NewOrgListOutputBodyWithDefaults instantiates a new OrgListOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *OrgListOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *OrgListOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *OrgListOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *OrgListOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetItems

`func (o *OrgListOutputBody) GetItems() []OrgFullView`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *OrgListOutputBody) GetItemsOk() (*[]OrgFullView, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *OrgListOutputBody) SetItems(v []OrgFullView)`

SetItems sets Items field to given value.


### SetItemsNil

`func (o *OrgListOutputBody) SetItemsNil(b bool)`

 SetItemsNil sets the value for Items to be an explicit nil

### UnsetItems
`func (o *OrgListOutputBody) UnsetItems()`

UnsetItems ensures that no value is present for Items, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


