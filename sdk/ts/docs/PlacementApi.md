# PlacementApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**backfillPlacementVmid**](PlacementApi.md#backfillplacementvmid) | **PATCH** /placements/{id} | Attach the externally-assigned Proxmox vmid to a placement. |
| [**listPlacements**](PlacementApi.md#listplacements) | **GET** /placements | List all placements in the tenant with pool, slot, and hypervisor context (fleet overview). |



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

