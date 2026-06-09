# \FleetAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateDatacenter**](FleetAPI.md#CreateDatacenter) | **Post** /datacenters | Create a datacenter in the caller&#39;s tenant.
[**CreateHypervisor**](FleetAPI.md#CreateHypervisor) | **Post** /hypervisors | Create a hypervisor in the caller&#39;s tenant.
[**CreateSlot**](FleetAPI.md#CreateSlot) | **Post** /slots | Create a slot (t-shirt-size VM template) in the caller&#39;s tenant.
[**DeleteDatacenter**](FleetAPI.md#DeleteDatacenter) | **Delete** /datacenters/{id} | Delete a datacenter.
[**DeleteHypervisor**](FleetAPI.md#DeleteHypervisor) | **Delete** /hypervisors/{id} | Delete a hypervisor.
[**DeleteSlot**](FleetAPI.md#DeleteSlot) | **Delete** /slots/{id} | Delete a slot.
[**DiscoverHypervisors**](FleetAPI.md#DiscoverHypervisors) | **Post** /datacenters/{id}/discover | Discover hypervisors from the datacenter&#39;s Proxmox cluster and upsert them (preserving reserved capacity and schedulable). Set async to run in the background.
[**GetDatacenter**](FleetAPI.md#GetDatacenter) | **Get** /datacenters/{id} | Fetch a datacenter by ID.
[**GetHypervisor**](FleetAPI.md#GetHypervisor) | **Get** /hypervisors/{id} | Fetch a hypervisor by ID.
[**GetSlot**](FleetAPI.md#GetSlot) | **Get** /slots/{id} | Fetch a slot by ID.
[**ListDatacenters**](FleetAPI.md#ListDatacenters) | **Get** /datacenters | List datacenters in the caller&#39;s tenant.
[**ListHypervisors**](FleetAPI.md#ListHypervisors) | **Get** /hypervisors | List hypervisors in the caller&#39;s tenant.
[**ListSlots**](FleetAPI.md#ListSlots) | **Get** /slots | List slots in the caller&#39;s tenant.
[**UpdateDatacenter**](FleetAPI.md#UpdateDatacenter) | **Put** /datacenters/{id} | Update a datacenter.
[**UpdateHypervisor**](FleetAPI.md#UpdateHypervisor) | **Put** /hypervisors/{id} | Update a hypervisor.
[**UpdateSlot**](FleetAPI.md#UpdateSlot) | **Put** /slots/{id} | Update a slot.



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
	resp, r, err := apiClient.FleetAPI.CreateDatacenter(context.Background()).DatacenterBody(datacenterBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.CreateDatacenter``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateDatacenter`: DatacenterView
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.CreateDatacenter`: %v\n", resp)
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
	resp, r, err := apiClient.FleetAPI.CreateHypervisor(context.Background()).HypervisorBody(hypervisorBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.CreateHypervisor``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateHypervisor`: HypervisorView
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.CreateHypervisor`: %v\n", resp)
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


## CreateSlot

> SlotView CreateSlot(ctx).SlotBody(slotBody).Execute()

Create a slot (t-shirt-size VM template) in the caller's tenant.

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
	slotBody := *openapiclient.NewSlotBody(int64(123), "Name_example", int64(123), int64(123)) // SlotBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.FleetAPI.CreateSlot(context.Background()).SlotBody(slotBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.CreateSlot``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateSlot`: SlotView
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.CreateSlot`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateSlotRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **slotBody** | [**SlotBody**](SlotBody.md) |  | 

### Return type

[**SlotView**](SlotView.md)

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
	r, err := apiClient.FleetAPI.DeleteDatacenter(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.DeleteDatacenter``: %v\n", err)
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
	r, err := apiClient.FleetAPI.DeleteHypervisor(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.DeleteHypervisor``: %v\n", err)
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


## DeleteSlot

> DeleteSlot(ctx, id).Execute()

Delete a slot.

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
	r, err := apiClient.FleetAPI.DeleteSlot(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.DeleteSlot``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDeleteSlotRequest struct via the builder pattern


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
	resp, r, err := apiClient.FleetAPI.DiscoverHypervisors(context.Background(), id).DiscoverInputBody(discoverInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.DiscoverHypervisors``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DiscoverHypervisors`: DiscoverOutputBody
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.DiscoverHypervisors`: %v\n", resp)
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
	resp, r, err := apiClient.FleetAPI.GetDatacenter(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.GetDatacenter``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDatacenter`: DatacenterView
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.GetDatacenter`: %v\n", resp)
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
	resp, r, err := apiClient.FleetAPI.GetHypervisor(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.GetHypervisor``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetHypervisor`: HypervisorView
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.GetHypervisor`: %v\n", resp)
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


## GetSlot

> SlotView GetSlot(ctx, id).Execute()

Fetch a slot by ID.

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
	resp, r, err := apiClient.FleetAPI.GetSlot(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.GetSlot``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSlot`: SlotView
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.GetSlot`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSlotRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**SlotView**](SlotView.md)

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
	resp, r, err := apiClient.FleetAPI.ListDatacenters(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.ListDatacenters``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListDatacenters`: DatacenterListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.ListDatacenters`: %v\n", resp)
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
	resp, r, err := apiClient.FleetAPI.ListHypervisors(context.Background()).DatacenterId(datacenterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.ListHypervisors``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListHypervisors`: HypervisorListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.ListHypervisors`: %v\n", resp)
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


## ListSlots

> SlotListOutputBody ListSlots(ctx).Execute()

List slots in the caller's tenant.

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
	resp, r, err := apiClient.FleetAPI.ListSlots(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.ListSlots``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListSlots`: SlotListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.ListSlots`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListSlotsRequest struct via the builder pattern


### Return type

[**SlotListOutputBody**](SlotListOutputBody.md)

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
	resp, r, err := apiClient.FleetAPI.UpdateDatacenter(context.Background(), id).DatacenterBody(datacenterBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.UpdateDatacenter``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateDatacenter`: DatacenterView
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.UpdateDatacenter`: %v\n", resp)
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
	resp, r, err := apiClient.FleetAPI.UpdateHypervisor(context.Background(), id).HypervisorBody(hypervisorBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.UpdateHypervisor``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateHypervisor`: HypervisorView
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.UpdateHypervisor`: %v\n", resp)
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


## UpdateSlot

> SlotView UpdateSlot(ctx, id).SlotBody(slotBody).Execute()

Update a slot.

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
	slotBody := *openapiclient.NewSlotBody(int64(123), "Name_example", int64(123), int64(123)) // SlotBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.FleetAPI.UpdateSlot(context.Background(), id).SlotBody(slotBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FleetAPI.UpdateSlot``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateSlot`: SlotView
	fmt.Fprintf(os.Stdout, "Response from `FleetAPI.UpdateSlot`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateSlotRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **slotBody** | [**SlotBody**](SlotBody.md) |  | 

### Return type

[**SlotView**](SlotView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

