
# CreatePoolInputBody


## Properties

Name | Type
------------ | -------------
`$schema` | string
`datacenterId` | string
`desiredCount` | number
`metadata` | any
`name` | string
`slotId` | string

## Example

```typescript
import type { CreatePoolInputBody } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "datacenterId": null,
  "desiredCount": null,
  "metadata": null,
  "name": null,
  "slotId": null,
} satisfies CreatePoolInputBody

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreatePoolInputBody
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


