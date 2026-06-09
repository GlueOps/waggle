# \SystemAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Health**](SystemAPI.md#Health) | **Get** /health | Health check



## Health

> HealthOutputBody Health(ctx).Execute()

Health check

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
	resp, r, err := apiClient.SystemAPI.Health(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemAPI.Health``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Health`: HealthOutputBody
	fmt.Fprintf(os.Stdout, "Response from `SystemAPI.Health`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiHealthRequest struct via the builder pattern


### Return type

[**HealthOutputBody**](HealthOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

