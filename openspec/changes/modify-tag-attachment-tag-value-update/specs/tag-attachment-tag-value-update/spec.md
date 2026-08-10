## ADDED Requirements

### Requirement: tag_value is updatable in place
The `tencentcloud_tag_attachment` resource SHALL allow the `tag_value` field to be updated in place without destroying and recreating the resource. The `tag_value` schema field SHALL NOT have `ForceNew: true`. When `tag_value` changes, the provider SHALL call the Tag cloud API `UpdateResourceTagValue` to modify the tag value associated with the resource while keeping the tag key unchanged.

#### Scenario: Update tag_value in place
- **WHEN** a user changes only the `tag_value` field on an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL call `UpdateResourceTagValue` with the `TagKey`, the new `TagValue`, and the `Resource` from the resource state
- **AND** the provider SHALL NOT delete the existing attachment nor create a new one
- **AND** after a successful update, the composite id SHALL be refreshed to `tagKey#newTagValue#resource`

#### Scenario: Update with no change to tag_value
- **WHEN** Terraform invokes the Update function but `tag_value` has not changed
- **THEN** the provider SHALL NOT call `UpdateResourceTagValue`
- **AND** the provider SHALL proceed to read the current state

### Requirement: tag_key and resource remain immutable
The `tencentcloud_tag_attachment` resource SHALL keep `tag_key` and `resource` as `ForceNew: true`. Changes to `tag_key` or `resource` SHALL trigger destroy-and-recreate behavior, because the `UpdateResourceTagValue` cloud API only modifies the tag value and does not support changing the tag key or the resource of an existing attachment.

#### Scenario: Change tag_key triggers recreate
- **WHEN** a user changes the `tag_key` field on an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL destroy the existing attachment and create a new one (ForceNew behavior)

#### Scenario: Change resource triggers recreate
- **WHEN** a user changes the `resource` field on an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL destroy the existing attachment and create a new one (ForceNew behavior)

### Requirement: Update API call uses retry with write timeout
The Update function SHALL wrap the `UpdateResourceTagValue` cloud API call in a `resource.Retry` using `tccommon.WriteRetryTimeout`. If the API call fails, the error SHALL be wrapped with `tccommon.RetryError`. The retry block SHALL only contain the API call; setting the composite id and other success-path operations SHALL happen outside the retry block, after the error-handling path.

#### Scenario: API call failure retries
- **WHEN** the `UpdateResourceTagValue` call fails with a retryable error
- **THEN** the provider SHALL retry the call up to the `WriteRetryTimeout` limit
- **AND** the composite id SHALL NOT be updated until the call succeeds

#### Scenario: API call success updates id then reads
- **WHEN** the `UpdateResourceTagValue` call succeeds
- **THEN** the provider SHALL set the composite id to `tagKey#newTagValue#resource`
- **AND** the provider SHALL call the Read function to refresh state
