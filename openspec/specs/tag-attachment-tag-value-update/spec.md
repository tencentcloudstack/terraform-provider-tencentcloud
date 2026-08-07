# tag-attachment-tag-value-update Specification

## Purpose
TBD - created by archiving change update-tag-attachment-tag-value. Update Purpose after archive.
## Requirements
### Requirement: tag_value is updatable in place
The `tencentcloud_tag_attachment` resource SHALL allow the `tag_value` field to be updated in place. The `tag_value` field SHALL NOT be marked `ForceNew`. Changes to `tag_value` SHALL trigger the resource `Update` function instead of destroying and recreating the attachment.

#### Scenario: Changing tag_value triggers Update
- **WHEN** a user changes only the `tag_value` field of an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL call the `Update` function
- **AND** the provider SHALL NOT delete and recreate the attachment

#### Scenario: tag_value is not ForceNew
- **WHEN** the `tencentcloud_tag_attachment` schema is defined
- **THEN** the `tag_value` field SHALL NOT have `ForceNew: true`

### Requirement: tag_key and resource remain immutable
The `tag_key` and `resource` fields of `tencentcloud_tag_attachment` SHALL remain `ForceNew: true`. Changing either field SHALL destroy and recreate the attachment as before.

#### Scenario: Changing tag_key recreates the attachment
- **WHEN** a user changes the `tag_key` field of an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL destroy and recreate the attachment

#### Scenario: Changing resource recreates the attachment
- **WHEN** a user changes the `resource` field of an existing `tencentcloud_tag_attachment` resource
- **THEN** the provider SHALL destroy and recreate the attachment

### Requirement: Update uses the UpdateResourceTagValue API in a single request
The `tencentcloud_tag_attachment` `Update` function SHALL modify the bound tag value in a single API request using the tag service `UpdateResourceTagValue` API. The request SHALL pass `TagKey`, `TagValue` (the new value), and `Resource` (the six-segment QCS string) fields. The provider SHALL NOT implement the update as a delete-then-add sequence.

#### Scenario: Update calls UpdateResourceTagValue with the new tag value
- **WHEN** the `Update` function is invoked because `tag_value` changed
- **THEN** the provider SHALL build an `UpdateResourceTagValueRequest` with `TagKey` set to the current `tag_key`, `TagValue` set to the new `tag_value`, and `Resource` set to the current `resource`
- **AND** the provider SHALL send the request to the `UpdateResourceTagValue` API in a single call

#### Scenario: Update does not delete then add
- **WHEN** the `Update` function is invoked because `tag_value` changed
- **THEN** the provider SHALL NOT call `DeleteResourceTag` followed by `AddResourceTag`

### Requirement: UpdateResourceTagValue call is retried with WriteRetryTimeout
The `Update` function SHALL wrap the `UpdateResourceTagValue` API call in `resource.Retry(tccommon.WriteRetryTimeout, ...)` and SHALL convert errors using `tccommon.RetryError(err)`. The retry block SHALL only perform the API call; setting the id and other success handling SHALL happen outside the retry block.

#### Scenario: Transient API failure is retried
- **WHEN** the `UpdateResourceTagValue` API returns a transient error
- **THEN** the provider SHALL retry the call within `tccommon.WriteRetryTimeout`
- **AND** the id SHALL NOT be set inside the retry block

#### Scenario: Persistent API failure is returned
- **WHEN** the `UpdateResourceTagValue` API returns a non-retryable error
- **THEN** the provider SHALL return the wrapped error to the caller

### Requirement: Composite id is refreshed after a successful update
The `tencentcloud_tag_attachment` composite id is `tagKey + FILED_SP + tagValue + FILED_SP + resource`. After a successful `UpdateResourceTagValue` call, the `Update` function SHALL recompute and set `d.SetId()` with the new `tag_value` so the id reflects the updated binding. The id format SHALL remain unchanged.

#### Scenario: Id reflects new tag value after update
- **WHEN** the `UpdateResourceTagValue` call succeeds for a tag value change from `A` to `B`
- **THEN** the provider SHALL set `d.SetId(tagKey + FILED_SP + "B" + FILED_SP + resource)`

#### Scenario: Id format is preserved
- **WHEN** the composite id is recomputed after update
- **THEN** the id SHALL use the existing `FILED_SP` separator and three-segment structure

### Requirement: Update reads tag_key and resource from the current id
The `Update` function SHALL parse the current composite id (`tagKey#oldTagValue#resource`) to validate the id and obtain the existing `tag_key` and `resource`. It SHALL read the new `tag_value` from `d.Get("tag_value")`. If the id cannot be split into three segments, the `Update` function SHALL return an error.

#### Scenario: tag_key and resource are obtained from the id
- **WHEN** the `Update` function runs
- **THEN** the provider SHALL split `d.Id()` by `FILED_SP` into three parts: `tagKey`, the old `tagValue`, and `resource`
- **AND** the provider SHALL use `tagKey` (first segment) and `resource` (third segment) to build the `UpdateResourceTagValue` request

#### Scenario: Broken id returns an error
- **WHEN** `d.Id()` does not split into exactly three segments by `FILED_SP`
- **THEN** the `Update` function SHALL return an error of the form `id is broken,%s`

### Requirement: Service layer provides UpdateTagTagValue helper
The tag service layer SHALL provide an `UpdateTagTagValue(ctx, tagKey, tagValue, resourceName string) error` function that builds an `UpdateResourceTagValueRequest`, applies `ratelimit.Check`, and calls `UpdateResourceTagValue` wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(err)` on failure.

#### Scenario: Helper builds the request from schema fields
- **WHEN** `UpdateTagTagValue` is called with `tagKey`, `tagValue`, and `resourceName`
- **THEN** the service function SHALL set `request.TagKey`, `request.TagValue`, and `request.Resource` accordingly
- **AND** the service function SHALL send the `UpdateResourceTagValue` request

#### Scenario: Helper applies rate limiting and retry
- **WHEN** `UpdateTagTagValue` calls the API
- **THEN** the service function SHALL call `ratelimit.Check(request.GetAction())`
- **AND** the service function SHALL wrap the call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`

### Requirement: tag_value in-place update is tested
The resource test file `resource_tc_tag_attachment_test.go` SHALL include an acceptance test step that updates the `tag_value` of an existing `tencentcloud_tag_attachment` and verifies the updated value is reflected in state.

#### Scenario: Test updates tag_value in place
- **WHEN** the acceptance test runs an update step that changes `tag_value`
- **THEN** the provider SHALL update the attachment in place
- **AND** the test SHALL verify `tag_value` matches the new value in state

