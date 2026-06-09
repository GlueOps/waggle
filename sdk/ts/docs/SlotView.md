
# SlotView


## Properties

Name | Type
------------ | -------------
`$schema` | string
`createdAt` | Date
`diskGb` | number
`id` | string
`name` | string
`ramGb` | number
`updatedAt` | Date
`vcpu` | number

## Example

```typescript
import type { SlotView } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "createdAt": null,
  "diskGb": null,
  "id": null,
  "name": null,
  "ramGb": null,
  "updatedAt": null,
  "vcpu": null,
} satisfies SlotView

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SlotView
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


