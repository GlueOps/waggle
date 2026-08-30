# @glueops/waggle-sdk@0.2.7

A TypeScript SDK client for the localhost API.

## Usage

First, install the SDK from npm.

```bash
npm install @glueops/waggle-sdk --save
```

Next, try it out.


```ts
import {
  Configuration,
  ApiKeysApi,
} from '@glueops/waggle-sdk';
import type { CreateApiKeyRequest } from '@glueops/waggle-sdk';

async function example() {
  console.log("🚀 Testing @glueops/waggle-sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearer
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ApiKeysApi(config);

  const body = {
    // CreateAPIKeyInputBody
    createAPIKeyInputBody: ...,
  } satisfies CreateApiKeyRequest;

  try {
    const data = await api.createApiKey(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```


## Documentation

### API Endpoints

All URIs are relative to */api/v1*

| Class | Method | HTTP request | Description
| ----- | ------ | ------------ | -------------
*ApiKeysApi* | [**createApiKey**](docs/ApiKeysApi.md#createapikey) | **POST** /api-keys | Mint an organization API key for automation (e.g. Terraform). The plaintext token is returned once.
*ApiKeysApi* | [**listApiKeys**](docs/ApiKeysApi.md#listapikeys) | **GET** /api-keys | List the organization\&#39;s API keys (secrets are never returned).
*ApiKeysApi* | [**revokeApiKey**](docs/ApiKeysApi.md#revokeapikey) | **DELETE** /api-keys/{id} | Revoke an organization API key. Idempotent from the caller\&#39;s view.
*AuthApi* | [**authAcceptInvite**](docs/AuthApi.md#authacceptinvite) | **POST** /auth/accept-invite | Accept an organization invite: set a password (if new) and sign in to the org.
*AuthApi* | [**authLogin**](docs/AuthApi.md#authlogin) | **POST** /auth/login | Exchange credentials for an access + refresh token. Returns membership list when no organization_id is given and multiple memberships exist.
*AuthApi* | [**authLogout**](docs/AuthApi.md#authlogout) | **POST** /auth/logout | Revoke the supplied refresh token\&#39;s session. Idempotent.
*AuthApi* | [**authMe**](docs/AuthApi.md#authme) | **GET** /auth/me | Return the authenticated account, its verified-or-pending emails, and its org memberships.
*AuthApi* | [**authRefresh**](docs/AuthApi.md#authrefresh) | **POST** /auth/refresh | Rotate a refresh token, returning a new access + refresh pair.
*AuthApi* | [**authSignup**](docs/AuthApi.md#authsignup) | **POST** /auth/signup | Create an account, organization, and first user; enqueue tenant provisioning.
*AuthApi* | [**authSwitchOrg**](docs/AuthApi.md#authswitchorg) | **POST** /auth/switch | Issue a new token pair scoped to another organization the account belongs to.
*AuthApi* | [**authVerifyEmail**](docs/AuthApi.md#authverifyemail) | **POST** /auth/verify-email | Consume a verification token to mark an email address verified. Idempotent.
*DatacentersApi* | [**createDatacenter**](docs/DatacentersApi.md#createdatacenter) | **POST** /datacenters | Create a datacenter in the caller\&#39;s tenant.
*DatacentersApi* | [**deleteDatacenter**](docs/DatacentersApi.md#deletedatacenter) | **DELETE** /datacenters/{id} | Delete a datacenter.
*DatacentersApi* | [**discoverHypervisors**](docs/DatacentersApi.md#discoverhypervisors) | **POST** /datacenters/{id}/discover | Discover hypervisors from the datacenter\&#39;s Proxmox cluster and upsert them (preserving reserved capacity and schedulable). Set async to run in the background.
*DatacentersApi* | [**getDatacenter**](docs/DatacentersApi.md#getdatacenter) | **GET** /datacenters/{id} | Fetch a datacenter by ID.
*DatacentersApi* | [**listDatacenters**](docs/DatacentersApi.md#listdatacenters) | **GET** /datacenters | List datacenters in the caller\&#39;s tenant.
*DatacentersApi* | [**updateDatacenter**](docs/DatacentersApi.md#updatedatacenter) | **PUT** /datacenters/{id} | Update a datacenter.
*HypervisorsApi* | [**createHypervisor**](docs/HypervisorsApi.md#createhypervisor) | **POST** /hypervisors | Create a hypervisor in the caller\&#39;s tenant.
*HypervisorsApi* | [**deleteHypervisor**](docs/HypervisorsApi.md#deletehypervisor) | **DELETE** /hypervisors/{id} | Delete a hypervisor.
*HypervisorsApi* | [**getHypervisor**](docs/HypervisorsApi.md#gethypervisor) | **GET** /hypervisors/{id} | Fetch a hypervisor by ID.
*HypervisorsApi* | [**listHypervisors**](docs/HypervisorsApi.md#listhypervisors) | **GET** /hypervisors | List hypervisors in the caller\&#39;s tenant.
*HypervisorsApi* | [**updateHypervisor**](docs/HypervisorsApi.md#updatehypervisor) | **PUT** /hypervisors/{id} | Update a hypervisor.
*OrganizationsApi* | [**addMember**](docs/OrganizationsApi.md#addmember) | **POST** /organizations/{id}/members | Add or invite a member by email (admin+; owner required to grant owner). Unknown emails get an invite link.
*OrganizationsApi* | [**createOrg**](docs/OrganizationsApi.md#createorg) | **POST** /organizations | Create an organization (you become its owner) and enqueue tenant provisioning.
*OrganizationsApi* | [**deleteOrg**](docs/OrganizationsApi.md#deleteorg) | **DELETE** /organizations/{id} | Delete an organization and enqueue tenant teardown (owner only).
*OrganizationsApi* | [**getOrg**](docs/OrganizationsApi.md#getorg) | **GET** /organizations/{id} | Get an organization the caller belongs to.
*OrganizationsApi* | [**listMembers**](docs/OrganizationsApi.md#listmembers) | **GET** /organizations/{id}/members | List an organization\&#39;s members.
*OrganizationsApi* | [**listOrgs**](docs/OrganizationsApi.md#listorgs) | **GET** /organizations | List the organizations the caller belongs to (with their role).
*OrganizationsApi* | [**removeMember**](docs/OrganizationsApi.md#removemember) | **DELETE** /organizations/{id}/members/{userId} | Remove a member (admin+; owner required to remove owners; never the last owner).
*OrganizationsApi* | [**updateMember**](docs/OrganizationsApi.md#updatemember) | **PATCH** /organizations/{id}/members/{userId} | Change a member\&#39;s role (admin+; owner required to touch owners).
*OrganizationsApi* | [**updateOrg**](docs/OrganizationsApi.md#updateorg) | **PATCH** /organizations/{id} | Rename an organization (admin or owner).
*PlacementsApi* | [**backfillPlacementVmid**](docs/PlacementsApi.md#backfillplacementvmid) | **PATCH** /placements/{id} | Attach the externally-assigned Proxmox vmid to a placement.
*PlacementsApi* | [**deletePlacement**](docs/PlacementsApi.md#deleteplacement) | **DELETE** /placements/{id} | Remove a placement. The pool\&#39;s desired_count is not adjusted; resize the pool to re-fill the vacancy.
*PlacementsApi* | [**getPlacement**](docs/PlacementsApi.md#getplacement) | **GET** /placements/{id} | Fetch a single placement with its pool, hypervisor, and vmid.
*PlacementsApi* | [**listPlacements**](docs/PlacementsApi.md#listplacements) | **GET** /placements | List all placements in the tenant with pool, slot, and hypervisor context (fleet overview).
*PoolsApi* | [**createPool**](docs/PoolsApi.md#createpool) | **POST** /pools | Create a node pool and place its VMs across hypervisors (anti-affinity spread, all-or-nothing). Placements are available at GET /pools/{id}/placements.
*PoolsApi* | [**deletePool**](docs/PoolsApi.md#deletepool) | **DELETE** /pools/{id} | Delete a pool and release all its placements.
*PoolsApi* | [**getPool**](docs/PoolsApi.md#getpool) | **GET** /pools/{id} | Fetch a pool. Its placements are available at GET /pools/{id}/placements.
*PoolsApi* | [**listPoolPlacements**](docs/PoolsApi.md#listpoolplacements) | **GET** /pools/{id}/placements | List a pool\&#39;s placements (hypervisor + optional vmid).
*PoolsApi* | [**listPools**](docs/PoolsApi.md#listpools) | **GET** /pools | List node pools in the caller\&#39;s tenant.
*PoolsApi* | [**resizePool**](docs/PoolsApi.md#resizepool) | **PATCH** /pools/{id} | Resize a pool\&#39;s desired count. Grow places new VMs (all-or-nothing); shrink removes newest placements (LIFO). Placements are available at GET /pools/{id}/placements.
*SlotsApi* | [**createSlot**](docs/SlotsApi.md#createslot) | **POST** /slots | Create a slot (t-shirt-size VM template) in the caller\&#39;s tenant.
*SlotsApi* | [**deleteSlot**](docs/SlotsApi.md#deleteslot) | **DELETE** /slots/{id} | Delete a slot.
*SlotsApi* | [**getSlot**](docs/SlotsApi.md#getslot) | **GET** /slots/{id} | Fetch a slot by ID.
*SlotsApi* | [**listSlots**](docs/SlotsApi.md#listslots) | **GET** /slots | List slots in the caller\&#39;s tenant.
*SlotsApi* | [**updateSlot**](docs/SlotsApi.md#updateslot) | **PUT** /slots/{id} | Update a slot.
*SystemApi* | [**health**](docs/SystemApi.md#health) | **GET** /health | Health check
*SystemApi* | [**version**](docs/SystemApi.md#version) | **GET** /version | Server version


### Models

- [AcceptInviteInputBody](docs/AcceptInviteInputBody.md)
- [AccountEmailView](docs/AccountEmailView.md)
- [AddMemberInputBody](docs/AddMemberInputBody.md)
- [AddMemberOutputBody](docs/AddMemberOutputBody.md)
- [ApiKeyView](docs/ApiKeyView.md)
- [AuthTokens](docs/AuthTokens.md)
- [BackfillVMIDInputBody](docs/BackfillVMIDInputBody.md)
- [CreateAPIKeyInputBody](docs/CreateAPIKeyInputBody.md)
- [CreateAPIKeyOutputBody](docs/CreateAPIKeyOutputBody.md)
- [CreateOrgInputBody](docs/CreateOrgInputBody.md)
- [CreatePoolInputBody](docs/CreatePoolInputBody.md)
- [DatacenterBody](docs/DatacenterBody.md)
- [DatacenterListOutputBody](docs/DatacenterListOutputBody.md)
- [DatacenterView](docs/DatacenterView.md)
- [DiscoverInputBody](docs/DiscoverInputBody.md)
- [DiscoverOutputBody](docs/DiscoverOutputBody.md)
- [ErrorDetail](docs/ErrorDetail.md)
- [ErrorModel](docs/ErrorModel.md)
- [FleetPlacementListOutputBody](docs/FleetPlacementListOutputBody.md)
- [FleetPlacementView](docs/FleetPlacementView.md)
- [HealthOutputBody](docs/HealthOutputBody.md)
- [HypervisorBody](docs/HypervisorBody.md)
- [HypervisorListOutputBody](docs/HypervisorListOutputBody.md)
- [HypervisorView](docs/HypervisorView.md)
- [ListAPIKeysOutputBody](docs/ListAPIKeysOutputBody.md)
- [LoginInputBody](docs/LoginInputBody.md)
- [LoginOutputBody](docs/LoginOutputBody.md)
- [LogoutInputBody](docs/LogoutInputBody.md)
- [MeOutputBody](docs/MeOutputBody.md)
- [MemberJSONView](docs/MemberJSONView.md)
- [MemberListOutputBody](docs/MemberListOutputBody.md)
- [MembershipView](docs/MembershipView.md)
- [OrgFullView](docs/OrgFullView.md)
- [OrgListOutputBody](docs/OrgListOutputBody.md)
- [OrgView](docs/OrgView.md)
- [PlacementListOutputBody](docs/PlacementListOutputBody.md)
- [PlacementView](docs/PlacementView.md)
- [PoolListOutputBody](docs/PoolListOutputBody.md)
- [PoolView](docs/PoolView.md)
- [RefreshInputBody](docs/RefreshInputBody.md)
- [RefreshOutputBody](docs/RefreshOutputBody.md)
- [ResizePoolInputBody](docs/ResizePoolInputBody.md)
- [SignupInputBody](docs/SignupInputBody.md)
- [SignupOutputBody](docs/SignupOutputBody.md)
- [SlotBody](docs/SlotBody.md)
- [SlotListOutputBody](docs/SlotListOutputBody.md)
- [SlotView](docs/SlotView.md)
- [SwitchOrgInputBody](docs/SwitchOrgInputBody.md)
- [UpdateMemberInputBody](docs/UpdateMemberInputBody.md)
- [UpdateOrgInputBody](docs/UpdateOrgInputBody.md)
- [VerifyEmailInputBody](docs/VerifyEmailInputBody.md)
- [VersionOutputBody](docs/VersionOutputBody.md)

### Authorization


Authentication schemes defined for the API:
<a id="bearer"></a>
#### bearer


- **Type**: HTTP Bearer Token authentication (JWT)

## About

This TypeScript SDK client supports the [Fetch API](https://fetch.spec.whatwg.org/)
and is automatically generated by the
[OpenAPI Generator](https://openapi-generator.tech) project:

- API version: `0.2.7`
- Package version: `0.2.7`
- Generator version: `7.22.0`
- Build package: `org.openapitools.codegen.languages.TypeScriptFetchClientCodegen`

The generated npm module supports the following:

- Environments
  * Node.js
  * Webpack
  * Browserify
- Language levels
  * ES5 - you must have a Promises/A+ library installed
  * ES6
- Module systems
  * CommonJS
  * ES6 module system


## Development

### Building

To build the TypeScript source code, you need to have Node.js and npm installed.
After cloning the repository, navigate to the project directory and run:

```bash
npm install
npm run build
```

### Publishing

Once you've built the package, you can publish it to npm:

```bash
npm publish
```

## License

[]()
