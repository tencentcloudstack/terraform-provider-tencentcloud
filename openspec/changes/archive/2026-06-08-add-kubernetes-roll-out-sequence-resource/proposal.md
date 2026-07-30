## Why

TKE (Tencent Kubernetes Engine) supports cluster roll-out sequences that allow users to define staged deployment workflows with configurable soak times between steps. Currently, there is no Terraform resource to manage these roll-out sequences, requiring users to manage them manually through the console or API. Adding a Terraform resource enables infrastructure-as-code management of roll-out sequences.

## What Changes

- Add a new Terraform resource `tencentcloud_kubernetes_roll_out_sequence` (RESOURCE_KIND_GENERAL) that supports full CRUD lifecycle management of TKE cluster roll-out sequences.
- The resource will support:
  - Creating a roll-out sequence with a name, sequence flows (steps with tags and soak times), and enabled status
  - Reading/describing roll-out sequences via the list API filtered by ID
  - Updating the name, sequence flows, and enabled status
  - Deleting a roll-out sequence by ID
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`

## Capabilities

### New Capabilities
- `kubernetes-roll-out-sequence-resource`: CRUD resource for managing TKE cluster roll-out sequences, including sequence flows with tags and soak times

### Modified Capabilities

## Impact

- New files:
  - `tencentcloud/services/tke/resource_tc_kubernetes_roll_out_sequence.go` - Resource implementation
  - `tencentcloud/services/tke/resource_tc_kubernetes_roll_out_sequence_test.go` - Unit tests with gomonkey mocks
  - `tencentcloud/services/tke/resource_tc_kubernetes_roll_out_sequence.md` - Example usage documentation
- Modified files:
  - `tencentcloud/provider.go` - Register the new resource
  - `tencentcloud/provider.md` - Add resource entry
- Dependencies:
  - Uses existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525` package (already vendored)
  - Cloud APIs: `CreateRollOutSequence`, `DescribeRollOutSequences`, `ModifyRollOutSequence`, `DeleteRollOutSequence`
