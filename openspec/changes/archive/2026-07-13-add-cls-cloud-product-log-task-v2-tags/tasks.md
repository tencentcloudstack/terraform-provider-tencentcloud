## 1. Schema Definition

- [x] 1.1 Add the `tags` parameter (TypeMap, Optional, Elem string) to the `ResourceTencentCloudClsCloudProductLogTaskV2()` schema map in `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`, with a description noting tags are bound to the associated logset and topic.
- [x] 1.2 Verify the `tags` parameter is NOT marked `ForceNew` (it must be mutable in-place via the modify API).

## 2. Create Logic

- [x] 2.1 In `resourceTencentCloudClsCloudProductLogTaskV2Create`, read the `tags` map from `d.GetOk("tags")` and convert it to a `[]*clsv20201016.Tag` slice (Key/Value string pointers), then assign to `request.Tags` before calling `CreateCloudProductLogCollection`.
- [x] 2.2 Ensure that when `tags` is not set, `request.Tags` is left nil (no behavior change for existing users).

## 3. Read Logic

- [x] 3.1 In `resourceTencentCloudClsCloudProductLogTaskV2Read`, after fetching the task via `DescribeClsCloudProductLogTaskById`, read tags from the response `CloudProductLogTaskInfo` (`TopicTags`, falling back to `LogsetTags` if `TopicTags` is empty) and flatten them into the `tags` map in state via `d.Set("tags", ...)`.
- [x] 3.2 Guard the tag read with nil checks (only set when the tags slice is non-nil and non-empty), consistent with the existing nil-guard pattern for other response fields.

## 4. Update Logic

- [x] 4.1 In `resourceTencentCloudClsCloudProductLogTaskV2Update`, add `"tags"` to the existing `mutableArgs` slice so a change to `tags` triggers the modify path.
- [x] 4.2 Inside the `needChange` block, read the `tags` map from `d.GetOk("tags")`, convert to `[]*clsv20201016.Tag`, and assign to `request.Tags` on the `ModifyCloudProductLogCollectionRequest` before the retry call. When tags are removed (empty map), send an empty slice so the API clears them.

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.md` to add a usage example demonstrating the `tags` parameter (both the default-new-logset/topic example and/or the existing-logset/topic example).
- [x] 5.2 Do NOT manually edit `website/docs/` files; the final website docs are generated via `make doc` during the finalize phase.

## 6. Unit Tests (gomonkey mocks)

- [x] 6.1 Add a unit test in `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2_test.go` (using gomonkey mocks, not the terraform test suite) covering the create flow with `tags` set, verifying `request.Tags` is populated correctly.
- [x] 6.2 Add a unit test covering the update flow when `tags` changes, verifying `ModifyCloudProductLogCollection` is called with the updated `Tags`.
- [x] 6.3 Add a unit test covering the read flow, verifying tags from the API response are flattened into state.
- [x] 6.4 Run the new unit tests with `go test -gcflags=all=-l` on the affected test file and ensure they pass.

## 7. Verification

- [x] 7.1 Verify the modified `resource_tc_cls_cloud_product_log_task_v2.go` compiles (handled by the downstream build/lint process; do NOT run `go build`/`go vet` manually).
- [x] 7.2 Verify all functions returning an error are checked (assign to `_ =` if the error is guaranteed nil) to avoid unused-variable errors.
- [x] 7.3 Confirm no existing schema field behavior changed (all existing parameters retain their Required/Optional/ForceNew/Computed settings).
