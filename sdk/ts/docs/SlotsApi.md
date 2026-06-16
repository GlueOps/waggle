# SlotsApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createSlot**](SlotsApi.md#createslot) | **POST** /slots | Create a slot (t-shirt-size VM template) in the caller\&#39;s tenant. |
| [**deleteSlot**](SlotsApi.md#deleteslot) | **DELETE** /slots/{id} | Delete a slot. |
| [**getSlot**](SlotsApi.md#getslot) | **GET** /slots/{id} | Fetch a slot by ID. |
| [**listSlots**](SlotsApi.md#listslots) | **GET** /slots | List slots in the caller\&#39;s tenant. |
| [**updateSlot**](SlotsApi.md#updateslot) | **PUT** /slots/{id} | Update a slot. |



## createSlot

> SlotView createSlot(slotBody)

Create a slot (t-shirt-size VM template) in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  SlotsApi,
} from '@glueops/waggle-sdk';
import type { CreateSlotRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new SlotsApi(config);

  const body = {
    // SlotBody
    slotBody: ...,
  } satisfies CreateSlotRequest;

  try {
    const data = await api.createSlot(body);
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
| **slotBody** | [SlotBody](SlotBody.md) |  | |

### Return type

[**SlotView**](SlotView.md)

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


## deleteSlot

> deleteSlot(id)

Delete a slot.

### Example

```ts
import {
  Configuration,
  SlotsApi,
} from '@glueops/waggle-sdk';
import type { DeleteSlotRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new SlotsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies DeleteSlotRequest;

  try {
    const data = await api.deleteSlot(body);
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


## getSlot

> SlotView getSlot(id)

Fetch a slot by ID.

### Example

```ts
import {
  Configuration,
  SlotsApi,
} from '@glueops/waggle-sdk';
import type { GetSlotRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new SlotsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies GetSlotRequest;

  try {
    const data = await api.getSlot(body);
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

[**SlotView**](SlotView.md)

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


## listSlots

> SlotListOutputBody listSlots(name)

List slots in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  SlotsApi,
} from '@glueops/waggle-sdk';
import type { ListSlotsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new SlotsApi(config);

  const body = {
    // string | Filter to the slot with this exact (unique) name. (optional)
    name: name_example,
  } satisfies ListSlotsRequest;

  try {
    const data = await api.listSlots(body);
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
| **name** | `string` | Filter to the slot with this exact (unique) name. | [Optional] [Defaults to `undefined`] |

### Return type

[**SlotListOutputBody**](SlotListOutputBody.md)

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


## updateSlot

> SlotView updateSlot(id, slotBody)

Update a slot.

### Example

```ts
import {
  Configuration,
  SlotsApi,
} from '@glueops/waggle-sdk';
import type { UpdateSlotRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new SlotsApi(config);

  const body = {
    // string
    id: id_example,
    // SlotBody
    slotBody: ...,
  } satisfies UpdateSlotRequest;

  try {
    const data = await api.updateSlot(body);
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
| **slotBody** | [SlotBody](SlotBody.md) |  | |

### Return type

[**SlotView**](SlotView.md)

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

