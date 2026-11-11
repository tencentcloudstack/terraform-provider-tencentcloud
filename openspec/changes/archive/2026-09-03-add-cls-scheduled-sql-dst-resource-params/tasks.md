## 1. Schema Definition

- [x] 1.1 Add `metric_names` (Optional, TypeList of TypeString) to the `dst_resource` block schema in `tencentcloud/services/cls/resource_tc_cls_scheduled_sql.go`
- [x] 1.2 Add `metric_labels` (Optional, TypeList of TypeString) to the `dst_resource` block schema
- [x] 1.3 Add `custom_time` (Optional, TypeString) to the `dst_resource` block schema
- [x] 1.4 Add `custom_metric_labels` (Optional, TypeList of Resource with `key` and `value` string fields) to the `dst_resource` block schema

## 2. Create Operation

- [x] 2.1 In `resourceTencentCloudClsScheduledSqlCreate`, add handling for `metric_names` — convert the schema list to `[]*string` and assign to `scheduledSqlResouceInfo.MetricNames`
- [x] 2.2 In `resourceTencentCloudClsScheduledSqlCreate`, add handling for `metric_labels` — convert the schema list to `[]*string` and assign to `scheduledSqlResouceInfo.MetricLabels`
- [x] 2.3 In `resourceTencentCloudClsScheduledSqlCreate`, add handling for `custom_time` — convert the schema string and assign to `scheduledSqlResouceInfo.CustomTime`
- [x] 2.4 In `resourceTencentCloudClsScheduledSqlCreate`, add handling for `custom_metric_labels` — iterate the schema list, build `[]*cls.MetricLabel` with `Key` and `Value`, and assign to `scheduledSqlResouceInfo.CustomMetricLabels`

## 3. Read Operation

- [x] 3.1 In `resourceTencentCloudClsScheduledSqlRead`, add nil-checked reading of `DstResource.MetricNames` and set `metric_names` in the `dst_resource` state map
- [x] 3.2 In `resourceTencentCloudClsScheduledSqlRead`, add nil-checked reading of `DstResource.MetricLabels` and set `metric_labels` in the `dst_resource` state map
- [x] 3.3 In `resourceTencentCloudClsScheduledSqlRead`, add nil-checked reading of `DstResource.CustomTime` and set `custom_time` in the `dst_resource` state map
- [x] 3.4 In `resourceTencentCloudClsScheduledSqlRead`, add nil-checked reading of `DstResource.CustomMetricLabels` and set `custom_metric_labels` (as a list of maps with `key` and `value`) in the `dst_resource` state map

## 4. Update Operation

- [x] 4.1 In `resourceTencentCloudClsScheduledSqlUpdate`, verify that `dst_resource` is already in the `immutableArgs` array — no additional immutableArgs changes needed since the new fields are sub-fields of `dst_resource`
- [x] 4.2 In the `d.HasChange("dst_resource")` block of Update, add handling for `metric_names`, `metric_labels`, `custom_time`, and `custom_metric_labels` — mirror the Create operation logic for building `scheduledSqlResouceInfo`

## 5. Unit Tests

- [x] 5.1 In `tencentcloud/services/cls/resource_tc_cls_scheduled_sql_test.go`, add a unit test case (using gomonkey mock) for Create with all new parameters (`metric_names`, `metric_labels`, `custom_time`, `custom_metric_labels`)
- [x] 5.2 Add a unit test case for Read populating the new parameters from the Describe response
- [x] 5.3 Add a unit test case for Update with the new parameters in `dst_resource`

## 6. Documentation

- [x] 6.1 Update `tencentcloud/services/cls/resource_tc_cls_scheduled_sql.md` — add an example usage demonstrating `metric_names`, `metric_labels`, `custom_time`, and `custom_metric_labels` within the `dst_resource` block

## 7. Finalization

- [ ] 7.1 Run `gofmt` formatting on modified Go files (via tfpacer-finalize skill)
- [ ] 7.2 Run `make doc` to generate website documentation (via tfpacer-finalize skill)
- [ ] 7.3 Generate changelog file (via tfpacer-finalize skill)
