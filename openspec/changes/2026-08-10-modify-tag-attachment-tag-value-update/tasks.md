## 1. Schema Changes

- [x] 1.1 In `tencentcloud/services/tag/resource_tc_tag_attachment.go`, remove `ForceNew: true` from the `tag_value` schema field (keep it `Required`, `TypeString`)
- [x] 1.2 Keep `tag_key` and `resource` schema fields as `ForceNew: true` (unchanged)
- [x] 1.3 Register the `Update` function in the `schema.Resource` returned by `ResourceTencentCloudTagAttachment()` (add `Update: resourceTencentCloudTagAttachmentUpdate`)

## 2. Update Function Implementation

- [x] 2.1 Add `resourceTencentCloudTagAttachmentUpdate` function in `tencentcloud/services/tag/resource_tc_tag_attachment.go` with `defer tccommon.LogElapsed("resource.tencentcloud_tag_attachment.update")()` and `defer tccommon.InconsistentCheck(d, meta)()`
- [x] 2.2 Parse the composite id via `strings.Split(d.Id(), tccommon.FILED_SP)`; return an error if the split length is not 3 (id broken)
- [x] 2.3 Obtain the old and new `tag_value` via `oldValue, newValue := d.GetChange("tag_value")`; read `tag_key` and `resource` from the id split
- [x] 2.4 Build the update payload: `replaceTags = {tag_key: newValue}` and `deleteKeys = nil` — per the `ModifyResourceTags` API contract, placing `{tag_key, newValue}` in `ReplaceTags` changes the existing value in place; `DeleteTags` is empty because `ReplaceTags` and `DeleteTags` must not share the same tag key
- [x] 2.5 Call `service.ModifyTags(ctx, resource, replaceTags, nil)` to invoke `ModifyResourceTags` atomically in a single request (reuse existing service method; do NOT add a new service method)
- [x] 2.6 After the `ModifyTags` call succeeds (no error), refresh the composite id: `d.SetId(tagKey + tccommon.FILED_SP + newValue.(string) + tccommon.FILED_SP + resource)`
- [x] 2.7 Return `resourceTencentCloudTagAttachmentRead(d, meta)` to refresh state

## 3. Unit Tests

- [x] 3.1 In `tencentcloud/services/tag/resource_tc_tag_attachment_test.go`, add a test step for the `Update` path that changes `tag_value` and verifies the attachment still exists with the new value (via the terraform test suite, per project convention for modified resources); the underlying `ModifyResourceTags` is called with `ReplaceTags` containing the new value and `DeleteTags` empty
- [x] 3.2 Add a unit test verifying the composite id is refreshed to `tagKey + FILED_SP + new_tag_value + FILED_SP + resource` after a successful update
- [x] 3.3 Add a unit test verifying that changing `tag_key` or `resource` still forces replacement (ForceNew behavior preserved)

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/tag/resource_tc_tag_attachment.md` example to reflect that `tag_value` can be updated in-place (add an example showing tag_value change)

## 5. Validation

- [x] 5.1 Verify the code compiles successfully (no `go build`/`go vet` commands run by hand; the build/lint step is executed by the downstream pipeline)
- [x] 5.2 Verify no lint errors (deferred to downstream pipeline)
