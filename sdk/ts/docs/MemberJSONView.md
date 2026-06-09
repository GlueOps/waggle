
# MemberJSONView


## Properties

Name | Type
------------ | -------------
`$schema` | string
`accountId` | string
`createdAt` | Date
`displayName` | string
`email` | string
`isActive` | boolean
`lastLoginAt` | Date
`pending` | boolean
`role` | string
`userId` | string

## Example

```typescript
import type { MemberJSONView } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "accountId": null,
  "createdAt": null,
  "displayName": null,
  "email": null,
  "isActive": null,
  "lastLoginAt": null,
  "pending": null,
  "role": null,
  "userId": null,
} satisfies MemberJSONView

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as MemberJSONView
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


