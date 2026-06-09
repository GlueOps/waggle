# \ApiKeysAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateApiKey**](ApiKeysAPI.md#CreateApiKey) | **Post** /api-keys | Mint an organization API key for automation (e.g. Terraform). The plaintext token is returned once.
[**ListApiKeys**](ApiKeysAPI.md#ListApiKeys) | **Get** /api-keys | List the organization&#39;s API keys (secrets are never returned).
[**RevokeApiKey**](ApiKeysAPI.md#RevokeApiKey) | **Delete** /api-keys/{id} | Revoke an organization API key. Idempotent from the caller&#39;s view.



## CreateApiKey

> CreateAPIKeyOutputBody CreateApiKey(ctx).CreateAPIKeyInputBody(createAPIKeyInputBody).Execute()

Mint an organization API key for automation (e.g. Terraform). The plaintext token is returned once.

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
	createAPIKeyInputBody := *openapiclient.NewCreateAPIKeyInputBody("Name_example") // CreateAPIKeyInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ApiKeysAPI.CreateApiKey(context.Background()).CreateAPIKeyInputBody(createAPIKeyInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ApiKeysAPI.CreateApiKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateApiKey`: CreateAPIKeyOutputBody
	fmt.Fprintf(os.Stdout, "Response from `ApiKeysAPI.CreateApiKey`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateApiKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createAPIKeyInputBody** | [**CreateAPIKeyInputBody**](CreateAPIKeyInputBody.md) |  | 

### Return type

[**CreateAPIKeyOutputBody**](CreateAPIKeyOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListApiKeys

> ListAPIKeysOutputBody ListApiKeys(ctx).Execute()

List the organization's API keys (secrets are never returned).

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
	resp, r, err := apiClient.ApiKeysAPI.ListApiKeys(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ApiKeysAPI.ListApiKeys``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListApiKeys`: ListAPIKeysOutputBody
	fmt.Fprintf(os.Stdout, "Response from `ApiKeysAPI.ListApiKeys`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListApiKeysRequest struct via the builder pattern


### Return type

[**ListAPIKeysOutputBody**](ListAPIKeysOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RevokeApiKey

> RevokeApiKey(ctx, id).Execute()

Revoke an organization API key. Idempotent from the caller's view.

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
	r, err := apiClient.ApiKeysAPI.RevokeApiKey(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ApiKeysAPI.RevokeApiKey``: %v\n", err)
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

Other parameters are passed through a pointer to a apiRevokeApiKeyRequest struct via the builder pattern


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

