# DatacentersApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createDatacenter**](DatacentersApi.md#createdatacenter) | **POST** /datacenters | Create a datacenter in the caller\&#39;s tenant. |
| [**deleteDatacenter**](DatacentersApi.md#deletedatacenter) | **DELETE** /datacenters/{id} | Delete a datacenter. |
| [**discoverHypervisors**](DatacentersApi.md#discoverhypervisors) | **POST** /datacenters/{id}/discover | Discover hypervisors from the datacenter\&#39;s Proxmox cluster and upsert them (preserving reserved capacity and schedulable). Set async to run in the background. |
| [**getDatacenter**](DatacentersApi.md#getdatacenter) | **GET** /datacenters/{id} | Fetch a datacenter by ID. |
| [**listDatacenters**](DatacentersApi.md#listdatacenters) | **GET** /datacenters | List datacenters in the caller\&#39;s tenant. |
| [**updateDatacenter**](DatacentersApi.md#updatedatacenter) | **PUT** /datacenters/{id} | Update a datacenter. |



## createDatacenter

> DatacenterView createDatacenter(datacenterBody)

Create a datacenter in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  DatacentersApi,
} from '@glueops/waggle-sdk';
import type { CreateDatacenterRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new DatacentersApi(config);

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


## deleteDatacenter

> deleteDatacenter(id)

Delete a datacenter.

### Example

```ts
import {
  Configuration,
  DatacentersApi,
} from '@glueops/waggle-sdk';
import type { DeleteDatacenterRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new DatacentersApi(config);

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


## discoverHypervisors

> DiscoverOutputBody discoverHypervisors(id, discoverInputBody)

Discover hypervisors from the datacenter\&#39;s Proxmox cluster and upsert them (preserving reserved capacity and schedulable). Set async to run in the background.

### Example

```ts
import {
  Configuration,
  DatacentersApi,
} from '@glueops/waggle-sdk';
import type { DiscoverHypervisorsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new DatacentersApi(config);

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
  DatacentersApi,
} from '@glueops/waggle-sdk';
import type { GetDatacenterRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new DatacentersApi(config);

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


## listDatacenters

> DatacenterListOutputBody listDatacenters()

List datacenters in the caller\&#39;s tenant.

### Example

```ts
import {
  Configuration,
  DatacentersApi,
} from '@glueops/waggle-sdk';
import type { ListDatacentersRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new DatacentersApi(config);

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


## updateDatacenter

> DatacenterView updateDatacenter(id, datacenterBody)

Update a datacenter.

### Example

```ts
import {
  Configuration,
  DatacentersApi,
} from '@glueops/waggle-sdk';
import type { UpdateDatacenterRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new DatacentersApi(config);

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

