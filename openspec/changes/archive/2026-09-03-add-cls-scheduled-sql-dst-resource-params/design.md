## Context

The `tencentcloud_cls_scheduled_sql` resource manages CLS (Cloud Log Service) scheduled SQL analysis tasks. When the destination topic is a metric topic (`biz_type = 1`), the CLS API supports advanced metric configuration through the `ScheduledSqlResouceInfo` struct, including:

- `MetricNames` (`[]*string`): Multiple metric names for multi-metric scenarios
- `MetricLabels` (`[]*string`): Metric dimension labels (excluding time-type)
- `CustomTime` (`*string`): Custom timestamp field for metrics
- `CustomMetricLabels` (`[]*MetricLabel`): Static dimension labels with `Key`/`Value` pairs

Currently, the Terraform resource only exposes `topic_id`, `region`, `biz_type`, and `metric_name` within the `dst_resource` block. The SDK already includes all the necessary structs (`ScheduledSqlResouceInfo` at line 22989 and `MetricLabel` at line 16921 in `models.go`), so no vendor changes are needed.

## Goals / Non-Goals

**Goals:**
- Add `metric_names`, `metric_labels`, `custom_time`, and `custom_metric_labels` as Optional parameters inside the existing `dst_resource` block schema.
- Implement full CRUD support: Create, Read, and Update for all new parameters.
- Maintain backward compatibility — all new fields are Optional.
- Add unit test cases using gomonkey mock for the new parameters.

**Non-Goals:**
- No changes to the `metric_name` (singular) field — it remains as-is.
- No changes to the service layer `DescribeClsScheduledSqlById` function — it already returns the full `ScheduledSqlTaskInfo` including `DstResource`.
- No changes to the Delete operation — delete does not involve `DstResource`.
- No changes to vendor/SDK — all required structs already exist.

## Decisions

### Decision 1: Schema field types for list parameters

**Choice**: Use `schema.TypeList` with `schema.TypeString` elements for `metric_names` and `metric_labels`.

**Rationale**: The SDK fields `MetricNames` and `MetricLabels` are both `[]*string`. Using `TypeList` of `TypeString` is the standard Terraform pattern for string arrays and matches how the SDK expects the data. This is consistent with how other list-of-string parameters are handled in the provider.

### Decision 2: Schema structure for custom_metric_labels

**Choice**: Use `schema.TypeList` with a `schema.Resource` element containing `key` and `value` string fields.

**Rationale**: The SDK field `CustomMetricLabels` is `[]*MetricLabel` where `MetricLabel` has `Key` and `Value` string fields. A TypeList of Resource blocks is the idiomatic Terraform approach for a list of structured objects with fixed fields. This allows users to write repeatable HCL blocks:

```hcl
custom_metric_labels {
  key   = "env"
  value = "production"
}
```

### Decision 3: Nil-safety in Read operation

**Choice**: Check `DstResource.MetricNames`, `DstResource.MetricLabels`, `DstResource.CustomTime`, and `DstResource.CustomMetricLabels` for nil before reading, consistent with the existing pattern for `topic_id`, `region`, `biz_type`, and `metric_name`.

**Rationale**: The CLS API may return nil for any of these fields when they are not configured. Following the existing code pattern prevents panics and avoids setting empty values in state.

### Decision 4: Update handling via immutableArgs

**Choice**: The `dst_resource` block is already in the `immutableArgs` array in the Update method, meaning any change to `dst_resource` (including the new sub-fields) returns an error indicating the argument cannot be changed.

**Rationale**: The current resource treats `dst_resource` as immutable in updates. The existing Update method builds `DstResource` from the schema when `d.HasChange("dst_resource")` is true, but the `immutableArgs` check runs first and returns an error for any change to `dst_resource`. The new parameters are sub-fields of `dst_resource`, so they automatically follow the same immutability behavior. The `immutableArgs` array does not need modification since it already includes `dst_resource`.

### Decision 5: List-to-SDK conversion using helper functions

**Choice**: Use `helper.String(v.(string))` for individual string conversions and iterate over the list to build `[]*string` slices for `MetricNames` and `MetricLabels`. For `CustomMetricLabels`, iterate and build `[]*cls.MetricLabel`.

**Rationale**: This matches the existing code patterns in the resource file (e.g., how `metric_name` is handled with `helper.String()`).

## Risks / Trade-offs

- **[Risk] API returns nil for new fields when not configured** → Mitigation: All new fields are Optional and nil-checked in Read before setting state. Existing resources without these fields will have them absent from state, which is the expected behavior for Optional fields.

- **[Risk] `custom_metric_labels` ordering** → The SDK uses `[]*MetricLabel` (an ordered slice), and Terraform `TypeList` preserves order. No ordering risk since both sides are lists.

- **[Trade-off] immutability of dst_resource** → The `dst_resource` block is already immutable in the Update method. Users cannot change any `dst_resource` sub-field after creation. This is the existing behavior and the new fields inherit it. No change needed.
