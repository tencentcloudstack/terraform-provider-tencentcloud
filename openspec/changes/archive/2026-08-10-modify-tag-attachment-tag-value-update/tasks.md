## 1. Service Layer Changes

- [x] 1.1 Add `ModifyTagTagAttachmentValue(ctx, tagKey, tagValue, resource string) (errRet error)` method to `TagService` in `tencentcloud/services/tag/service_tencentcloud_tag.go`
- [x] 1.2 In `ModifyTagTagAttachmentValue`, construct an `UpdateResourceTagValueRequest` with `TagKey`, `TagValue` (the new value), and `Resource`, call `UpdateResourceTagValue` on the tag client, and return the error on failure

## 2. Schema Changes

- [x] 2.1 In `ResourceTencentCloudTagAttachment()` in `tencentcloud/services/tag/resource_tc_tag_attachment.go`, remove `ForceNew: true` from the `tag_value` field
- [x] 2.2 Register the `Update` function (`resourceTencentCloudTagAttachmentUpdate`) in the resource definition

## 3. Update Function Implementation

- [x] 3.1 Implement `resourceTencentCloudTagAttachmentUpdate` in `tencentcloud/services/tag/resource_tc_tag_attachment.go` with `defer tccommon.LogElapsed()` and `defer tccommon.InconsistentCheck()`
- [x] 3.2 When `d.HasChange("tag_value")` is true, extract `tag_key`, the new `tag_value`, and `resource` from the schema data
- [x] 3.3 Wrap the `ModifyTagTagAttachmentValue` service call with `resource.Retry(tccommon.WriteRetryTimeout, ...)` and `tccommon.RetryError`
- [x] 3.4 After a successful API call, update the composite id to `tagKey + FILED_SP + newTagValue + FILED_SP + resource`
- [x] 3.5 Return `resourceTencentCloudTagAttachmentRead(d, meta)` at the end of the Update function

## 4. Test Changes

- [x] 4.1 Add an acceptance test step in `tencentcloud/services/tag/resource_tc_tag_attachment_test.go` that updates `tag_value` and verifies the in-place update succeeds (no destroy/recreate)

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/tag/resource_tc_tag_attachment.md` to reflect that `tag_value` is updatable (add an example showing tag_value modification)

## 6. Validation

- [x] 6.1 Verify the code compiles successfully
- [x] 6.2 Verify no lint errors
