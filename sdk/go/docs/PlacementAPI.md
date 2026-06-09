# \PlacementAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BackfillPlacementVmid**](PlacementAPI.md#BackfillPlacementVmid) | **Patch** /placements/{id} | Attach the externally-assigned Proxmox vmid to a placement.
[**CreatePool**](PlacementAPI.md#CreatePool) | **Post** /pools | Create a node pool and place its VMs across hypervisors (anti-affinity spread, all-or-nothing).
[**DeletePool**](PlacementAPI.md#DeletePool) | **Delete** /pools/{id} | Delete a pool and release all its placements.
[**GetPool**](PlacementAPI.md#GetPool) | **Get** /pools/{id} | Fetch a pool and its current placements.
[**ListPlacements**](PlacementAPI.md#ListPlacements) | **Get** /placements | List all placements in the tenant with pool, slot, and hypervisor context (fleet overview).
[**ListPoolPlacements**](PlacementAPI.md#ListPoolPlacements) | **Get** /pools/{id}/placements | List a pool&#39;s placements (hypervisor + optional vmid).
[**ListPools**](PlacementAPI.md#ListPools) | **Get** /pools | List node pools in the caller&#39;s tenant.
[**ResizePool**](PlacementAPI.md#ResizePool) | **Patch** /pools/{id} | Resize a pool&#39;s desired count. Grow places new VMs (all-or-nothing); shrink removes newest placements (LIFO).



## BackfillPlacementVmid

> PlacementView BackfillPlacementVmid(ctx, id).BackfillVMIDInputBody(backfillVMIDInputBody).Execute()

Attach the externally-assigned Proxmox vmid to a placement.

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
	backfillVMIDInputBody := *openapiclient.NewBackfillVMIDInputBody(int64(123)) // BackfillVMIDInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlacementAPI.BackfillPlacementVmid(context.Background(), id).BackfillVMIDInputBody(backfillVMIDInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementAPI.BackfillPlacementVmid``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BackfillPlacementVmid`: PlacementView
	fmt.Fprintf(os.Stdout, "Response from `PlacementAPI.BackfillPlacementVmid`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiBackfillPlacementVmidRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **backfillVMIDInputBody** | [**BackfillVMIDInputBody**](BackfillVMIDInputBody.md) |  | 

### Return type

[**PlacementView**](PlacementView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreatePool

> PoolResultBody CreatePool(ctx).CreatePoolInputBody(createPoolInputBody).Execute()

Create a node pool and place its VMs across hypervisors (anti-affinity spread, all-or-nothing).

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
	resp, r, err := apiClient.PlacementAPI.CreatePool(context.Background()).CreatePoolInputBody(createPoolInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementAPI.CreatePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreatePool`: PoolResultBody
	fmt.Fprintf(os.Stdout, "Response from `PlacementAPI.CreatePool`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreatePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createPoolInputBody** | [**CreatePoolInputBody**](CreatePoolInputBody.md) |  | 

### Return type

[**PoolResultBody**](PoolResultBody.md)

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
	r, err := apiClient.PlacementAPI.DeletePool(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementAPI.DeletePool``: %v\n", err)
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

> PoolResultBody GetPool(ctx, id).Execute()

Fetch a pool and its current placements.

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
	resp, r, err := apiClient.PlacementAPI.GetPool(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementAPI.GetPool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPool`: PoolResultBody
	fmt.Fprintf(os.Stdout, "Response from `PlacementAPI.GetPool`: %v\n", resp)
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

[**PoolResultBody**](PoolResultBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPlacements

> FleetPlacementListOutputBody ListPlacements(ctx).Execute()

List all placements in the tenant with pool, slot, and hypervisor context (fleet overview).

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
	resp, r, err := apiClient.PlacementAPI.ListPlacements(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementAPI.ListPlacements``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPlacements`: FleetPlacementListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `PlacementAPI.ListPlacements`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListPlacementsRequest struct via the builder pattern


### Return type

[**FleetPlacementListOutputBody**](FleetPlacementListOutputBody.md)

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
	resp, r, err := apiClient.PlacementAPI.ListPoolPlacements(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementAPI.ListPoolPlacements``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPoolPlacements`: PlacementListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `PlacementAPI.ListPoolPlacements`: %v\n", resp)
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
	resp, r, err := apiClient.PlacementAPI.ListPools(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementAPI.ListPools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPools`: PoolListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `PlacementAPI.ListPools`: %v\n", resp)
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

> PoolResultBody ResizePool(ctx, id).ResizePoolInputBody(resizePoolInputBody).Execute()

Resize a pool's desired count. Grow places new VMs (all-or-nothing); shrink removes newest placements (LIFO).

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
	resp, r, err := apiClient.PlacementAPI.ResizePool(context.Background(), id).ResizePoolInputBody(resizePoolInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementAPI.ResizePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ResizePool`: PoolResultBody
	fmt.Fprintf(os.Stdout, "Response from `PlacementAPI.ResizePool`: %v\n", resp)
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

[**PoolResultBody**](PoolResultBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

