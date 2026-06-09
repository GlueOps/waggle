
# ListAPIKeysOutputBody


## Properties

Name | Type
------------ | -------------
`$schema` | string
`items` | [Array&lt;ApiKeyView&gt;](ApiKeyView.md)

## Example

```typescript
import type { ListAPIKeysOutputBody } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "items": null,
} satisfies ListAPIKeysOutputBody

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ListAPIKeysOutputBody
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


