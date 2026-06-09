
# MembershipView


## Properties

Name | Type
------------ | -------------
`organizationId` | string
`organizationName` | string
`organizationSlug` | string

## Example

```typescript
import type { MembershipView } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "organizationId": null,
  "organizationName": null,
  "organizationSlug": null,
} satisfies MembershipView

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as MembershipView
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


