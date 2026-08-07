## 1. Schema Changes

- [x] 1.1 Remove `ForceNew: true` from the `tag_value` schema field in `ResourceTencentCloudTagAttachment()` in `tencentcloud/services/tag/resource_tc_tag_attachment.go` (keep `ForceNew` on `tag_key` and `resource`)
- [x] 1.2 Register the `Update: resourceTencentCloudTagAttachmentUpdate` callback in the `schema.Resource` definition

## 2. Service Layer Changes

- [x] 2.1 Add `UpdateTagAttachmentTagValue` function to `tencentcloud/services/tag/service_tencentcloud_tag.go` that wraps the `UpdateResourceTagValue` API call with `resource.Retry(tccommon.WriteRetryTimeout, ...)`, `ratelimit.Check`, and `tccommon.RetryError` error wrapping

## 3. Update Function Implementation

- [x] 3.1 Implement `resourceTencentCloudTagAttachmentUpdate` in `tencentcloud/services/tag/resource_tc_tag_attachment.go` that parses the current composite ID into `tagKey`, `tagValue`, `resource`
- [x] 3.2 In the Update function, detect when `tag_value` has changed (`d.HasChange("tag_value")`); when changed, build the `UpdateResourceTagValueRequest` with `TagKey`, the new `TagValue`, and `Resource`, and call `UpdateTagAttachmentTagValue`
- [x] 3.3 After a successful update call, re-set the composite ID via `d.SetId(tagKey + tccommon.FILED_SP + newTagValue + tccommon.FILED_SP + resource)`
- [x] 3.4 Call `resourceTencentCloudTagAttachmentRead(d, meta)` at the end of Update to refresh state

## 4. Unit Tests

- [x] 4.1 Add unit test cases in `tencentcloud/services/tag/resource_tc_tag_attachment_test.go` (using gomonkey mocks) covering the tag-value update flow: successful update via `UpdateResourceTagValue`, ID re-set verification, and update-error propagation

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/tag/resource_tc_tag_attachment.md` to reflect that `tag_value` is updatable (remove ForceNew note; add example showing tag_value update)

## 6. Validation

- [x] 6.1 Verify the code compiles successfully
- [x] 6.2 Verify no lint errors
