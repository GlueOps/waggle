# SystemApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**health**](SystemApi.md#health) | **GET** /health | Health check |



## health

> HealthOutputBody health()

Health check

### Example

```ts
import {
  Configuration,
  SystemApi,
} from '@glueops/waggle-sdk';
import type { HealthRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const api = new SystemApi();

  try {
    const data = await api.health();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**HealthOutputBody**](HealthOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`, `application/problem+json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

