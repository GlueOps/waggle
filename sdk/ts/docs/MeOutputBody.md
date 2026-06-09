
# MeOutputBody


## Properties

Name | Type
------------ | -------------
`$schema` | string
`accountId` | string
`currentOrganization` | [OrgView](OrgView.md)
`displayName` | string
`emails` | [Array&lt;AccountEmailView&gt;](AccountEmailView.md)
`lastLoginAt` | Date
`memberships` | [Array&lt;MembershipView&gt;](MembershipView.md)

## Example

```typescript
import type { MeOutputBody } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "accountId": null,
  "currentOrganization": null,
  "displayName": null,
  "emails": null,
  "lastLoginAt": null,
  "memberships": null,
} satisfies MeOutputBody

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as MeOutputBody
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


