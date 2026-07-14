## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.go` with constructor `ResourceTencentCloudTeoIncreasePlanQuotaOperation()`, following the operation-type code style of `resource_tc_teo_confirm_origin_acl_update_operation.go`
- [x] 1.2 Define schema: `plan_id` (TypeString, Required, ForceNew), `quota_type` (TypeString, Required, ForceNew), `quota_number` (TypeInt, Required, ForceNew), `deal_name` (TypeString, Computed)
- [x] 1.3 Implement Create: call `IncreasePlanQuota` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`, nil-safe set `deal_name` from response, `d.SetId(helper.BuildToken())`, then call Read
- [x] 1.4 Implement Read: no-op (return nil)
- [x] 1.5 Implement Delete: no-op (return nil)
- [x] 1.6 Ensure all response-value access is nil-safe and every API call uses the retry mechanism

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_increase_plan_quota` in `tencentcloud/provider.go` ResourcesMap (alphabetically among teo resources)
- [x] 2.2 Add the resource entry to `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.md` example file, format following `resource_tc_teo_confirm_origin_acl_update_operation.md`
- [ ] 3.2 Generate `website/docs/r/teo_increase_plan_quota_operation.html.markdown` via `make doc` and confirm `website/tencentcloud.erb` link entry (deferred to tfpacer-finalize)

## 4. Unit Test

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation_test.go`, using gomonkey to mock the `IncreasePlanQuota` API call, following the pattern of existing TEO operation tests (e.g., `resource_tc_teo_import_zone_config_operation_test.go`)

## 5. Verification

- [ ] 5.1 Run `go build ./tencentcloud/...` and ensure no compile errors (deferred to tfpacer-finalize)
- [ ] 5.2 Run `go vet ./tencentcloud/services/teo/` and confirm no newly introduced errors (deferred to tfpacer-finalize)
- [ ] 5.3 Run `make doc` and verify generated website doc is consistent with the `.md` example (deferred to tfpacer-finalize)