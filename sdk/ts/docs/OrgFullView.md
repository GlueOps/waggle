
# OrgFullView


## Properties

Name | Type
------------ | -------------
`$schema` | string
`createdAt` | Date
`domain` | string
`id` | string
`name` | string
`role` | string
`slug` | string
`status` | string

## Example

```typescript
import type { OrgFullView } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "createdAt": null,
  "domain": null,
  "id": null,
  "name": null,
  "role": null,
  "slug": null,
  "status": null,
} satisfies OrgFullView

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as OrgFullView
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


