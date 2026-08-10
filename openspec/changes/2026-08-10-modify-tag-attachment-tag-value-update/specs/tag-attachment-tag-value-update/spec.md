## ADDED Requirements

### Requirement: tag_value is updatable in-place
The `tencentcloud_tag_attachment` resource SHALL NOT mark `tag_value` as `ForceNew`. Changes to `tag_value` in the Terraform configuration SHALL trigger an in-place `Update` operation rather than destroying and recreating the attachment.

#### Scenario: Update tag_value in-place
- **WHEN** a user changes `tag_value` on an existing `tencentcloud_tag_attachment` resource (keeping `tag_key` and `resource` unchanged)
- **THEN** the provider SHALL perform an in-place update via the `ModifyResourceTags` cloud API instead of destroying and recreating the resource

#### Scenario: tag_value no longer forces replacement
- **WHEN** the `tag_value` field is changed in the Terraform configuration
- **THEN** the provider SHALL NOT execute `DeleteResourceTag` followed by `AddResourceTag`
- **AND** the provider SHALL NOT require the attachment to be destroyed and recreated

### Requirement: In-place update uses ModifyResourceTags atomically
The `tencentcloud_tag_attachment` `Update` function SHALL update the tag value by calling the `ModifyResourceTags` cloud API in a single atomic request, placing the new `{tag_key: new_tag_value}` in `ReplaceTags` (and leaving `DeleteTags` empty), so the resource is never left without the tag during the update. Per the `ModifyResourceTags` API contract, when a resource already has the `tag_key` associated, placing `{tag_key, new_value}` in `ReplaceTags` changes the existing value to the new value in place.

#### Scenario: Atomic update of tag value
- **WHEN** the `Update` function is invoked because `tag_value` changed from `old_value` to `new_value`
- **THEN** the provider SHALL call `ModifyResourceTags` with `Resource` set to the attachment's `resource`, `ReplaceTags` containing `{TagKey: tag_key, TagValue: new_value}`, and `DeleteTags` empty
- **AND** the replace operation SHALL be sent in a single API request

#### Scenario: New tag value placed in ReplaceTags
- **WHEN** updating `tag_value` from `old_value` to `new_value`
- **THEN** the `{tag_key, new_value}` association SHALL be passed via the `ReplaceTags` field of the `ModifyResourceTags` request
- **AND** the `DeleteTags` field SHALL be empty, because `ReplaceTags` and `DeleteTags` must not contain the same tag key per the API contract

### Requirement: tag_key and resource remain ForceNew
The `tag_key` and `resource` fields of `tencentcloud_tag_attachment` SHALL remain `ForceNew: true`, because they identify which resource and tag key the attachment refers to and cannot be changed without recreating the attachment.

#### Scenario: Changing tag_key forces replacement
- **WHEN** a user changes `tag_key` on an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL force replacement (destroy + create) of the attachment

#### Scenario: Changing resource forces replacement
- **WHEN** a user changes `resource` on an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL force replacement (destroy + create) of the attachment

### Requirement: Composite id refreshed after in-place update
After a successful in-place update of `tag_value`, the provider SHALL refresh the composite id to `tag_key + FILED_SP + new_tag_value + FILED_SP + resource` so that subsequent Read operations resolve the correct attachment.

#### Scenario: Id reflects new tag value after update
- **WHEN** the `Update` function successfully updates `tag_value` from `old_value` to `new_value`
- **THEN** the provider SHALL set `d.SetId(tagKey + FILED_SP + new_tag_value + FILED_SP + resource)`
- **AND** the Read function SHALL be invoked to refresh state using the updated id

### Requirement: Update reuses TagService.ModifyTags
The `Update` function SHALL reuse the existing `TagService.ModifyTags` service-layer method (which wraps `ModifyResourceTags` with retry handling) to perform the update, rather than introducing a new service-layer method or calling the SDK directly.

#### Scenario: Update invokes ModifyTags service method
- **WHEN** the `Update` function performs the in-place tag value update
- **THEN** the provider SHALL call `TagService.ModifyTags(ctx, resource, replaceTags, deleteKeys)` with `replaceTags = {tag_key: new_value}` and `deleteKeys = nil`
- **AND** no new service-layer method SHALL be added for this operation
