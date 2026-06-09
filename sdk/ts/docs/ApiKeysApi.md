# ApiKeysApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createApiKey**](ApiKeysApi.md#createapikey) | **POST** /api-keys | Mint an organization API key for automation (e.g. Terraform). The plaintext token is returned once. |
| [**listApiKeys**](ApiKeysApi.md#listapikeys) | **GET** /api-keys | List the organization\&#39;s API keys (secrets are never returned). |
| [**revokeApiKey**](ApiKeysApi.md#revokeapikey) | **DELETE** /api-keys/{id} | Revoke an organization API key. Idempotent from the caller\&#39;s view. |



## createApiKey

> CreateAPIKeyOutputBody createApiKey(createAPIKeyInputBody)

Mint an organization API key for automation (e.g. Terraform). The plaintext token is returned once.

### Example

```ts
import {
  Configuration,
  ApiKeysApi,
} from '@glueops/waggle-sdk';
import type { CreateApiKeyRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ApiKeysApi(config);

  const body = {
    // CreateAPIKeyInputBody
    createAPIKeyInputBody: ...,
  } satisfies CreateApiKeyRequest;

  try {
    const data = await api.createApiKey(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **createAPIKeyInputBody** | [CreateAPIKeyInputBody](CreateAPIKeyInputBody.md) |  | |

### Return type

[**CreateAPIKeyOutputBody**](CreateAPIKeyOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`, `application/problem+json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listApiKeys

> ListAPIKeysOutputBody listApiKeys()

List the organization\&#39;s API keys (secrets are never returned).

### Example

```ts
import {
  Configuration,
  ApiKeysApi,
} from '@glueops/waggle-sdk';
import type { ListApiKeysRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ApiKeysApi(config);

  try {
    const data = await api.listApiKeys();
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

[**ListAPIKeysOutputBody**](ListAPIKeysOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`, `application/problem+json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## revokeApiKey

> revokeApiKey(id)

Revoke an organization API key. Idempotent from the caller\&#39;s view.

### Example

```ts
import {
  Configuration,
  ApiKeysApi,
} from '@glueops/waggle-sdk';
import type { RevokeApiKeyRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ApiKeysApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies RevokeApiKeyRequest;

  try {
    const data = await api.revokeApiKey(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `string` |  | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/problem+json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No Content |  -  |
| **0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

