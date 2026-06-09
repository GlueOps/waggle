
# AddMemberOutputBody


## Properties

Name | Type
------------ | -------------
`$schema` | string
`invited` | boolean
`member` | [MemberJSONView](MemberJSONView.md)

## Example

```typescript
import type { AddMemberOutputBody } from '@glueops/waggle-sdk'

// TODO: Update the object below with actual values
const example = {
  "$schema": null,
  "invited": null,
  "member": null,
} satisfies AddMemberOutputBody

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as AddMemberOutputBody
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


