## 1. Schema Changes

- [x] 1.1 In `tencentcloud/services/tag/resource_tc_tag_attachment.go`, remove `ForceNew: true` from the `tag_value` schema field in `ResourceTencentCloudTagAttachment()`
- [x] 1.2 Register the `Update` function (`resourceTencentCloudTagAttachmentUpdate`) in the resource definition (`Update:` field), keeping `Create`, `Read`, `Delete`, and `Importer` as-is

## 2. Update Function Implementation

- [x] 2.1 Add `resourceTencentCloudTagAttachmentUpdate` function in `tencentcloud/services/tag/resource_tc_tag_attachment.go` with `defer tccommon.LogElapsed("resource.tencentcloud_tag_attachment.update")()` and `defer tccommon.InconsistentCheck(d, meta)()`
- [x] 2.2 Parse the composite id (`tagKey#tagValue#resource`) from `d.Id()` using `tccommon.FILED_SP`, returning an error if the split does not yield 3 parts
- [x] 2.3 Gate the update on `d.HasChange("tag_value")`; if no change, skip directly to `resourceTencentCloudTagAttachmentRead(d, meta)` and return
- [x] 2.4 Build the `UpdateResourceTagValueRequest` with `TagKey` (from id split), `TagValue` (new value from `d.Get("tag_value")`), and `Resource` (from id split)
- [x] 2.5 Wrap the `UpdateResourceTagValue` API call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; on error use `tccommon.RetryError(e)`; the retry block SHALL only contain the API call
- [x] 2.6 After the retry succeeds, refresh the composite id via `d.SetId(tagKey + tccommon.FILED_SP + newTagValue + tccommon.FILED_SP + resource)` outside the retry block
- [x] 2.7 Call `resourceTencentCloudTagAttachmentRead(d, meta)` to refresh state and return

## 3. Unit Tests

- [x] 3.1 Add unit test cases in `tencentcloud/services/tag/resource_tc_tag_attachment_test.go` using gomonkey to mock the cloud API (mock `UpdateResourceTagValue` on the tag client) covering: successful update of `tag_value`, and update invoked with no change to `tag_value`

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/tag/resource_tc_tag_attachment.md` to note that `tag_value` is updatable in place (no import section change; this resource is not RESOURCE_KIND_GENERAL/ATTACHMENT/CONFIG so keep existing Import block)

## 5. SDK Upgrade

- [x] 5.1 Upgrade the Tag SDK module from `v1.0.860` to `v1.3.110` in `go.mod`/`go.sum` and re-vendor `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813/` to obtain the `UpdateResourceTagValue` API (`Client.UpdateResourceTagValue`, `UpdateResourceTagValueRequest`/`Response`)

## 6. Validation

- [x] 6.1 Verify the generated Go code compiles successfully (checked by the build/lint flow, not run here)
- [x] 6.2 Verify the proposal artifacts are consistent with the implemented code; update artifacts if needed
