# PlacementApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**backfillPlacementVmid**](PlacementApi.md#backfillplacementvmid) | **PATCH** /placements/{id} | Attach the externally-assigned Proxmox vmid to a placement. |
| [**createPool**](PlacementApi.md#createpool) | **POST** /pools | Create a node pool and place its VMs across hypervisors (anti-affinity spread, all-or-nothing). |
| [**deletePool**](PlacementApi.md#deletepool) | **DELETE** /pools/{id} | Delete a pool and release all its placements. |
| [**getPool**](PlacementApi.md#getpool) | **GET** /pools/{id} | Fetch a pool and its current placements. |
| [**listPlacements**](PlacementApi.md#listplacements) | **GET** /placements | List all placements in the tenant with pool, slot, and hypervisor context (fleet overview). |
| [**listPoolPlacements**](PlacementApi.md#listpoolplacements) | **GET** /pools/{id}/placements | List a pool\&#39;s placements (hypervisor + optional vmid). |
| [**listPools**](PlacementApi.md#listpools) | **GET** /pools | List node pools in the caller\&#39;s tenant. |
| [**resizePool**](PlacementApi.md#resizepool) | **PATCH** /pools/{id} | Resize a pool\&#39;s desired count. Grow places new VMs (all-or-nothing); shrink removes newest placements (LIFO). |



## backfillPlacementVmid

> PlacementView backfillPlacementVmid(id, backfillVMIDInputBody)

Attach the externally-assigned Proxmox vmid to a placement.

### Example

```ts
import {
  Configuration,
  PlacementApi,
} from '@glueops/waggle-sdk';
import type { BackfillPlacementVmidRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementApi(config);

  const body = {
    // string
    id: id_example,
    // BackfillVMIDInputBody
    backfillVMIDInputBody: ...,
  } satisfies BackfillPlacementVmidRequest;

  try {
    const data = await api.backfillPlacementVmid(body);
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
| **backfillVMIDInputBody** | [BackfillVMIDInputBody](BackfillVMIDInputBody.md) |  | |

### Return type

[**PlacementView**](PlacementView.md)

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


## createPool

> PoolResultBody createPool(createPoolInputBody)

Create a node pool and place its VMs across hypervisors (anti-affinity spread, all-or-nothing).

### Example

```ts
import {
  Configuration,
  PlacementApi,
} from '@glueops/waggle-sdk';
import type { CreatePoolRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementApi(config);

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

[**PoolResultBody**](PoolResultBody.md)

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
  PlacementApi,
} from '@glueops/waggle-sdk';
import type { DeletePoolRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementApi(config);

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

> PoolResultBody getPool(id)

Fetch a pool and its current placements.

### Example

```ts
import {
  Configuration,
  PlacementApi,
} from '@glueops/waggle-sdk';
import type { GetPoolRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementApi(config);

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

[**PoolResultBody**](PoolResultBody.md)

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


## listPlacements

> FleetPlacementListOutputBody listPlacements()

List all placements in the tenant with pool, slot, and hypervisor context (fleet overview).

### Example

```ts
import {
  Configuration,
  PlacementApi,
} from '@glueops/waggle-sdk';
import type { ListPlacementsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementApi(config);

  try {
    const data = await api.listPlacements();
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

[**FleetPlacementListOutputBody**](FleetPlacementListOutputBody.md)

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
  PlacementApi,
} from '@glueops/waggle-sdk';
import type { ListPoolPlacementsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementApi(config);

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
  PlacementApi,
} from '@glueops/waggle-sdk';
import type { ListPoolsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementApi(config);

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

> PoolResultBody resizePool(id, resizePoolInputBody)

Resize a pool\&#39;s desired count. Grow places new VMs (all-or-nothing); shrink removes newest placements (LIFO).

### Example

```ts
import {
  Configuration,
  PlacementApi,
} from '@glueops/waggle-sdk';
import type { ResizePoolRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementApi(config);

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

[**PoolResultBody**](PoolResultBody.md)

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

