## ADDED Requirements

### Requirement: Support custom_image parameter in native node pool

The `tencentcloud_kubernetes_native_node_pool` resource SHALL support specifying a custom image ID via the `custom_image` field inside the `native` block. The field SHALL be mapped to `Native.CustomImage` in the `CreateNodePool` and `ModifyNodePool` API requests, and SHALL be populated from `Native.CustomImage` in the Read response.

#### Scenario: Create node pool with custom_image

- **WHEN** user specifies `custom_image` inside the `native` block of a `tencentcloud_kubernetes_native_node_pool` resource
- **THEN** the provider SHALL pass `Native.CustomImage` to the `CreateNodePool` API request
- **AND** the node pool SHALL be created using the specified custom image

#### Scenario: Update node pool custom_image

- **WHEN** user changes the `custom_image` value inside the `native` block of an existing `tencentcloud_kubernetes_native_node_pool` resource
- **THEN** the provider SHALL pass the new `Native.CustomImage` to the `ModifyNodePool` API request
- **AND** the node pool SHALL be updated in place without recreation

#### Scenario: Read custom_image from API

- **WHEN** the provider reads a `tencentcloud_kubernetes_native_node_pool` resource whose API response includes `Native.CustomImage`
- **THEN** the provider SHALL populate `custom_image` in the `native` block of Terraform state
- **AND** when `Native.CustomImage` is nil, the provider SHALL NOT overwrite the state value

### Requirement: custom_image schema definition

The `custom_image` field SHALL be defined inside the `native` block schema of the `tencentcloud_kubernetes_native_node_pool` resource.

#### Scenario: custom_image field specification

- **WHEN** defining the `custom_image` schema field
- **THEN** it SHALL be of type String
- **AND** it SHALL be Optional
- **AND** it SHALL NOT have ForceNew set (it is mutable via `ModifyNodePool`)
- **AND** it SHALL have a description explaining it specifies the custom image ID

### Requirement: Backward compatibility

The addition of the `custom_image` parameter SHALL NOT break existing `tencentcloud_kubernetes_native_node_pool` configurations.

#### Scenario: Existing configurations without custom_image

- **WHEN** user applies an existing `tencentcloud_kubernetes_native_node_pool` configuration that does not specify `custom_image`
- **THEN** the node pool SHALL be created/updated successfully without passing `Native.CustomImage`
- **AND** no changes SHALL be required to existing Terraform configurations

#### Scenario: State file compatibility

- **WHEN** the provider reads state from a node pool created before this feature was added
- **THEN** the provider SHALL handle the missing `CustomImage` field gracefully
- **AND** no spurious diffs SHALL be generated
