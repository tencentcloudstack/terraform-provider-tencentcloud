## 1. Schema Definition & Resource Structure

- [x] 1.1 Create resource file `tencentcloud/services/cbs/resource_tc_cbs_copy_snapshot_cross_region_attachment.go` with the resource definition function `ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment()` returning `*schema.Resource`
- [x] 1.2 Define schema fields: `snapshot_id` (Required, ForceNew), `destination_regions` (Required, ForceNew, TypeList/TypeString), `snapshot_name` (Optional, ForceNew), `delete_bind_images` (Optional, ForceNew, default false), `snapshot_copy_result_set` (Computed, TypeList with nested schema containing snapshot_id, code, message, destination_region)
- [x] 1.3 Add Timeouts block to the schema for async operation support (Create timeout)
- [x] 1.4 Register Create, Read, Delete functions (no Update) and Importer with ImportStatePassthrough

## 2. Create Function Implementation

- [x] 2.1 Implement `resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentCreate` function with `defer tccommon.LogElapsed()` and context setup
- [x] 2.2 Build CopySnapshotCrossRegions request with snapshot_id, destination_regions, and optional snapshot_name from schema
- [x] 2.3 Call CopySnapshotCrossRegions API inside `resource.Retry(tccommon.WriteRetryTimeout)` wrapper
- [x] 2.4 Check response nil-ness (Response, SnapshotCopyResultSet) and return NonRetryableError if empty, log logId and d.Id() before checking
- [x] 2.5 Check each SnapshotCopyResult entry's Code field; if not "Success", return error with Code and Message
- [x] 2.6 Set composite ID using `snapshot_id + FILED_SP + copied_snapshot_id` (using the first copied snapshot ID from the result set)
- [x] 2.7 Implement async polling: for each copied snapshot_id in SnapshotCopyResultSet, poll DescribeSnapshots until SnapshotState is NORMAL using `resource.Retry` with `tccommon.ReadRetryTimeout`
- [x] 2.8 Save snapshot_copy_result_set to state with all fields (snapshot_id, code, message, destination_region) from each SnapshotCopyResult
- [x] 2.9 Call Read function at the end of Create

## 3. Read Function Implementation

- [x] 3.1 Implement `resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentRead` function with `defer tccommon.LogElapsed()` and `defer tccommon.InconsistentCheck()`
- [x] 3.2 Parse composite ID from d.Id() using FILED_SP separator, validate format, extract snapshot_id and copied_snapshot_id
- [x] 3.3 Call DescribeSnapshots API inside `resource.Retry(tccommon.ReadRetryTimeout)` using copied_snapshot_id as SnapshotIds filter
- [x] 3.4 If DescribeSnapshots returns empty result, log `log.Printf("[CRUD] cbs_copy_snapshot_cross_region id=%s", d.Id())` then set d.SetId("")
- [x] 3.5 If snapshot found, populate state fields: snapshot_id, snapshot_name, and snapshot_copy_result_set from the response

## 4. Delete Function Implementation

- [x] 4.1 Implement `resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentDelete` function with `defer tccommon.LogElapsed()`
- [x] 4.2 Read all copied snapshot IDs from state's snapshot_copy_result_set field
- [x] 4.3 Build DeleteSnapshots request with SnapshotIds list and DeleteBindImages flag from schema
- [x] 4.4 Call DeleteSnapshots API inside `resource.Retry(tccommon.WriteRetryTimeout)` wrapper

## 5. Provider Registration & Documentation

- [x] 5.1 Register resource in `tencentcloud/provider.go` ResourcesMap with key `"tencentcloud_cbs_copy_snapshot_cross_region"` and value `cbs.ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment()`
- [x] 5.2 Update `tencentcloud/provider.go` file comment to include the new resource name under CBS Resource section
- [x] 5.3 Update `tencentcloud/provider.md` to include `tencentcloud_cbs_copy_snapshot_cross_region` under Cloud Block Storage(CBS) Resource section
- [x] 5.4 Create documentation file `tencentcloud/services/cbs/resource_tc_cbs_copy_snapshot_cross_region_attachment.md` with one-line description ("Provides a CBS copy snapshot cross region resource."), Example Usage section, and Import section (since this is RESOURCE_KIND_ATTACHMENT type)

## 6. Unit Tests

- [x] 6.1 Create test file `tencentcloud/services/cbs/resource_tc_cbs_copy_snapshot_cross_region_attachment_test.go`
- [x] 6.2 Implement Create function unit test using gomonkey to mock CopySnapshotCrossRegions and DescribeSnapshots API calls, verify composite ID construction and snapshot_copy_result_set state
- [x] 6.3 Implement Read function unit test using gomonkey to mock DescribeSnapshots API call, verify ID parsing and state population
- [x] 6.4 Implement Delete function unit test using gomonkey to mock DeleteSnapshots API call, verify copied snapshot IDs extraction and delete request construction
- [x] 6.5 Run unit tests with `go test -gcflags=all=-l` to verify all test cases pass

## 7. Final Verification

- [x] 7.1 Verify all generated Go code compiles correctly (check imports, function signatures, struct usage match the CBS SDK)
- [x] 7.2 Verify schema field types and descriptions are correct
- [x] 7.3 Verify API parameter mapping matches the CBS SDK request/response structures (CopySnapshotCrossRegions, DescribeSnapshots, DeleteSnapshots)
- [x] 7.4 Verify composite ID format and parsing logic are consistent across Create, Read, and Delete functions
