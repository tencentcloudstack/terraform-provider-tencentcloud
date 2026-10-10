## 1. Service Layer

- [x] 1.1 Append `DescribeAggregateConfigDeliver(ctx, accountGroupId string) (*configv20220802.DescribeAggregateConfigDeliverResponseParams, error)` method to `tencentcloud/services/config/service_tencentcloud_config.go`. Build request, set `request.AccountGroupId`, wrap call in `resource.Retry(tccommon.ReadRetryTimeout)` using `me.client.UseConfigV20220802Client().DescribeAggregateConfigDeliverWithContext(ctx, request)`, check `result == nil || result.Response == nil` → `NonRetryableError`, return `response.Response`.

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver.go` with `ResourceTencentCloudConfigUpdateAggregateConfigDeliver()` schema: `account_group_id` (string, required, ForceNew), `status` (int, required), `deliver_name`/`target_arn`/`deliver_prefix`/`deliver_type` (string, optional), `deliver_uin` (int, optional), `deliver_content_type` (int, optional), `create_time` (string, computed); include `Importer` (ImportStatePassthrough).
- [x] 2.2 Implement Create handler: read `account_group_id` from `d.Get`, `d.SetId(account_group_id)`, log id, then delegate to Update.
- [x] 2.3 Implement Read handler: call `service.DescribeAggregateConfigDeliver(ctx, d.Id())`; guard `respData == nil` → `log.Printf("[CRUD] ...")` + `d.SetId("")`; set each non-nil field to schema.
- [x] 2.4 Implement Update handler: build `UpdateAggregateConfigDeliverRequest`, set `AccountGroupId` from `d.Id()`, populate mutable fields from `d.GetOk`/`d.GetOkExists`, wrap in `resource.Retry(tccommon.WriteRetryTimeout)` with `tccommon.RetryError(e)`, then call Read to refresh.
- [x] 2.5 Implement Delete handler: no-op (return nil).

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_config_update_aggregate_config_deliver` in `tencentcloud/provider.go` ResourcesMap (below the existing `tencentcloud_config_deliver_config` entry).
- [x] 3.2 Add `tencentcloud_config_update_aggregate_config_deliver` to `tencentcloud/provider.md` resource list.

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver.md` with one-line description (mentioning Config), Example Usage (HCL), and Import section (using account_group_id as the import id).

## 5. Tests

- [x] 5.1 Create `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver_test.go` with gomonkey-based mock unit tests covering Create/Read/Update/Delete business logic (no terraform test suite, per project rules for new resources).