# \HypervisorsAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateHypervisor**](HypervisorsAPI.md#CreateHypervisor) | **Post** /hypervisors | Create a hypervisor in the caller&#39;s tenant.
[**DeleteHypervisor**](HypervisorsAPI.md#DeleteHypervisor) | **Delete** /hypervisors/{id} | Delete a hypervisor.
[**GetHypervisor**](HypervisorsAPI.md#GetHypervisor) | **Get** /hypervisors/{id} | Fetch a hypervisor by ID.
[**ListHypervisors**](HypervisorsAPI.md#ListHypervisors) | **Get** /hypervisors | List hypervisors in the caller&#39;s tenant.
[**UpdateHypervisor**](HypervisorsAPI.md#UpdateHypervisor) | **Put** /hypervisors/{id} | Update a hypervisor.



## CreateHypervisor

> HypervisorView CreateHypervisor(ctx).HypervisorBody(hypervisorBody).Execute()

Create a hypervisor in the caller's tenant.

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
	hypervisorBody := *openapiclient.NewHypervisorBody(int64(123), int64(123), "DatacenterId_example", int64(123), int64(123), "Name_example", int64(123), int64(123)) // HypervisorBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.HypervisorsAPI.CreateHypervisor(context.Background()).HypervisorBody(hypervisorBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HypervisorsAPI.CreateHypervisor``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateHypervisor`: HypervisorView
	fmt.Fprintf(os.Stdout, "Response from `HypervisorsAPI.CreateHypervisor`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateHypervisorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **hypervisorBody** | [**HypervisorBody**](HypervisorBody.md) |  | 

### Return type

[**HypervisorView**](HypervisorView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteHypervisor

> DeleteHypervisor(ctx, id).Execute()

Delete a hypervisor.

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
	r, err := apiClient.HypervisorsAPI.DeleteHypervisor(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HypervisorsAPI.DeleteHypervisor``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDeleteHypervisorRequest struct via the builder pattern


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


## GetHypervisor

> HypervisorView GetHypervisor(ctx, id).Execute()

Fetch a hypervisor by ID.

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
	resp, r, err := apiClient.HypervisorsAPI.GetHypervisor(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HypervisorsAPI.GetHypervisor``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetHypervisor`: HypervisorView
	fmt.Fprintf(os.Stdout, "Response from `HypervisorsAPI.GetHypervisor`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetHypervisorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**HypervisorView**](HypervisorView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListHypervisors

> HypervisorListOutputBody ListHypervisors(ctx).DatacenterId(datacenterId).Execute()

List hypervisors in the caller's tenant.

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
	datacenterId := "datacenterId_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.HypervisorsAPI.ListHypervisors(context.Background()).DatacenterId(datacenterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HypervisorsAPI.ListHypervisors``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListHypervisors`: HypervisorListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `HypervisorsAPI.ListHypervisors`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListHypervisorsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **datacenterId** | **string** |  | 

### Return type

[**HypervisorListOutputBody**](HypervisorListOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateHypervisor

> HypervisorView UpdateHypervisor(ctx, id).HypervisorBody(hypervisorBody).Execute()

Update a hypervisor.

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
	hypervisorBody := *openapiclient.NewHypervisorBody(int64(123), int64(123), "DatacenterId_example", int64(123), int64(123), "Name_example", int64(123), int64(123)) // HypervisorBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.HypervisorsAPI.UpdateHypervisor(context.Background(), id).HypervisorBody(hypervisorBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HypervisorsAPI.UpdateHypervisor``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateHypervisor`: HypervisorView
	fmt.Fprintf(os.Stdout, "Response from `HypervisorsAPI.UpdateHypervisor`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateHypervisorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **hypervisorBody** | [**HypervisorBody**](HypervisorBody.md) |  | 

### Return type

[**HypervisorView**](HypervisorView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

