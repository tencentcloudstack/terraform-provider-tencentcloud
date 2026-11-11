## 1. Service Layer

- [x] 1.1 Append `DescribeDBCustomDisasterRecoverGroupById(ctx, disasterRecoverGroupId string)` to `service_tencentcloud_dbdc.go` — wraps `DescribeDBCustomDisasterRecoverGroups` with `DisasterRecoverGroupIds=[id]`, `Limit=100` (documented max), `resource.Retry(tccommon.ReadRetryTimeout)` + `ratelimit.Check`; returns the single matching `*dbdcv20201029.DisasterRecoverGroup` (nil if not found), nil/length-safe
- [x] 1.2 Reuse existing `waitDBCustomTaskSucceeded` helper (already in `resource_tc_dbdc_db_custom_cluster.go`, same `dbdc` package) for async delete polling

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_dbdc_db_custom_disaster_recover_group.go` with `ResourceTencentCloudDbdcDbCustomDisasterRecoverGroup()` schema: required `name`; optional `type`(ForceNew), `strategy`(ForceNew), `affinity`(mutable), `tags`(ForceNew, TypeList of {key,value}), `client_token`(ForceNew); computed `disaster_recover_group_id`, `status`, `node_quota_total`, `current_num`, `created_time`, `node_ids`; `Importer` passthrough; `Timeouts` block (Create/Delete)
- [x] 2.2 Implement Create: build `CreateDBCustomDisasterRecoverGroupRequest` (name, type, strategy, affinity, tags[], client_token), `resource.Retry(WriteRetryTimeout)` calling `CreateDBCustomDisasterRecoverGroupWithContext`; guard `result==nil||result.Response==nil`; check `DisasterRecoverGroupId` non-nil/non-empty (print logId + id first, return `NonRetryableError` if empty); `d.SetId(disasterRecoverGroupId)`; poll `DescribeDBCustomDisasterRecoverGroupById` until `Status=="Available"` (NonRetryableError on `CreateFailed`); return Read
- [x] 2.3 Implement Read: call `DescribeDBCustomDisasterRecoverGroupById(ctx, d.Id())`; if nil → `log.Printf("[CRUD] disaster_recover_group id=%s", d.Id())` then `d.SetId("")`; set each field with nil guards (name, type, strategy, affinity, tags→list of {key,value}, status, node_quota_total, current_num, created_time, node_ids, disaster_recover_group_id)
- [x] 2.4 Implement Update: build `immutableArgs` = ["type","strategy","tags","client_token"]; if any has change → return error; if `name` or `affinity` changed → `ModifyDBCustomDisasterRecoverGroupAttributeRequest{DisasterRecoverGroupId, Name, Affinity}`, `resource.Retry(WriteRetryTimeout)`; then return Read
- [x] 2.5 Implement Delete: `DeleteDBCustomDisasterRecoverGroupsRequest{DisasterRecoverGroupIds=[d.Id()]}`, `resource.Retry(WriteRetryTimeout)`; guard response; if `TaskId != nil` → `waitDBCustomTaskSucceeded(ctx, &service, *TaskId, d.Timeout(schema.TimeoutDelete))`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_dbdc_db_custom_disaster_recover_group` in `provider.go` ResourcesMap (next to existing `tencentcloud_dbdc_db_custom_cluster` entries), reference `dbdc.ResourceTencentCloudDbdcDbCustomDisasterRecoverGroup()`

## 4. Documentation & Tests

- [x] 4.1 Create `resource_tc_dbdc_db_custom_disaster_recover_group.md` — one-line description mentioning DB Custom (dbdc), Example Usage HCL (name + optional type/strategy/affinity/tags), Import section (by `disaster_recover_group_id`); no Argument/Attribute Reference sections (auto-generated)
- [x] 4.2 Create `resource_tc_dbdc_db_custom_disaster_recover_group_test.go` — gomonkey mocks for `UseDbdcV20201029Client` (Create/Describe/Modify/Delete/TaskStatus), unit tests covering create-read-update-delete flow; check all returned errors; do NOT use terraform test suite

## 5. Verification

- [x] 5.1 Confirm generated `.go` files are compilable in the current environment (no `go build`/`go vet` execution; static review against SDK signatures and conventions)
- [x] 5.2 Confirm `provider.go` registration and `.md` doc are consistent with the `tencentcloud_igtm_strategy` reference style
