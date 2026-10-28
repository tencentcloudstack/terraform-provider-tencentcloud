## Context

The `tencentcloud_monitor_tmp_instance` resource manages Prometheus instances via TencentCloud Monitor APIs. The cloud APIs support an `InstanceAttributes` parameter (type: `[]*PrometheusRuleKV`) that carries special instance properties. One supported attribute is `LongTermStorageRetentionTime` (archive storage retention in days, value range 60-730).

- `CreatePrometheusMultiTenantInstancePostPayMode` supports `InstanceAttributes` to set `LongTermStorageRetentionTime` at creation time.
- `ModifyPrometheusInstanceAttributes` supports `InstanceAttributes` to update `LongTermStorageRetentionTime` after creation.
- `DescribePrometheusInstances` returns `InstanceAttributes` in the `PrometheusInstancesItem` struct.

Current resource file: `tencentcloud/services/tmp/resource_tc_monitor_tmp_instance.go`
Service layer: `tencentcloud/services/monitor/service_tencentcloud_monitor.go`

## Goals / Non-Goals

**Goals:**
- Add `long_term_storage_retention_time` as an Optional integer parameter (range 60-730) to the resource schema.
- Pass `InstanceAttributes` with key `LongTermStorageRetentionTime` to the Create API when provided.
- Pass `InstanceAttributes` with key `LongTermStorageRetentionTime` to the Modify API when the value changes.
- Read back `LongTermStorageRetentionTime` from the Describe API response's `InstanceAttributes` array.
- Update unit tests to verify Create/Read/Update with `long_term_storage_retention_time`.

**Non-Goals:**
- Supporting other `InstanceAttributes` keys (e.g., `CreatedFrom`, `FreeTrialExpireAt`) — out of scope.
- Adding validation logic for value range in schema (rely on API-side validation).

## Decisions

1. **Parameter mutability**: The parameter is mutable after creation via `ModifyPrometheusInstanceAttributes` API, so updates are supported.

2. **Schema definition**: The parameter will be `Optional` + `Computed` (int type). Optional because users may or may not want to configure archive storage. Computed because the value can be read back from the API (may be set server-side).

3. **Read implementation**: Iterate over `InstanceAttributes` array from `PrometheusInstancesItem`, find the entry with key `LongTermStorageRetentionTime`, parse the string value as integer, and set in state.

4. **Wire format**: The API uses `PrometheusRuleKV` struct (`Key` string, `Value` string), so the integer value is converted to/from string for API communication.

## Risks / Trade-offs

- [Risk] The `InstanceAttributes` array in the API response may contain other keys besides `LongTermStorageRetentionTime` → Mitigation: Only read the specific key we care about, ignore others.
- [Risk] Backward compatibility concern if existing state files don't have this field → Mitigation: The field is Optional+Computed, so existing configurations without it will simply read the value from the API on next refresh without breaking.
