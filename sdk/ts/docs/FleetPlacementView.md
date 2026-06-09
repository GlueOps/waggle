
# FleetPlacementView


## Properties

Name | Type
------------ | -------------
`createdAt` | Date
`diskGb` | number
`hypervisorId` | string
`hypervisorName` | string
`id` | string
`poolId` | string
`poolName` | string
`ramGb` | number
`slotName` | string
`vcpu` | number
`vmid` | number

## Example

```typescript
import type { FleetPlacementView } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "createdAt": null,
  "diskGb": null,
  "hypervisorId": null,
  "hypervisorName": null,
  "id": null,
  "poolId": null,
  "poolName": null,
  "ramGb": null,
  "slotName": null,
  "vcpu": null,
  "vmid": null,
} satisfies FleetPlacementView

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as FleetPlacementView
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


