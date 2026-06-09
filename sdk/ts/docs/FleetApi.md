# FleetApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createDatacenter**](FleetApi.md#createdatacenter) | **POST** /datacenters | Create a datacenter in the caller\&#39;s tenant. |
| [**createHypervisor**](FleetApi.md#createhypervisor) | **POST** /hypervisors | Create a hypervisor in the caller\&#39;s tenant. |
| [**createSlot**](FleetApi.md#createslot) | **POST** /slots | Create a slot (t-shirt-size VM template) in the caller\&#39;s tenant. |
| [**deleteDatacenter**](FleetApi.md#deletedatacenter) | **DELETE** /datacenters/{id} | Delete a datacenter. |
| [**deleteHypervisor**](FleetApi.md#deletehypervisor) | **DELETE** /hypervisors/{id} | Delete a hypervisor. |
| [**deleteSlot**](FleetApi.md#deleteslot) | **DELETE** /slots/{id} | Delete a slot. |
| [**discoverHypervisors**](FleetApi.md#discoverhypervisors) | **POST** /datacenters/{id}/discover | Discover hypervisors from the datacenter\&#39;s Proxmox cluster and upsert them (preserving reserved capacity and schedulable). Set async to run in the background. |
| [**getDatacenter**](FleetApi.md#getdatacenter) | **GET** /datacenters/{id} | Fetch a datacenter by ID. |
| [**getHypervisor**](FleetApi.md#gethypervisor) | **GET** /hypervisors/{id} | Fetch a hypervisor by ID. |
| [**getSlot**](FleetApi.md#getslot) | **GET** /slots/{id} | Fetch a slot by ID. |
| [**listDatacenters**](FleetApi.md#listdatacenters) | **GET** /datacenters | List datacenters in the caller\&#39;s tenant. |
| [**listHypervisors**](FleetApi.md#listhypervisors) | **GET** /hypervisors | List hypervisors in the caller\&#39;s tenant. |
| [**listSlots**](FleetApi.md#listslots) | **GET** /slots | List slots in the caller\&#39;s tenant. |
| [**updateDatacenter**](FleetApi.md#updatedatacenter) | **PUT** /datacenters/{id} | Update a datacenter. |
| [**updateHypervisor**](FleetApi.md#updatehypervisor) | **PUT** /hypervisors/{id} | Update a hypervisor. |
| [**updateSlot**](FleetApi.md#updateslot) | **PUT** /slots/{id} | Update a slot. |



## createDatacenter

> DatacenterView createDatacenter(datacenterBody)

Create a datacenter in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { CreateDatacenterRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

  const body = {
    // DatacenterBody
    datacenterBody: ...,
  } satisfies CreateDatacenterRequest;

  try {
    const data = await api.createDatacenter(body);
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
| **datacenterBody** | [DatacenterBody](DatacenterBody.md) |  | |

### Return type

[**DatacenterView**](DatacenterView.md)

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


## createHypervisor

> HypervisorView createHypervisor(hypervisorBody)

Create a hypervisor in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { CreateHypervisorRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

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


## createSlot

> SlotView createSlot(slotBody)

Create a slot (t-shirt-size VM template) in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { CreateSlotRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

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


## deleteDatacenter

> deleteDatacenter(id)

Delete a datacenter.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { DeleteDatacenterRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies DeleteDatacenterRequest;

  try {
    const data = await api.deleteDatacenter(body);
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


## deleteHypervisor

> deleteHypervisor(id)

Delete a hypervisor.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { DeleteHypervisorRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

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


## deleteSlot

> deleteSlot(id)

Delete a slot.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { DeleteSlotRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

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


## discoverHypervisors

> DiscoverOutputBody discoverHypervisors(id, discoverInputBody)

Discover hypervisors from the datacenter\&#39;s Proxmox cluster and upsert them (preserving reserved capacity and schedulable). Set async to run in the background.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { DiscoverHypervisorsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

  const body = {
    // string
    id: id_example,
    // DiscoverInputBody
    discoverInputBody: ...,
  } satisfies DiscoverHypervisorsRequest;

  try {
    const data = await api.discoverHypervisors(body);
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
| **discoverInputBody** | [DiscoverInputBody](DiscoverInputBody.md) |  | |

### Return type

[**DiscoverOutputBody**](DiscoverOutputBody.md)

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


## getDatacenter

> DatacenterView getDatacenter(id)

Fetch a datacenter by ID.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { GetDatacenterRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies GetDatacenterRequest;

  try {
    const data = await api.getDatacenter(body);
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

[**DatacenterView**](DatacenterView.md)

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


## getHypervisor

> HypervisorView getHypervisor(id)

Fetch a hypervisor by ID.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { GetHypervisorRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

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


## getSlot

> SlotView getSlot(id)

Fetch a slot by ID.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { GetSlotRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

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


## listDatacenters

> DatacenterListOutputBody listDatacenters()

List datacenters in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { ListDatacentersRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

  try {
    const data = await api.listDatacenters();
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

[**DatacenterListOutputBody**](DatacenterListOutputBody.md)

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
  FleetApi,
} from '@glueops/waggle-sdk';
import type { ListHypervisorsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

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


## listSlots

> SlotListOutputBody listSlots()

List slots in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { ListSlotsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

  try {
    const data = await api.listSlots();
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


## updateDatacenter

> DatacenterView updateDatacenter(id, datacenterBody)

Update a datacenter.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { UpdateDatacenterRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

  const body = {
    // string
    id: id_example,
    // DatacenterBody
    datacenterBody: ...,
  } satisfies UpdateDatacenterRequest;

  try {
    const data = await api.updateDatacenter(body);
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
| **datacenterBody** | [DatacenterBody](DatacenterBody.md) |  | |

### Return type

[**DatacenterView**](DatacenterView.md)

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


## updateHypervisor

> HypervisorView updateHypervisor(id, hypervisorBody)

Update a hypervisor.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { UpdateHypervisorRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

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


## updateSlot

> SlotView updateSlot(id, slotBody)

Update a slot.

### Example

```ts
import {
  Configuration,
  FleetApi,
} from '@glueops/waggle-sdk';
import type { UpdateSlotRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FleetApi(config);

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

