# \PlacementAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BackfillPlacementVmid**](PlacementAPI.md#BackfillPlacementVmid) | **Patch** /placements/{id} | Attach the externally-assigned Proxmox vmid to a placement.
[**ListPlacements**](PlacementAPI.md#ListPlacements) | **Get** /placements | List all placements in the tenant with pool, slot, and hypervisor context (fleet overview).



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

