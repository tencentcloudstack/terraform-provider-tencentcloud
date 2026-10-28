## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/cls/resource_tc_cls_open_service_operation.go` with constructor `ResourceTencentCloudClsOpenServiceOperation()`, operation-type code style following `tencentcloud_dlc_update_row_filter_operation`
- [x] 1.2 Define schema: computed `status` (Int) only (both APIs take no input parameters)
- [x] 1.3 Implement Create: call `OpenClsService` inside `resource.Retry(WriteRetryTimeout, ...)`, then `d.SetId(helper.BuildToken())` and call Read
- [x] 1.4 Implement Read: call `GetClsService` inside `resource.Retry(ReadRetryTimeout, ...)`, nil-safe set computed `status`
- [x] 1.5 Implement Delete: no-op (return nil)
- [x] 1.6 Ensure all response-value access is nil-safe and every API call uses the retry mechanism

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_cls_open_service_operation` in `tencentcloud/provider.go` ResourcesMap
- [x] 2.2 Add the resource entry to `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/cls/resource_tc_cls_open_service_operation.md` example file, format following `resource_tc_dlc_update_row_filter_operation.md` (Example Usage + Import)
- [x] 3.2 Generate `website/docs/r/cls_open_service_operation.html.markdown` via `make doc` (do not hand-write) and confirm `website/tencentcloud.erb` link entry

## 4. Unit Test

- [x] 4.1 Create `tencentcloud/services/cls/resource_tc_cls_open_service_operation_test.go`, naming/format following `resource_tc_dlc_update_row_filter_operation_test.go`

## 5. Verification

- [x] 5.1 Run `gofmt`/`go build ./tencentcloud/...` and ensure no compile errors
- [x] 5.2 Run `go vet ./tencentcloud/services/cls/` and `read_lints` to confirm no newly introduced errors
- [x] 5.3 Run `make doc` and verify generated website doc is consistent with the `.md` example
