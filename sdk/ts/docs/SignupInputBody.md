
# SignupInputBody


## Properties

Name | Type
------------ | -------------
`$schema` | string
`displayName` | string
`email` | string
`organizationName` | string
`password` | string

## Example

```typescript
import type { SignupInputBody } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "displayName": null,
  "email": null,
  "organizationName": null,
  "password": null,
} satisfies SignupInputBody

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SignupInputBody
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


