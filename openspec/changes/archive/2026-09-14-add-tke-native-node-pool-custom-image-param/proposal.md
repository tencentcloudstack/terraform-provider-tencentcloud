## Why

The `tencentcloud_kubernetes_native_node_pool` resource currently does not support the `CustomImage` parameter. TKE native node pools allow users to specify a custom image ID when creating or updating a node pool, but this field is absent from the Terraform schema. Users who want to provision native node pool nodes based on a custom image cannot do so through Terraform today.

## What Changes

- Add a new optional schema field `custom_image` (String) to the `native` block of the `tencentcloud_kubernetes_native_node_pool` resource.
- Map this field to `Native.CustomImage` in the `CreateNodePool` API request (`CreateNativeNodePoolParam.CustomImage`).
- Map this field to `Native.CustomImage` in the `ModifyNodePool` API request (`UpdateNativeNodePoolParam.CustomImage`).
- Populate `custom_image` from `NativeNodePoolInfo.CustomImage` in the Read handler so Terraform state stays in sync with the cloud.
- The field is mutable (not ForceNew) since `ModifyNodePool` supports updating it.

## Capabilities

### New Capabilities

None - this is an enhancement to an existing resource.

### Modified Capabilities

- `tke-native-node-pool-resource`: Add support for the `custom_image` parameter in the `native` block of the `tencentcloud_kubernetes_native_node_pool` resource

## Impact

**Affected Code:**
- `tencentcloud/services/tke/resource_tc_kubernetes_native_node_pool.go`: Schema definition (native block), Create function, Read function, Update function
- `tencentcloud/services/tke/resource_tc_kubernetes_native_node_pool.md`: Documentation with the new parameter

**User Impact:**
- **Non-breaking**: Existing configurations continue to work unchanged; the new field is optional.
- Users can now specify a custom image ID when creating or updating native node pools.

**Dependencies:**
- Requires no SDK changes. The `CustomImage` field already exists in `CreateNativeNodePoolParam`, `UpdateNativeNodePoolParam`, and `NativeNodePoolInfo` structs in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20220501`.
