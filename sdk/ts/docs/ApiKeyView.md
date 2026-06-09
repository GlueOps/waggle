
# ApiKeyView


## Properties

Name | Type
------------ | -------------
`createdAt` | Date
`expiresAt` | Date
`id` | string
`lastUsedAt` | Date
`name` | string
`prefix` | string
`revokedAt` | Date

## Example

```typescript
import type { ApiKeyView } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "createdAt": null,
  "expiresAt": null,
  "id": null,
  "lastUsedAt": null,
  "name": null,
  "prefix": null,
  "revokedAt": null,
} satisfies ApiKeyView

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiKeyView
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


