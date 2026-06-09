# \OrganizationsAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddMember**](OrganizationsAPI.md#AddMember) | **Post** /organizations/{id}/members | Add or invite a member by email (admin+; owner required to grant owner). Unknown emails get an invite link.
[**CreateOrg**](OrganizationsAPI.md#CreateOrg) | **Post** /organizations | Create an organization (you become its owner) and enqueue tenant provisioning.
[**DeleteOrg**](OrganizationsAPI.md#DeleteOrg) | **Delete** /organizations/{id} | Delete an organization and enqueue tenant teardown (owner only).
[**GetOrg**](OrganizationsAPI.md#GetOrg) | **Get** /organizations/{id} | Get an organization the caller belongs to.
[**ListMembers**](OrganizationsAPI.md#ListMembers) | **Get** /organizations/{id}/members | List an organization&#39;s members.
[**ListOrgs**](OrganizationsAPI.md#ListOrgs) | **Get** /organizations | List the organizations the caller belongs to (with their role).
[**RemoveMember**](OrganizationsAPI.md#RemoveMember) | **Delete** /organizations/{id}/members/{userId} | Remove a member (admin+; owner required to remove owners; never the last owner).
[**UpdateMember**](OrganizationsAPI.md#UpdateMember) | **Patch** /organizations/{id}/members/{userId} | Change a member&#39;s role (admin+; owner required to touch owners).
[**UpdateOrg**](OrganizationsAPI.md#UpdateOrg) | **Patch** /organizations/{id} | Rename an organization (admin or owner).



## AddMember

> AddMemberOutputBody AddMember(ctx, id).AddMemberInputBody(addMemberInputBody).Execute()

Add or invite a member by email (admin+; owner required to grant owner). Unknown emails get an invite link.

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
	id := "id_example" // string | 
	addMemberInputBody := *openapiclient.NewAddMemberInputBody("Email_example") // AddMemberInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.AddMember(context.Background(), id).AddMemberInputBody(addMemberInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.AddMember``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddMember`: AddMemberOutputBody
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.AddMember`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddMemberRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addMemberInputBody** | [**AddMemberInputBody**](AddMemberInputBody.md) |  | 

### Return type

[**AddMemberOutputBody**](AddMemberOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateOrg

> OrgFullView CreateOrg(ctx).CreateOrgInputBody(createOrgInputBody).Execute()

Create an organization (you become its owner) and enqueue tenant provisioning.

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
	createOrgInputBody := *openapiclient.NewCreateOrgInputBody("Name_example") // CreateOrgInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.CreateOrg(context.Background()).CreateOrgInputBody(createOrgInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.CreateOrg``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateOrg`: OrgFullView
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.CreateOrg`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateOrgRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createOrgInputBody** | [**CreateOrgInputBody**](CreateOrgInputBody.md) |  | 

### Return type

[**OrgFullView**](OrgFullView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOrg

> DeleteOrg(ctx, id).Execute()

Delete an organization and enqueue tenant teardown (owner only).

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
	id := "id_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OrganizationsAPI.DeleteOrg(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.DeleteOrg``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOrgRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOrg

> OrgFullView GetOrg(ctx, id).Execute()

Get an organization the caller belongs to.

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
	id := "id_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.GetOrg(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.GetOrg``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOrg`: OrgFullView
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.GetOrg`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOrgRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**OrgFullView**](OrgFullView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListMembers

> MemberListOutputBody ListMembers(ctx, id).Execute()

List an organization's members.

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
	id := "id_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.ListMembers(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.ListMembers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListMembers`: MemberListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.ListMembers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiListMembersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**MemberListOutputBody**](MemberListOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOrgs

> OrgListOutputBody ListOrgs(ctx).Execute()

List the organizations the caller belongs to (with their role).

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
	resp, r, err := apiClient.OrganizationsAPI.ListOrgs(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.ListOrgs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOrgs`: OrgListOutputBody
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.ListOrgs`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListOrgsRequest struct via the builder pattern


### Return type

[**OrgListOutputBody**](OrgListOutputBody.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveMember

> RemoveMember(ctx, id, userId).Execute()

Remove a member (admin+; owner required to remove owners; never the last owner).

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
	id := "id_example" // string | 
	userId := "userId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OrganizationsAPI.RemoveMember(context.Background(), id, userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.RemoveMember``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 
**userId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveMemberRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMember

> MemberJSONView UpdateMember(ctx, id, userId).UpdateMemberInputBody(updateMemberInputBody).Execute()

Change a member's role (admin+; owner required to touch owners).

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
	id := "id_example" // string | 
	userId := "userId_example" // string | 
	updateMemberInputBody := *openapiclient.NewUpdateMemberInputBody("Role_example") // UpdateMemberInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.UpdateMember(context.Background(), id, userId).UpdateMemberInputBody(updateMemberInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.UpdateMember``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMember`: MemberJSONView
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.UpdateMember`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 
**userId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMemberRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateMemberInputBody** | [**UpdateMemberInputBody**](UpdateMemberInputBody.md) |  | 

### Return type

[**MemberJSONView**](MemberJSONView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateOrg

> OrgFullView UpdateOrg(ctx, id).UpdateOrgInputBody(updateOrgInputBody).Execute()

Rename an organization (admin or owner).

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
	id := "id_example" // string | 
	updateOrgInputBody := *openapiclient.NewUpdateOrgInputBody("Name_example") // UpdateOrgInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.UpdateOrg(context.Background(), id).UpdateOrgInputBody(updateOrgInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.UpdateOrg``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateOrg`: OrgFullView
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.UpdateOrg`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateOrgRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateOrgInputBody** | [**UpdateOrgInputBody**](UpdateOrgInputBody.md) |  | 

### Return type

[**OrgFullView**](OrgFullView.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

