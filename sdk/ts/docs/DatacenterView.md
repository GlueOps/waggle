
# DatacenterView


## Properties

Name | Type
------------ | -------------
`$schema` | string
`cpuOvercommitRatio` | number
`createdAt` | Date
`hasToken` | boolean
`id` | string
`insecureSkipVerify` | boolean
`name` | string
`updatedAt` | Date
`url` | string

## Example

```typescript
import type { DatacenterView } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "cpuOvercommitRatio": null,
  "createdAt": null,
  "hasToken": null,
  "id": null,
  "insecureSkipVerify": null,
  "name": null,
  "updatedAt": null,
  "url": null,
} satisfies DatacenterView

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DatacenterView
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


