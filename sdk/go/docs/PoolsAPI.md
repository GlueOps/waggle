# \PoolsAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreatePool**](PoolsAPI.md#CreatePool) | **Post** /pools | Create a node pool and place its VMs across hypervisors (anti-affinity spread, all-or-nothing). Placements are available at GET /pools/{id}/placements.
[**DeletePool**](PoolsAPI.md#DeletePool) | **Delete** /pools/{id} | Delete a pool and release all its placements.
[**GetPool**](PoolsAPI.md#GetPool) | **Get** /pools/{id} | Fetch a pool. Its placements are available at GET /pools/{id}/placements.
[**ListPoolPlacements**](PoolsAPI.md#ListPoolPlacements) | **Get** /pools/{id}/placements | List a pool&#39;s placements (hypervisor + optional vmid).
[**ListPools**](PoolsAPI.md#ListPools) | **Get** /pools | List node pools in the caller&#39;s tenant.
[**ResizePool**](PoolsAPI.md#ResizePool) | **Patch** /pools/{id} | Resize a pool&#39;s desired count. Grow places new VMs (all-or-nothing); shrink removes newest placements (LIFO). Placements are available at GET /pools/{id}/placements.



## CreatePool

> PoolView CreatePool(ctx).CreatePoolInputBody(createPoolInputBody).Execute()

Create a node pool and place its VMs across hypervisors (anti-affinity spread, all-or-nothing). Placements are available at GET /pools/{id}/placements.

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
	createPoolInputBody := *openapiclient.NewCreatePoolInputBody("DatacenterId_example", int64(123), "Name_example", "SlotId_example") // CreatePoolInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoolsAPI.CreatePool(context.Background()).CreatePoolInputBody(createPoolInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.CreatePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreatePool`: PoolView
	fmt.Fprintf(os.Stdout, "Response from `PoolsAPI.CreatePool`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreatePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createPoolInputBody** | [**CreatePoolInputBody**](CreatePoolInputBody.md) |  | 

### Return type

[**PoolView**](PoolView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeletePool

> DeletePool(ctx, id).Execute()

Delete a pool and release all its placements.

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
	r, err := apiClient.PoolsAPI.DeletePool(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.DeletePool``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDeletePoolRequest struct via the builder pattern


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


## GetPool

> PoolView GetPool(ctx, id).Execute()

Fetch a pool. Its placements are available at GET /pools/{id}/placements.

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
	resp, r, err := apiClient.PoolsAPI.GetPool(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.GetPool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPool`: PoolView
	fmt.Fprintf(os.Stdout, "Response from `PoolsAPI.GetPool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**PoolView**](PoolView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPoolPlacements

> PlacementListOutputBody ListPoolPlacements(ctx, id).Execute()

List a pool's placements (hypervisor + optional vmid).

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
	resp, r, err := apiClient.PoolsAPI.ListPoolPlacements(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.ListPoolPlacements``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPoolPlacements`: PlacementListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `PoolsAPI.ListPoolPlacements`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiListPoolPlacementsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**PlacementListOutputBody**](PlacementListOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPools

> PoolListOutputBody ListPools(ctx).Execute()

List node pools in the caller's tenant.

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
	resp, r, err := apiClient.PoolsAPI.ListPools(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.ListPools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPools`: PoolListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `PoolsAPI.ListPools`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListPoolsRequest struct via the builder pattern


### Return type

[**PoolListOutputBody**](PoolListOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResizePool

> PoolView ResizePool(ctx, id).ResizePoolInputBody(resizePoolInputBody).Execute()

Resize a pool's desired count. Grow places new VMs (all-or-nothing); shrink removes newest placements (LIFO). Placements are available at GET /pools/{id}/placements.

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
	resizePoolInputBody := *openapiclient.NewResizePoolInputBody(int64(123)) // ResizePoolInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoolsAPI.ResizePool(context.Background(), id).ResizePoolInputBody(resizePoolInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.ResizePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ResizePool`: PoolView
	fmt.Fprintf(os.Stdout, "Response from `PoolsAPI.ResizePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiResizePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **resizePoolInputBody** | [**ResizePoolInputBody**](ResizePoolInputBody.md) |  | 

### Return type

[**PoolView**](PoolView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

