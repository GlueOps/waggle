
# HypervisorView


## Properties

Name | Type
------------ | -------------
`$schema` | string
`cpuBookable` | number
`cpuEffectiveTotal` | number
`cpuOvercommitRatio` | number
`cpuReserved` | number
`cpuTotal` | number
`cpuUsed` | number
`createdAt` | Date
`datacenterId` | string
`diskGbBookable` | number
`diskGbReserved` | number
`diskGbTotal` | number
`diskGbUsed` | number
`id` | string
`lastSyncedAt` | Date
`name` | string
`ramGbBookable` | number
`ramGbReserved` | number
`ramGbTotal` | number
`ramGbUsed` | number
`schedulable` | boolean
`updatedAt` | Date

## Example

```typescript
import type { HypervisorView } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "cpuBookable": null,
  "cpuEffectiveTotal": null,
  "cpuOvercommitRatio": null,
  "cpuReserved": null,
  "cpuTotal": null,
  "cpuUsed": null,
  "createdAt": null,
  "datacenterId": null,
  "diskGbBookable": null,
  "diskGbReserved": null,
  "diskGbTotal": null,
  "diskGbUsed": null,
  "id": null,
  "lastSyncedAt": null,
  "name": null,
  "ramGbBookable": null,
  "ramGbReserved": null,
  "ramGbTotal": null,
  "ramGbUsed": null,
  "schedulable": null,
  "updatedAt": null,
} satisfies HypervisorView

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as HypervisorView
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


