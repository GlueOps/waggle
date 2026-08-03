# PlacementsApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**backfillPlacementVmid**](PlacementsApi.md#backfillplacementvmid) | **PATCH** /placements/{id} | Attach the externally-assigned Proxmox vmid to a placement. |
| [**deletePlacement**](PlacementsApi.md#deleteplacement) | **DELETE** /placements/{id} | Remove a placement. The pool\&#39;s desired_count is not adjusted; resize the pool to re-fill the vacancy. |
| [**getPlacement**](PlacementsApi.md#getplacement) | **GET** /placements/{id} | Fetch a single placement with its pool, hypervisor, and vmid. |
| [**listPlacements**](PlacementsApi.md#listplacements) | **GET** /placements | List all placements in the tenant with pool, slot, and hypervisor context (fleet overview). |



## backfillPlacementVmid

> PlacementView backfillPlacementVmid(id, backfillVMIDInputBody)

Attach the externally-assigned Proxmox vmid to a placement.

### Example

```ts
import {
  Configuration,
  PlacementsApi,
} from '@glueops/waggle-sdk';
import type { BackfillPlacementVmidRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementsApi(config);

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


## deletePlacement

> deletePlacement(id)

Remove a placement. The pool\&#39;s desired_count is not adjusted; resize the pool to re-fill the vacancy.

### Example

```ts
import {
  Configuration,
  PlacementsApi,
} from '@glueops/waggle-sdk';
import type { DeletePlacementRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies DeletePlacementRequest;

  try {
    const data = await api.deletePlacement(body);
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


## getPlacement

> PlacementView getPlacement(id)

Fetch a single placement with its pool, hypervisor, and vmid.

### Example

```ts
import {
  Configuration,
  PlacementsApi,
} from '@glueops/waggle-sdk';
import type { GetPlacementRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies GetPlacementRequest;

  try {
    const data = await api.getPlacement(body);
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

[**PlacementView**](PlacementView.md)

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
  PlacementsApi,
} from '@glueops/waggle-sdk';
import type { ListPlacementsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PlacementsApi(config);

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

