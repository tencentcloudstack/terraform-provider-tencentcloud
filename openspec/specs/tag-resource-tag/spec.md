# tag-resource-tag Specification

## Purpose
TBD - created by archiving change add-tag-resource-tag. Update Purpose after archive.
## Requirements
### Requirement: Resource manages a single tag binding on a single cloud resource

The `tencentcloud_tag_resource_tag` resource SHALL manage the full lifecycle (create, read, update, delete) of a single tag key/value binding to a single cloud resource described by its six-segment resource coordinate. The resource SHALL expose `tag_key`, `tag_value`, and `resource` as user-facing schema arguments.

#### Scenario: Create binds a tag to a resource

- **WHEN** a user applies a `tencentcloud_tag_resource_tag` configuration with `tag_key`, `tag_value`, and `resource` set
- **THEN** the provider SHALL call `AddResourceTag` with the supplied `TagKey`, `TagValue`, and `Resource`
- **AND** the provider SHALL set the Terraform ID to `tag_key` + `resource` joined by the field separator
- **AND** if the API response is nil or empty, the provider SHALL return a `NonRetryableError` without writing an empty ID

#### Scenario: Update changes only the tag value

- **WHEN** a user changes only `tag_value` on an existing `tencentcloud_tag_resource_tag`
- **THEN** the provider SHALL call `UpdateResourceTagValue` with the stored `TagKey`, the new `TagValue`, and the stored `Resource`
- **AND** the provider SHALL NOT recreate the resource

#### Scenario: Changing tag_key or resource recreates the binding

- **WHEN** a user changes `tag_key` or `resource` on an existing `tencentcloud_tag_resource_tag`
- **THEN** the provider SHALL treat the change as a ForceNew (destroy + create), because the cloud API cannot rename the key or move a binding in place

#### Scenario: Delete unbinds the tag from the resource

- **WHEN** a user destroys a `tencentcloud_tag_resource_tag`
- **THEN** the provider SHALL call `DeleteResourceTag` with the stored `TagKey` and `Resource`
- **AND** after a successful delete the resource SHALL be removed from state

### Requirement: Read locates the binding by tag key

The Read operation SHALL call `GetResources` with the raw six-segment `resource` in `ResourceList`, iterate `ResourceTagMappingList` to find the mapping whose `Resource` matches, and SHALL select the tag whose `TagKey` equals the stored `tag_key`. The matched tag's `TagValue` SHALL be written into state.

#### Scenario: Binding exists

- **WHEN** the Read operation queries `GetResources` for the resource and a tag with `TagKey` equal to the stored `tag_key` is present in the matched `ResourceTagMapping`
- **THEN** the provider SHALL set `tag_value` from the matched tag's `TagValue`
- **AND** the provider SHALL set `resource` from the stored value
- **AND** the provider SHALL set `tag_key` from the stored value

#### Scenario: Binding not found

- **WHEN** the Read operation queries `GetResources` and no tag matches the stored `tag_key` (the binding was deleted out-of-band)
- **THEN** the provider SHALL log the current `logId` and `d.Id()` to preserve the scene
- **AND** the provider SHALL set `d.SetId("")` to remove the resource from state

### Requirement: Nil checks before Set

The Read operation SHALL check whether each response field is nil before calling `d.Set()` for that field. A nil field SHALL NOT cause a panic.

#### Scenario: Nil tag value in response

- **WHEN** the matched tag has a nil `TagValue`
- **THEN** the provider SHALL skip calling `d.Set("tag_value", ...)` for that field

### Requirement: Retry and error handling

The Create, Update, and Delete operations SHALL wrap cloud API calls in `tccommon` retry timeouts. On API failure the provider SHALL return errors wrapped with `tccommon.RetryError()`. The retry block SHALL contain only the API call; setting the ID and other post-success actions SHALL occur outside the retry block. The Read operation SHALL call the service-layer `GetResources` lookup directly (no retry wrapper), with `ratelimit.Check` before the API call.

#### Scenario: Transient API failure retried

- **WHEN** a cloud API call fails with a retryable error during Create/Update/Delete
- **THEN** the provider SHALL retry within the configured timeout
- **AND** if retries are exhausted, the provider SHALL return the wrapped error

#### Scenario: Non-retryable error on empty create response

- **WHEN** `AddResourceTag` returns a nil or empty response
- **THEN** the provider SHALL return a `NonRetryableError` and SHALL NOT set the ID

### Requirement: Provider registration and documentation

The provider SHALL register `tencentcloud_tag_resource_tag` in `tencentcloud/provider.go` and SHALL provide a markdown doc following the provider's doc conventions.

#### Scenario: Resource is available in the provider

- **WHEN** the provider is built
- **THEN** `tencentcloud_tag_resource_tag` SHALL be a registered resource type
- **AND** the resource SHALL importable via `terraform import` using the composite ID `tag_key` + `resource`

