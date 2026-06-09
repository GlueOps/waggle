
# PlacementView


## Properties

Name | Type
------------ | -------------
`$schema` | string
`createdAt` | Date
`hypervisorId` | string
`hypervisorName` | string
`id` | string
`vmid` | number

## Example

```typescript
import type { PlacementView } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "createdAt": null,
  "hypervisorId": null,
  "hypervisorName": null,
  "id": null,
  "vmid": null,
} satisfies PlacementView

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as PlacementView
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


