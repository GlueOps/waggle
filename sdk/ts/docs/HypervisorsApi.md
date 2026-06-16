# HypervisorsApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createHypervisor**](HypervisorsApi.md#createhypervisor) | **POST** /hypervisors | Create a hypervisor in the caller\&#39;s tenant. |
| [**deleteHypervisor**](HypervisorsApi.md#deletehypervisor) | **DELETE** /hypervisors/{id} | Delete a hypervisor. |
| [**getHypervisor**](HypervisorsApi.md#gethypervisor) | **GET** /hypervisors/{id} | Fetch a hypervisor by ID. |
| [**listHypervisors**](HypervisorsApi.md#listhypervisors) | **GET** /hypervisors | List hypervisors in the caller\&#39;s tenant. |
| [**updateHypervisor**](HypervisorsApi.md#updatehypervisor) | **PUT** /hypervisors/{id} | Update a hypervisor. |



## createHypervisor

> HypervisorView createHypervisor(hypervisorBody)

Create a hypervisor in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  HypervisorsApi,
} from '@glueops/waggle-sdk';
import type { CreateHypervisorRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new HypervisorsApi(config);

  const body = {
    // HypervisorBody
    hypervisorBody: ...,
  } satisfies CreateHypervisorRequest;

  try {
    const data = await api.createHypervisor(body);
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
| **hypervisorBody** | [HypervisorBody](HypervisorBody.md) |  | |

### Return type

[**HypervisorView**](HypervisorView.md)

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


## deleteHypervisor

> deleteHypervisor(id)

Delete a hypervisor.

### Example

```ts
import {
  Configuration,
  HypervisorsApi,
} from '@glueops/waggle-sdk';
import type { DeleteHypervisorRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new HypervisorsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies DeleteHypervisorRequest;

  try {
    const data = await api.deleteHypervisor(body);
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


## getHypervisor

> HypervisorView getHypervisor(id)

Fetch a hypervisor by ID.

### Example

```ts
import {
  Configuration,
  HypervisorsApi,
} from '@glueops/waggle-sdk';
import type { GetHypervisorRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new HypervisorsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies GetHypervisorRequest;

  try {
    const data = await api.getHypervisor(body);
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

[**HypervisorView**](HypervisorView.md)

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


## listHypervisors

> HypervisorListOutputBody listHypervisors(datacenterId)

List hypervisors in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  HypervisorsApi,
} from '@glueops/waggle-sdk';
import type { ListHypervisorsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new HypervisorsApi(config);

  const body = {
    // string (optional)
    datacenterId: datacenterId_example,
  } satisfies ListHypervisorsRequest;

  try {
    const data = await api.listHypervisors(body);
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
| **datacenterId** | `string` |  | [Optional] [Defaults to `undefined`] |

### Return type

[**HypervisorListOutputBody**](HypervisorListOutputBody.md)

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


## updateHypervisor

> HypervisorView updateHypervisor(id, hypervisorBody)

Update a hypervisor.

### Example

```ts
import {
  Configuration,
  HypervisorsApi,
} from '@glueops/waggle-sdk';
import type { UpdateHypervisorRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new HypervisorsApi(config);

  const body = {
    // string
    id: id_example,
    // HypervisorBody
    hypervisorBody: ...,
  } satisfies UpdateHypervisorRequest;

  try {
    const data = await api.updateHypervisor(body);
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
| **hypervisorBody** | [HypervisorBody](HypervisorBody.md) |  | |

### Return type

[**HypervisorView**](HypervisorView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`, `application/problem+json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

