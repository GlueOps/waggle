
# HypervisorBody


## Properties

Name | Type
------------ | -------------
`$schema` | string
`cpuOvercommitRatio` | number
`cpuReserved` | number
`cpuTotal` | number
`datacenterId` | string
`diskGbReserved` | number
`diskGbTotal` | number
`name` | string
`ramGbReserved` | number
`ramGbTotal` | number
`schedulable` | boolean

## Example

```typescript
import type { HypervisorBody } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "cpuOvercommitRatio": null,
  "cpuReserved": null,
  "cpuTotal": null,
  "datacenterId": null,
  "diskGbReserved": null,
  "diskGbTotal": null,
  "name": null,
  "ramGbReserved": null,
  "ramGbTotal": null,
  "schedulable": null,
} satisfies HypervisorBody

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as HypervisorBody
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


