
# SignupOutputBody


## Properties

Name | Type
------------ | -------------
`$schema` | string
`accessExpiresAt` | Date
`accessToken` | string
`accountId` | string
`organization` | [OrgView](OrgView.md)
`refreshExpiresAt` | Date
`refreshToken` | string

## Example

```typescript
import type { SignupOutputBody } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "accessExpiresAt": null,
  "accessToken": null,
  "accountId": null,
  "organization": null,
  "refreshExpiresAt": null,
  "refreshToken": null,
} satisfies SignupOutputBody

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SignupOutputBody
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


