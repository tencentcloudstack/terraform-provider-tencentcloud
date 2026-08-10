## ADDED Requirements

### Requirement: tag_value field is updatable
The `tencentcloud_tag_attachment` resource SHALL allow the `tag_value` field to be modified in place. The `tag_value` field SHALL NOT have `ForceNew: true`. When a user changes `tag_value` in the Terraform configuration, the provider SHALL call the `UpdateResourceTagValue` API to update the tag value on the bound resource instead of destroying and recreating the tag binding.

#### Scenario: Modify tag_value triggers in-place update
- **WHEN** a user changes the `tag_value` field of an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL NOT destroy and recreate the resource
- **AND** the provider SHALL call the `UpdateResourceTagValue` API with `TagKey`, `TagValue` (the new value), and `Resource` to update the tag value in place

#### Scenario: tag_value field schema
- **GIVEN** the schema definition of `tencentcloud_tag_attachment`
- **WHEN** inspecting the `tag_value` field
- **THEN** the field SHALL be `Required`, `TypeString`
- **AND** the field SHALL NOT have `ForceNew: true`

#### Scenario: tag_key and resource remain immutable
- **GIVEN** the schema definition of `tencentcloud_tag_attachment`
- **WHEN** the user changes `tag_key` or `resource`
- **THEN** the provider SHALL trigger resource replacement (ForceNew) for those fields
- **AND** the `UpdateResourceTagValue` API SHALL only be invoked when `tag_value` changes

### Requirement: Update function for tag_value changes
The `tencentcloud_tag_attachment` resource SHALL implement an `Update` function that detects changes to the `tag_value` field and calls the `UpdateResourceTagValue` API.

#### Scenario: Update detects tag_value change
- **WHEN** the `Update` function is invoked and `d.HasChange("tag_value")` is true
- **THEN** the provider SHALL extract `tag_key`, the new `tag_value`, and `resource` from the schema data
- **AND** the provider SHALL call the `ModifyTagTagAttachmentValue` service method with `tagKey`, `newTagValue`, and `resource`
- **AND** after a successful API call, the provider SHALL update the composite id to `tagKey#newTagValue#resource`

#### Scenario: Update with no tag_value change is a no-op
- **WHEN** the `Update` function is invoked and `tag_value` has not changed
- **THEN** the provider SHALL NOT call the `UpdateResourceTagValue` API
- **AND** the provider SHALL return the result of the `Read` function

### Requirement: Composite id update after tag_value change
After a successful `tag_value` update, the provider SHALL update the resource id because the composite id embeds `tag_value` (format: `tagKey#tagValue#resource`).

#### Scenario: id is updated to reflect new tag_value
- **WHEN** the `UpdateResourceTagValue` API call succeeds
- **THEN** the provider SHALL set the id to `tagKey + FILED_SP + newTagValue + FILED_SP + resource`
- **AND** the subsequent `Read` SHALL use the new id to look up the binding

### Requirement: ModifyTagTagAttachmentValue service method
The `TagService` SHALL provide a `ModifyTagTagAttachmentValue` method that wraps the `UpdateResourceTagValue` API call.

#### Scenario: Service method constructs and sends the request
- **WHEN** `ModifyTagTagAttachmentValue(ctx, tagKey, tagValue, resource)` is called
- **THEN** the method SHALL construct an `UpdateResourceTagValueRequest` with `TagKey`, `TagValue`, and `Resource`
- **AND** the method SHALL call `UpdateResourceTagValue` on the tag client
- **AND** the method SHALL return an error if the API call fails

#### Scenario: Retry on transient failures
- **WHEN** the `UpdateResourceTagValue` API call fails with a retryable error
- **THEN** the resource `Update` function SHALL wrap the call with `resource.Retry(tccommon.WriteRetryTimeout, ...)` and `tccommon.RetryError`
