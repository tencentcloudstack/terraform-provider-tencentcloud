## Why

The `tencentcloud_cls_scheduled_sql` resource currently only supports basic destination resource parameters (`topic_id`, `region`, `biz_type`, `metric_name`). When converting logs to metrics (BizType=1), users need to configure multiple metric names, metric labels, custom timestamps, and custom metric labels — all of which are already supported by the TencentCloud CLS API (`ScheduledSqlResouceInfo` struct in the SDK). Without these parameters, users cannot fully configure scheduled SQL tasks that convert logs to metrics through Terraform, forcing them to use the console or API directly.

## What Changes

- Add `metric_names` parameter (Optional, list of strings) to the `dst_resource` block of `tencentcloud_cls_scheduled_sql` — maps to `DstResource.MetricNames` in Create/Modify APIs and `DstResource.MetricNames` in Describe response.
- Add `metric_labels` parameter (Optional, list of strings) to the `dst_resource` block — maps to `DstResource.MetricLabels`.
- Add `custom_time` parameter (Optional, string) to the `dst_resource` block — maps to `DstResource.CustomTime`.
- Add `custom_metric_labels` parameter (Optional, list of objects with `key` and `value` string fields) to the `dst_resource` block — maps to `DstResource.CustomMetricLabels` (a list of `MetricLabel` structs).

All new parameters are Optional and backward compatible. No existing schema fields are modified or removed.

## Capabilities

### New Capabilities
- `cls-scheduled-sql-dst-resource-metric-params`: Adds metric-related destination resource parameters (`metric_names`, `metric_labels`, `custom_time`, `custom_metric_labels`) to the `tencentcloud_cls_scheduled_sql` resource for configuring log-to-metric conversion in scheduled SQL tasks.

### Modified Capabilities
<!-- None - no existing spec-level behavior changes -->

## Impact

- **Resource file**: `tencentcloud/services/cls/resource_tc_cls_scheduled_sql.go` — schema definition, Create, Read, and Update methods.
- **Resource test file**: `tencentcloud/services/cls/resource_tc_cls_scheduled_sql_test.go` — add unit test cases for new parameters using gomonkey mock.
- **Resource doc file**: `tencentcloud/services/cls/resource_tc_cls_scheduled_sql.md` — add usage examples with new parameters.
- **SDK dependency**: No changes needed — `ScheduledSqlResouceInfo`, `MetricLabel` structs already exist in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016/models.go`.
- **Backward compatibility**: Fully backward compatible — all new fields are Optional.
