
# PoolView


## Properties

Name | Type
------------ | -------------
`$schema` | string
`createdAt` | Date
`datacenterId` | string
`desiredCount` | number
`id` | string
`metadata` | any
`name` | string
`slotId` | string
`updatedAt` | Date

## Example

```typescript
import type { PoolView } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "createdAt": null,
  "datacenterId": null,
  "desiredCount": null,
  "id": null,
  "metadata": null,
  "name": null,
  "slotId": null,
  "updatedAt": null,
} satisfies PoolView

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as PoolView
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


