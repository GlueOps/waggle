
# PoolResultBody


## Properties

Name | Type
------------ | -------------
`$schema` | string
`placements` | [Array&lt;PlacementView&gt;](PlacementView.md)
`pool` | [PoolView](PoolView.md)

## Example

```typescript
import type { PoolResultBody } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "placements": null,
  "pool": null,
} satisfies PoolResultBody

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as PoolResultBody
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


