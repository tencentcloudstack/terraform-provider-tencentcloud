# tag-attachment-tag-value-update Specification

## Purpose
TBD - created by syncing change modify-tag-attachment-tag-value-update. Update Purpose after sync.
## Requirements
### Requirement: tag_value is updatable in-place

The `tencentcloud_tag_attachment` resource SHALL NOT mark `tag_value` as `ForceNew`. When a user changes `tag_value` in the configuration, the provider SHALL update the tag value in-place via the `UpdateResourceTagValue` cloud API rather than destroying and recreating the attachment.

#### Scenario: Update tag_value triggers UpdateResourceTagValue

- **WHEN** a user changes `tag_value` on an existing `tencentcloud_tag_attachment` resource (keeping `tag_key` and `resource` unchanged)
- **THEN** the provider SHALL call `UpdateResourceTagValue` with the `TagKey`, the new `TagValue`, and the `Resource` six-segment description to atomically change the tag value on the resource

#### Scenario: tag_key change remains ForceNew

- **WHEN** a user changes `tag_key` on an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL destroy and recreate the attachment (ForceNew behavior is retained for `tag_key`)

#### Scenario: resource change remains ForceNew

- **WHEN** a user changes `resource` on an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL destroy and recreate the attachment (ForceNew behavior is retained for `resource`)

### Requirement: Update function calls UpdateResourceTagValue

The `tencentcloud_tag_attachment` resource SHALL implement an `Update` function that, when `tag_value` has changed, calls the `UpdateResourceTagValue` API (tag v20180813) with `TagKey`, `TagValue`, and `Resource` parameters.

#### Scenario: Successful tag value update

- **WHEN** the Update function detects `tag_value` has changed and calls `UpdateResourceTagValue`
- **AND** the API call succeeds
- **THEN** the provider SHALL re-set the composite ID to `tagKey + FILED_SP + newTagValue + FILED_SP + resource` so it reflects the new tag value
- **AND** the provider SHALL call Read to refresh state from the cloud

#### Scenario: UpdateResourceTagValue API error

- **WHEN** the `UpdateResourceTagValue` API call fails
- **THEN** the provider SHALL return the wrapped retry error and SHALL NOT re-set the ID or call Read

#### Scenario: tag_value unchanged

- **WHEN** the Update function is invoked but `tag_value` has not changed (only other non-ForceNew fields changed, if any in future)
- **THEN** the provider SHALL NOT call `UpdateResourceTagValue`

### Requirement: Service-layer wrapper for UpdateResourceTagValue

The provider SHALL include a service-layer function `UpdateTagAttachmentTagValue` in `service_tencentcloud_tag.go` that wraps the `UpdateResourceTagValue` API call with retry handling (`tccommon.WriteRetryTimeout`), `ratelimit.Check`, and `tccommon.RetryError` error wrapping, consistent with existing service functions.

#### Scenario: Service function retry behavior

- **WHEN** `UpdateTagAttachmentTagValue` is called and the API returns a retryable error
- **THEN** the function SHALL retry within `tccommon.WriteRetryTimeout` and return the wrapped error on exhaustion

### Requirement: Composite ID re-set after update

After a successful `UpdateResourceTagValue` call, the Update function SHALL re-set the resource ID using the new `tag_value` so the composite ID (`tagKey + FILED_SP + tagValue + FILED_SP + resource`) matches the current tag binding.

#### Scenario: ID reflects new tag value

- **WHEN** `tag_value` is updated from `oldValue` to `newValue`
- **THEN** after a successful update the resource ID SHALL be `tagKey + FILED_SP + newValue + FILED_SP + resource`
