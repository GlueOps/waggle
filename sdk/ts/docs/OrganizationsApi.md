# OrganizationsApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**addMember**](OrganizationsApi.md#addmember) | **POST** /organizations/{id}/members | Add or invite a member by email (admin+; owner required to grant owner). Unknown emails get an invite link. |
| [**createOrg**](OrganizationsApi.md#createorg) | **POST** /organizations | Create an organization (you become its owner) and enqueue tenant provisioning. |
| [**deleteOrg**](OrganizationsApi.md#deleteorg) | **DELETE** /organizations/{id} | Delete an organization and enqueue tenant teardown (owner only). |
| [**getOrg**](OrganizationsApi.md#getorg) | **GET** /organizations/{id} | Get an organization the caller belongs to. |
| [**listMembers**](OrganizationsApi.md#listmembers) | **GET** /organizations/{id}/members | List an organization\&#39;s members. |
| [**listOrgs**](OrganizationsApi.md#listorgs) | **GET** /organizations | List the organizations the caller belongs to (with their role). |
| [**removeMember**](OrganizationsApi.md#removemember) | **DELETE** /organizations/{id}/members/{userId} | Remove a member (admin+; owner required to remove owners; never the last owner). |
| [**updateMember**](OrganizationsApi.md#updatemember) | **PATCH** /organizations/{id}/members/{userId} | Change a member\&#39;s role (admin+; owner required to touch owners). |
| [**updateOrg**](OrganizationsApi.md#updateorg) | **PATCH** /organizations/{id} | Rename an organization (admin or owner). |



## addMember

> AddMemberOutputBody addMember(id, addMemberInputBody)

Add or invite a member by email (admin+; owner required to grant owner). Unknown emails get an invite link.

### Example

```ts
import {
  Configuration,
  OrganizationsApi,
} from '@glueops/waggle-sdk';
import type { AddMemberRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new OrganizationsApi(config);

  const body = {
    // string
    id: id_example,
    // AddMemberInputBody
    addMemberInputBody: ...,
  } satisfies AddMemberRequest;

  try {
    const data = await api.addMember(body);
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
| **addMemberInputBody** | [AddMemberInputBody](AddMemberInputBody.md) |  | |

### Return type

[**AddMemberOutputBody**](AddMemberOutputBody.md)

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


## createOrg

> OrgFullView createOrg(createOrgInputBody)

Create an organization (you become its owner) and enqueue tenant provisioning.

### Example

```ts
import {
  Configuration,
  OrganizationsApi,
} from '@glueops/waggle-sdk';
import type { CreateOrgRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new OrganizationsApi(config);

  const body = {
    // CreateOrgInputBody
    createOrgInputBody: ...,
  } satisfies CreateOrgRequest;

  try {
    const data = await api.createOrg(body);
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
| **createOrgInputBody** | [CreateOrgInputBody](CreateOrgInputBody.md) |  | |

### Return type

[**OrgFullView**](OrgFullView.md)

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


## deleteOrg

> deleteOrg(id)

Delete an organization and enqueue tenant teardown (owner only).

### Example

```ts
import {
  Configuration,
  OrganizationsApi,
} from '@glueops/waggle-sdk';
import type { DeleteOrgRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new OrganizationsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies DeleteOrgRequest;

  try {
    const data = await api.deleteOrg(body);
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


## getOrg

> OrgFullView getOrg(id)

Get an organization the caller belongs to.

### Example

```ts
import {
  Configuration,
  OrganizationsApi,
} from '@glueops/waggle-sdk';
import type { GetOrgRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new OrganizationsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies GetOrgRequest;

  try {
    const data = await api.getOrg(body);
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

[**OrgFullView**](OrgFullView.md)

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


## listMembers

> MemberListOutputBody listMembers(id)

List an organization\&#39;s members.

### Example

```ts
import {
  Configuration,
  OrganizationsApi,
} from '@glueops/waggle-sdk';
import type { ListMembersRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new OrganizationsApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies ListMembersRequest;

  try {
    const data = await api.listMembers(body);
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

[**MemberListOutputBody**](MemberListOutputBody.md)

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


## listOrgs

> OrgListOutputBody listOrgs()

List the organizations the caller belongs to (with their role).

### Example

```ts
import {
  Configuration,
  OrganizationsApi,
} from '@glueops/waggle-sdk';
import type { ListOrgsRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new OrganizationsApi(config);

  try {
    const data = await api.listOrgs();
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

[**OrgListOutputBody**](OrgListOutputBody.md)

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


## removeMember

> removeMember(id, userId)

Remove a member (admin+; owner required to remove owners; never the last owner).

### Example

```ts
import {
  Configuration,
  OrganizationsApi,
} from '@glueops/waggle-sdk';
import type { RemoveMemberRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new OrganizationsApi(config);

  const body = {
    // string
    id: id_example,
    // string
    userId: userId_example,
  } satisfies RemoveMemberRequest;

  try {
    const data = await api.removeMember(body);
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
| **userId** | `string` |  | [Defaults to `undefined`] |

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


## updateMember

> MemberJSONView updateMember(id, userId, updateMemberInputBody)

Change a member\&#39;s role (admin+; owner required to touch owners).

### Example

```ts
import {
  Configuration,
  OrganizationsApi,
} from '@glueops/waggle-sdk';
import type { UpdateMemberRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new OrganizationsApi(config);

  const body = {
    // string
    id: id_example,
    // string
    userId: userId_example,
    // UpdateMemberInputBody
    updateMemberInputBody: ...,
  } satisfies UpdateMemberRequest;

  try {
    const data = await api.updateMember(body);
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
| **userId** | `string` |  | [Defaults to `undefined`] |
| **updateMemberInputBody** | [UpdateMemberInputBody](UpdateMemberInputBody.md) |  | |

### Return type

[**MemberJSONView**](MemberJSONView.md)

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


## updateOrg

> OrgFullView updateOrg(id, updateOrgInputBody)

Rename an organization (admin or owner).

### Example

```ts
import {
  Configuration,
  OrganizationsApi,
} from '@glueops/waggle-sdk';
import type { UpdateOrgRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new OrganizationsApi(config);

  const body = {
    // string
    id: id_example,
    // UpdateOrgInputBody
    updateOrgInputBody: ...,
  } satisfies UpdateOrgRequest;

  try {
    const data = await api.updateOrg(body);
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
| **updateOrgInputBody** | [UpdateOrgInputBody](UpdateOrgInputBody.md) |  | |

### Return type

[**OrgFullView**](OrgFullView.md)

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

