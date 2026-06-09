# AuthApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**authAcceptInvite**](AuthApi.md#authacceptinvite) | **POST** /auth/accept-invite | Accept an organization invite: set a password (if new) and sign in to the org. |
| [**authLogin**](AuthApi.md#authlogin) | **POST** /auth/login | Exchange credentials for an access + refresh token. Returns membership list when no organization_id is given and multiple memberships exist. |
| [**authLogout**](AuthApi.md#authlogout) | **POST** /auth/logout | Revoke the supplied refresh token\&#39;s session. Idempotent. |
| [**authMe**](AuthApi.md#authme) | **GET** /auth/me | Return the authenticated account, its verified-or-pending emails, and its org memberships. |
| [**authRefresh**](AuthApi.md#authrefresh) | **POST** /auth/refresh | Rotate a refresh token, returning a new access + refresh pair. |
| [**authSignup**](AuthApi.md#authsignup) | **POST** /auth/signup | Create an account, organization, and first user; enqueue tenant provisioning. |
| [**authSwitchOrg**](AuthApi.md#authswitchorg) | **POST** /auth/switch | Issue a new token pair scoped to another organization the account belongs to. |
| [**authVerifyEmail**](AuthApi.md#authverifyemail) | **POST** /auth/verify-email | Consume a verification token to mark an email address verified. Idempotent. |



## authAcceptInvite

> LoginOutputBody authAcceptInvite(acceptInviteInputBody)

Accept an organization invite: set a password (if new) and sign in to the org.

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/waggle-sdk';
import type { AuthAcceptInviteRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const api = new AuthApi();

  const body = {
    // AcceptInviteInputBody
    acceptInviteInputBody: ...,
  } satisfies AuthAcceptInviteRequest;

  try {
    const data = await api.authAcceptInvite(body);
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
| **acceptInviteInputBody** | [AcceptInviteInputBody](AcceptInviteInputBody.md) |  | |

### Return type

[**LoginOutputBody**](LoginOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`, `application/problem+json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## authLogin

> LoginOutputBody authLogin(loginInputBody)

Exchange credentials for an access + refresh token. Returns membership list when no organization_id is given and multiple memberships exist.

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/waggle-sdk';
import type { AuthLoginRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const api = new AuthApi();

  const body = {
    // LoginInputBody
    loginInputBody: ...,
  } satisfies AuthLoginRequest;

  try {
    const data = await api.authLogin(body);
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
| **loginInputBody** | [LoginInputBody](LoginInputBody.md) |  | |

### Return type

[**LoginOutputBody**](LoginOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`, `application/problem+json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## authLogout

> authLogout(logoutInputBody)

Revoke the supplied refresh token\&#39;s session. Idempotent.

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/waggle-sdk';
import type { AuthLogoutRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const api = new AuthApi();

  const body = {
    // LogoutInputBody
    logoutInputBody: ...,
  } satisfies AuthLogoutRequest;

  try {
    const data = await api.authLogout(body);
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
| **logoutInputBody** | [LogoutInputBody](LogoutInputBody.md) |  | |

### Return type

`void` (Empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/problem+json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No Content |  -  |
| **0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## authMe

> MeOutputBody authMe()

Return the authenticated account, its verified-or-pending emails, and its org memberships.

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/waggle-sdk';
import type { AuthMeRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new AuthApi(config);

  try {
    const data = await api.authMe();
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

[**MeOutputBody**](MeOutputBody.md)

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


## authRefresh

> RefreshOutputBody authRefresh(refreshInputBody)

Rotate a refresh token, returning a new access + refresh pair.

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/waggle-sdk';
import type { AuthRefreshRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const api = new AuthApi();

  const body = {
    // RefreshInputBody
    refreshInputBody: ...,
  } satisfies AuthRefreshRequest;

  try {
    const data = await api.authRefresh(body);
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
| **refreshInputBody** | [RefreshInputBody](RefreshInputBody.md) |  | |

### Return type

[**RefreshOutputBody**](RefreshOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`, `application/problem+json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## authSignup

> SignupOutputBody authSignup(signupInputBody)

Create an account, organization, and first user; enqueue tenant provisioning.

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/waggle-sdk';
import type { AuthSignupRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const api = new AuthApi();

  const body = {
    // SignupInputBody
    signupInputBody: ...,
  } satisfies AuthSignupRequest;

  try {
    const data = await api.authSignup(body);
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
| **signupInputBody** | [SignupInputBody](SignupInputBody.md) |  | |

### Return type

[**SignupOutputBody**](SignupOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`, `application/problem+json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## authSwitchOrg

> LoginOutputBody authSwitchOrg(switchOrgInputBody)

Issue a new token pair scoped to another organization the account belongs to.

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/waggle-sdk';
import type { AuthSwitchOrgRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new AuthApi(config);

  const body = {
    // SwitchOrgInputBody
    switchOrgInputBody: ...,
  } satisfies AuthSwitchOrgRequest;

  try {
    const data = await api.authSwitchOrg(body);
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
| **switchOrgInputBody** | [SwitchOrgInputBody](SwitchOrgInputBody.md) |  | |

### Return type

[**LoginOutputBody**](LoginOutputBody.md)

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


## authVerifyEmail

> authVerifyEmail(verifyEmailInputBody)

Consume a verification token to mark an email address verified. Idempotent.

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/waggle-sdk';
import type { AuthVerifyEmailRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const api = new AuthApi();

  const body = {
    // VerifyEmailInputBody
    verifyEmailInputBody: ...,
  } satisfies AuthVerifyEmailRequest;

  try {
    const data = await api.authVerifyEmail(body);
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
| **verifyEmailInputBody** | [VerifyEmailInputBody](VerifyEmailInputBody.md) |  | |

### Return type

`void` (Empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/problem+json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No Content |  -  |
| **0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

