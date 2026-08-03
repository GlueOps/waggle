# DatacenterBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**CpuOvercommitRatio** | Pointer to **float64** |  | [optional] 
**InsecureSkipVerify** | Pointer to **bool** |  | [optional] 
**Name** | **string** |  | 
**Token** | Pointer to **string** |  | [optional] 
**Url** | **string** |  | 

## Methods

### NewDatacenterBody

`func NewDatacenterBody(name string, url string, ) *DatacenterBody`

NewDatacenterBody instantiates a new DatacenterBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDatacenterBodyWithDefaults

`func NewDatacenterBodyWithDefaults() *DatacenterBody`

NewDatacenterBodyWithDefaults instantiates a new DatacenterBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *DatacenterBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *DatacenterBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *DatacenterBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *DatacenterBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetCpuOvercommitRatio

`func (o *DatacenterBody) GetCpuOvercommitRatio() float64`

GetCpuOvercommitRatio returns the CpuOvercommitRatio field if non-nil, zero value otherwise.

### GetCpuOvercommitRatioOk

`func (o *DatacenterBody) GetCpuOvercommitRatioOk() (*float64, bool)`

GetCpuOvercommitRatioOk returns a tuple with the CpuOvercommitRatio field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpuOvercommitRatio

`func (o *DatacenterBody) SetCpuOvercommitRatio(v float64)`

SetCpuOvercommitRatio sets CpuOvercommitRatio field to given value.

### HasCpuOvercommitRatio

`func (o *DatacenterBody) HasCpuOvercommitRatio() bool`

HasCpuOvercommitRatio returns a boolean if a field has been set.

### GetInsecureSkipVerify

`func (o *DatacenterBody) GetInsecureSkipVerify() bool`

GetInsecureSkipVerify returns the InsecureSkipVerify field if non-nil, zero value otherwise.

### GetInsecureSkipVerifyOk

`func (o *DatacenterBody) GetInsecureSkipVerifyOk() (*bool, bool)`

GetInsecureSkipVerifyOk returns a tuple with the InsecureSkipVerify field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInsecureSkipVerify

`func (o *DatacenterBody) SetInsecureSkipVerify(v bool)`

SetInsecureSkipVerify sets InsecureSkipVerify field to given value.

### HasInsecureSkipVerify

`func (o *DatacenterBody) HasInsecureSkipVerify() bool`

HasInsecureSkipVerify returns a boolean if a field has been set.

### GetName

`func (o *DatacenterBody) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DatacenterBody) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DatacenterBody) SetName(v string)`

SetName sets Name field to given value.


### GetToken

`func (o *DatacenterBody) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *DatacenterBody) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *DatacenterBody) SetToken(v string)`

SetToken sets Token field to given value.

### HasToken

`func (o *DatacenterBody) HasToken() bool`

HasToken returns a boolean if a field has been set.

### GetUrl

`func (o *DatacenterBody) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *DatacenterBody) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *DatacenterBody) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


