
# DatacenterListOutputBody


## Properties

Name | Type
------------ | -------------
`$schema` | string
`items` | [Array&lt;DatacenterView&gt;](DatacenterView.md)

## Example

```typescript
import type { DatacenterListOutputBody } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "items": null,
} satisfies DatacenterListOutputBody

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DatacenterListOutputBody
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


