## ADDED Requirements

### Requirement: In-place update of tag_value on tag attachment
The `tencentcloud_tag_attachment` resource SHALL support in-place update of the `tag_value` field. When `tag_value` changes (and `tag_key` and `resource` do not), the provider SHALL call the `UpdateResourceTagValue` API (request fields: `TagKey`, `TagValue`, `Resource`) to modify the associated tag value in a single step, rather than destroying and recreating the resource.

#### Scenario: Update tag_value triggers UpdateResourceTagValue
- **WHEN** a user changes only `tag_value` on an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL call the `UpdateResourceTagValue` API with `TagKey` set to the configured `tag_key`, `TagValue` set to the new `tag_value`, and `Resource` set to the configured `resource`
- **AND** the provider SHALL NOT call `DeleteResourceTag` or `AddResourceTag`

#### Scenario: tag_value is no longer ForceNew
- **WHEN** the `tencentcloud_tag_attachment` resource schema is defined
- **THEN** the `tag_value` field SHALL NOT have `ForceNew: true`
- **AND** an `Update` function SHALL be registered on the resource

#### Scenario: Successful tag_value update rebuilds composite ID
- **WHEN** the `UpdateResourceTagValue` API call succeeds
- **THEN** the provider SHALL rebuild the composite ID as `tagKey + FILED_SP + newTagValue + FILED_SP + resource`
- **AND** the provider SHALL call Read to refresh state

#### Scenario: UpdateResourceTagValue failure is retried
- **WHEN** the `UpdateResourceTagValue` API call returns a retryable error
- **THEN** the provider SHALL retry the call within `tccommon.WriteRetryTimeout`
- **AND** on persistent failure, the provider SHALL return the wrapped error via `tccommon.RetryError`

### Requirement: tag_key and resource remain immutable
The `tencentcloud_tag_attachment` resource SHALL keep `tag_key` and `resource` as `ForceNew: true`. Changes to `tag_key` or `resource` SHALL trigger resource recreation (Create then Delete), consistent with the `UpdateResourceTagValue` API semantics where the tag key and resource cannot change.

#### Scenario: Changing tag_key recreates the resource
- **WHEN** a user changes `tag_key` on an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL destroy the existing resource and create a new one

#### Scenario: Changing resource recreates the resource
- **WHEN** a user changes `resource` on an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL destroy the existing resource and create a new one

### Requirement: Update detects tag_value change and reads the new value
The `Update` function SHALL detect that `tag_value` has changed (via `d.HasChange("tag_value")`) and read the new `tag_value` from the configured value (via `d.GetOk("tag_value")`), and use the new value as the `TagValue` parameter of the `UpdateResourceTagValue` API request. The `UpdateResourceTagValue` API only accepts the new `TagValue` (it does not need the old value), so the old `tag_value` need not be read.

#### Scenario: tag_value change is detected and new value is used
- **WHEN** the Update function is invoked for a `tag_value` change
- **THEN** the provider SHALL detect the change using `d.HasChange("tag_value")`
- **AND** the provider SHALL read the new `tag_value` from config using `d.GetOk("tag_value")`
- **AND** the provider SHALL pass the new `tag_value` as `TagValue` in the `UpdateResourceTagValue` request
