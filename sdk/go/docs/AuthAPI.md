# \AuthAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AuthAcceptInvite**](AuthAPI.md#AuthAcceptInvite) | **Post** /auth/accept-invite | Accept an organization invite: set a password (if new) and sign in to the org.
[**AuthLogin**](AuthAPI.md#AuthLogin) | **Post** /auth/login | Exchange credentials for an access + refresh token. Returns membership list when no organization_id is given and multiple memberships exist.
[**AuthLogout**](AuthAPI.md#AuthLogout) | **Post** /auth/logout | Revoke the supplied refresh token&#39;s session. Idempotent.
[**AuthMe**](AuthAPI.md#AuthMe) | **Get** /auth/me | Return the authenticated account, its verified-or-pending emails, and its org memberships.
[**AuthRefresh**](AuthAPI.md#AuthRefresh) | **Post** /auth/refresh | Rotate a refresh token, returning a new access + refresh pair.
[**AuthSignup**](AuthAPI.md#AuthSignup) | **Post** /auth/signup | Create an account, organization, and first user; enqueue tenant provisioning.
[**AuthSwitchOrg**](AuthAPI.md#AuthSwitchOrg) | **Post** /auth/switch | Issue a new token pair scoped to another organization the account belongs to.
[**AuthVerifyEmail**](AuthAPI.md#AuthVerifyEmail) | **Post** /auth/verify-email | Consume a verification token to mark an email address verified. Idempotent.



## AuthAcceptInvite

> LoginOutputBody AuthAcceptInvite(ctx).AcceptInviteInputBody(acceptInviteInputBody).Execute()

Accept an organization invite: set a password (if new) and sign in to the org.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	acceptInviteInputBody := *openapiclient.NewAcceptInviteInputBody("Token_example") // AcceptInviteInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.AuthAcceptInvite(context.Background()).AcceptInviteInputBody(acceptInviteInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.AuthAcceptInvite``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AuthAcceptInvite`: LoginOutputBody
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.AuthAcceptInvite`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthAcceptInviteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **acceptInviteInputBody** | [**AcceptInviteInputBody**](AcceptInviteInputBody.md) |  | 

### Return type

[**LoginOutputBody**](LoginOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AuthLogin

> LoginOutputBody AuthLogin(ctx).LoginInputBody(loginInputBody).Execute()

Exchange credentials for an access + refresh token. Returns membership list when no organization_id is given and multiple memberships exist.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	loginInputBody := *openapiclient.NewLoginInputBody("Email_example", "Password_example") // LoginInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.AuthLogin(context.Background()).LoginInputBody(loginInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.AuthLogin``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AuthLogin`: LoginOutputBody
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.AuthLogin`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthLoginRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginInputBody** | [**LoginInputBody**](LoginInputBody.md) |  | 

### Return type

[**LoginOutputBody**](LoginOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AuthLogout

> AuthLogout(ctx).LogoutInputBody(logoutInputBody).Execute()

Revoke the supplied refresh token's session. Idempotent.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	logoutInputBody := *openapiclient.NewLogoutInputBody("RefreshToken_example") // LogoutInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AuthAPI.AuthLogout(context.Background()).LogoutInputBody(logoutInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.AuthLogout``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthLogoutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **logoutInputBody** | [**LogoutInputBody**](LogoutInputBody.md) |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AuthMe

> MeOutputBody AuthMe(ctx).Execute()

Return the authenticated account, its verified-or-pending emails, and its org memberships.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.AuthMe(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.AuthMe``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AuthMe`: MeOutputBody
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.AuthMe`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiAuthMeRequest struct via the builder pattern


### Return type

[**MeOutputBody**](MeOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AuthRefresh

> RefreshOutputBody AuthRefresh(ctx).RefreshInputBody(refreshInputBody).Execute()

Rotate a refresh token, returning a new access + refresh pair.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	refreshInputBody := *openapiclient.NewRefreshInputBody("RefreshToken_example") // RefreshInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.AuthRefresh(context.Background()).RefreshInputBody(refreshInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.AuthRefresh``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AuthRefresh`: RefreshOutputBody
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.AuthRefresh`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthRefreshRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **refreshInputBody** | [**RefreshInputBody**](RefreshInputBody.md) |  | 

### Return type

[**RefreshOutputBody**](RefreshOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AuthSignup

> SignupOutputBody AuthSignup(ctx).SignupInputBody(signupInputBody).Execute()

Create an account, organization, and first user; enqueue tenant provisioning.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	signupInputBody := *openapiclient.NewSignupInputBody("Email_example", "OrganizationName_example", "Password_example") // SignupInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.AuthSignup(context.Background()).SignupInputBody(signupInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.AuthSignup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AuthSignup`: SignupOutputBody
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.AuthSignup`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthSignupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **signupInputBody** | [**SignupInputBody**](SignupInputBody.md) |  | 

### Return type

[**SignupOutputBody**](SignupOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AuthSwitchOrg

> LoginOutputBody AuthSwitchOrg(ctx).SwitchOrgInputBody(switchOrgInputBody).Execute()

Issue a new token pair scoped to another organization the account belongs to.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	switchOrgInputBody := *openapiclient.NewSwitchOrgInputBody("OrganizationId_example") // SwitchOrgInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.AuthSwitchOrg(context.Background()).SwitchOrgInputBody(switchOrgInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.AuthSwitchOrg``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AuthSwitchOrg`: LoginOutputBody
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.AuthSwitchOrg`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthSwitchOrgRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **switchOrgInputBody** | [**SwitchOrgInputBody**](SwitchOrgInputBody.md) |  | 

### Return type

[**LoginOutputBody**](LoginOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AuthVerifyEmail

> AuthVerifyEmail(ctx).VerifyEmailInputBody(verifyEmailInputBody).Execute()

Consume a verification token to mark an email address verified. Idempotent.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/waggle/sdk/go/waggle"
)

func main() {
	verifyEmailInputBody := *openapiclient.NewVerifyEmailInputBody("Token_example") // VerifyEmailInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AuthAPI.AuthVerifyEmail(context.Background()).VerifyEmailInputBody(verifyEmailInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.AuthVerifyEmail``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthVerifyEmailRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **verifyEmailInputBody** | [**VerifyEmailInputBody**](VerifyEmailInputBody.md) |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

