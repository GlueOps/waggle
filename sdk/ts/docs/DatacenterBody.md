
# DatacenterBody


## Properties

Name | Type
------------ | -------------
`$schema` | string
`cpuOvercommitRatio` | number
`insecureSkipVerify` | boolean
`name` | string
`token` | string
`url` | string

## Example

```typescript
import type { DatacenterBody } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "cpuOvercommitRatio": null,
  "insecureSkipVerify": null,
  "name": null,
  "token": null,
  "url": null,
} satisfies DatacenterBody

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DatacenterBody
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


