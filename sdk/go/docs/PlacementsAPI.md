# \PlacementsAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BackfillPlacementVmid**](PlacementsAPI.md#BackfillPlacementVmid) | **Patch** /placements/{id} | Attach the externally-assigned Proxmox vmid to a placement.
[**DeletePlacement**](PlacementsAPI.md#DeletePlacement) | **Delete** /placements/{id} | Remove a placement. The pool&#39;s desired_count is not adjusted; resize the pool to re-fill the vacancy.
[**GetPlacement**](PlacementsAPI.md#GetPlacement) | **Get** /placements/{id} | Fetch a single placement with its pool, hypervisor, and vmid.
[**ListPlacements**](PlacementsAPI.md#ListPlacements) | **Get** /placements | List all placements in the tenant with pool, slot, and hypervisor context (fleet overview).



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
	resp, r, err := apiClient.PlacementsAPI.BackfillPlacementVmid(context.Background(), id).BackfillVMIDInputBody(backfillVMIDInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementsAPI.BackfillPlacementVmid``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BackfillPlacementVmid`: PlacementView
	fmt.Fprintf(os.Stdout, "Response from `PlacementsAPI.BackfillPlacementVmid`: %v\n", resp)
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


## DeletePlacement

> DeletePlacement(ctx, id).Execute()

Remove a placement. The pool's desired_count is not adjusted; resize the pool to re-fill the vacancy.

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
	r, err := apiClient.PlacementsAPI.DeletePlacement(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementsAPI.DeletePlacement``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDeletePlacementRequest struct via the builder pattern


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


## GetPlacement

> PlacementView GetPlacement(ctx, id).Execute()

Fetch a single placement with its pool, hypervisor, and vmid.

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
	resp, r, err := apiClient.PlacementsAPI.GetPlacement(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementsAPI.GetPlacement``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPlacement`: PlacementView
	fmt.Fprintf(os.Stdout, "Response from `PlacementsAPI.GetPlacement`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPlacementRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**PlacementView**](PlacementView.md)

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
	resp, r, err := apiClient.PlacementsAPI.ListPlacements(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacementsAPI.ListPlacements``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPlacements`: FleetPlacementListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `PlacementsAPI.ListPlacements`: %v\n", resp)
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

