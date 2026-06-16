# PoolsApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createPool**](PoolsApi.md#createpool) | **POST** /pools | Create a node pool and place its VMs across hypervisors (anti-affinity spread, all-or-nothing). Placements are available at GET /pools/{id}/placements. |
| [**deletePool**](PoolsApi.md#deletepool) | **DELETE** /pools/{id} | Delete a pool and release all its placements. |
| [**getPool**](PoolsApi.md#getpool) | **GET** /pools/{id} | Fetch a pool. Its placements are available at GET /pools/{id}/placements. |
| [**listPoolPlacements**](PoolsApi.md#listpoolplacements) | **GET** /pools/{id}/placements | List a pool\&#39;s placements (hypervisor + optional vmid). |
| [**listPools**](PoolsApi.md#listpools) | **GET** /pools | List node pools in the caller\&#39;s tenant. |
| [**resizePool**](PoolsApi.md#resizepool) | **PATCH** /pools/{id} | Resize a pool\&#39;s desired count. Grow places new VMs (all-or-nothing); shrink removes newest placements (LIFO). Placements are available at GET /pools/{id}/placements. |



## createPool

> PoolView createPool(createPoolInputBody)

Create a node pool and place its VMs across hypervisors (anti-affinity spread, all-or-nothing). Placements are available at GET /pools/{id}/placements.

### Example

```ts
import {
  Configuration,
  PoolsApi,
} from '@glueops/waggle-sdk';
import type { CreatePoolRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PoolsApi(config);

  const body = {
    // CreatePoolInputBody
    createPoolInputBody: ...,
  } satisfies CreatePoolRequest;

  try {
    const data = await api.createPool(body);
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
| **createPoolInputBody** | [CreatePoolInputBody](CreatePoolInputBody.md) |  | |

### Return type

[**PoolView**](PoolView.md)

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


## deletePool

> deletePool(id)

Delete a pool and release all its placements.

### Example

```ts
import {
  Configuration,
  PoolsApi,
} from '@glueops/waggle-sdk';
import type { DeletePoolRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PoolsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies DeletePoolRequest;

  try {
    const data = await api.deletePool(body);
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


## getPool

> PoolView getPool(id)

Fetch a pool. Its placements are available at GET /pools/{id}/placements.

### Example

```ts
import {
  Configuration,
  PoolsApi,
} from '@glueops/waggle-sdk';
import type { GetPoolRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PoolsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies GetPoolRequest;

  try {
    const data = await api.getPool(body);
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

[**PoolView**](PoolView.md)

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


## listPoolPlacements

> PlacementListOutputBody listPoolPlacements(id)

List a pool\&#39;s placements (hypervisor + optional vmid).

### Example

```ts
import {
  Configuration,
  PoolsApi,
} from '@glueops/waggle-sdk';
import type { ListPoolPlacementsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PoolsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies ListPoolPlacementsRequest;

  try {
    const data = await api.listPoolPlacements(body);
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

[**PlacementListOutputBody**](PlacementListOutputBody.md)

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


## listPools

> PoolListOutputBody listPools()

List node pools in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  PoolsApi,
} from '@glueops/waggle-sdk';
import type { ListPoolsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PoolsApi(config);

  try {
    const data = await api.listPools();
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

[**PoolListOutputBody**](PoolListOutputBody.md)

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


## resizePool

> PoolView resizePool(id, resizePoolInputBody)

Resize a pool\&#39;s desired count. Grow places new VMs (all-or-nothing); shrink removes newest placements (LIFO). Placements are available at GET /pools/{id}/placements.

### Example

```ts
import {
  Configuration,
  PoolsApi,
} from '@glueops/waggle-sdk';
import type { ResizePoolRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PoolsApi(config);

  const body = {
    // string
    id: id_example,
    // ResizePoolInputBody
    resizePoolInputBody: ...,
  } satisfies ResizePoolRequest;

  try {
    const data = await api.resizePool(body);
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
| **resizePoolInputBody** | [ResizePoolInputBody](ResizePoolInputBody.md) |  | |

### Return type

[**PoolView**](PoolView.md)

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

