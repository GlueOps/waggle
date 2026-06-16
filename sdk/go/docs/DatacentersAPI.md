# \DatacentersAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateDatacenter**](DatacentersAPI.md#CreateDatacenter) | **Post** /datacenters | Create a datacenter in the caller&#39;s tenant.
[**DeleteDatacenter**](DatacentersAPI.md#DeleteDatacenter) | **Delete** /datacenters/{id} | Delete a datacenter.
[**DiscoverHypervisors**](DatacentersAPI.md#DiscoverHypervisors) | **Post** /datacenters/{id}/discover | Discover hypervisors from the datacenter&#39;s Proxmox cluster and upsert them (preserving reserved capacity and schedulable). Set async to run in the background.
[**GetDatacenter**](DatacentersAPI.md#GetDatacenter) | **Get** /datacenters/{id} | Fetch a datacenter by ID.
[**ListDatacenters**](DatacentersAPI.md#ListDatacenters) | **Get** /datacenters | List datacenters in the caller&#39;s tenant.
[**UpdateDatacenter**](DatacentersAPI.md#UpdateDatacenter) | **Put** /datacenters/{id} | Update a datacenter.



## CreateDatacenter

> DatacenterView CreateDatacenter(ctx).DatacenterBody(datacenterBody).Execute()

Create a datacenter in the caller's tenant.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	datacenterBody := *openapiclient.NewDatacenterBody("Name_example", "Url_example") // DatacenterBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DatacentersAPI.CreateDatacenter(context.Background()).DatacenterBody(datacenterBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DatacentersAPI.CreateDatacenter``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateDatacenter`: DatacenterView
	fmt.Fprintf(os.Stdout, "Response from `DatacentersAPI.CreateDatacenter`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateDatacenterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **datacenterBody** | [**DatacenterBody**](DatacenterBody.md) |  | 

### Return type

[**DatacenterView**](DatacenterView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteDatacenter

> DeleteDatacenter(ctx, id).Execute()

Delete a datacenter.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	id := "id_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DatacentersAPI.DeleteDatacenter(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DatacentersAPI.DeleteDatacenter``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDatacenterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DiscoverHypervisors

> DiscoverOutputBody DiscoverHypervisors(ctx, id).DiscoverInputBody(discoverInputBody).Execute()

Discover hypervisors from the datacenter's Proxmox cluster and upsert them (preserving reserved capacity and schedulable). Set async to run in the background.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	id := "id_example" // string | 
	discoverInputBody := *openapiclient.NewDiscoverInputBody() // DiscoverInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DatacentersAPI.DiscoverHypervisors(context.Background(), id).DiscoverInputBody(discoverInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DatacentersAPI.DiscoverHypervisors``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DiscoverHypervisors`: DiscoverOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DatacentersAPI.DiscoverHypervisors`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDiscoverHypervisorsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **discoverInputBody** | [**DiscoverInputBody**](DiscoverInputBody.md) |  | 

### Return type

[**DiscoverOutputBody**](DiscoverOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDatacenter

> DatacenterView GetDatacenter(ctx, id).Execute()

Fetch a datacenter by ID.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	id := "id_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DatacentersAPI.GetDatacenter(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DatacentersAPI.GetDatacenter``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDatacenter`: DatacenterView
	fmt.Fprintf(os.Stdout, "Response from `DatacentersAPI.GetDatacenter`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDatacenterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DatacenterView**](DatacenterView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListDatacenters

> DatacenterListOutputBody ListDatacenters(ctx).Execute()

List datacenters in the caller's tenant.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DatacentersAPI.ListDatacenters(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DatacentersAPI.ListDatacenters``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListDatacenters`: DatacenterListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DatacentersAPI.ListDatacenters`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListDatacentersRequest struct via the builder pattern


### Return type

[**DatacenterListOutputBody**](DatacenterListOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateDatacenter

> DatacenterView UpdateDatacenter(ctx, id).DatacenterBody(datacenterBody).Execute()

Update a datacenter.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	id := "id_example" // string | 
	datacenterBody := *openapiclient.NewDatacenterBody("Name_example", "Url_example") // DatacenterBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DatacentersAPI.UpdateDatacenter(context.Background(), id).DatacenterBody(datacenterBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DatacentersAPI.UpdateDatacenter``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateDatacenter`: DatacenterView
	fmt.Fprintf(os.Stdout, "Response from `DatacentersAPI.UpdateDatacenter`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateDatacenterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **datacenterBody** | [**DatacenterBody**](DatacenterBody.md) |  | 

### Return type

[**DatacenterView**](DatacenterView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

