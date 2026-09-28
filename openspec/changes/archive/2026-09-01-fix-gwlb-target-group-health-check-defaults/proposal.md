## Why

When users create a `tencentcloud_gwlb_target_group` resource with a `health_check` block but without explicitly setting `timeout`, `interval_time`, `health_num`, or `un_health_num`, Terraform sends the Go int zero value (`0`) for these fields to the GWLB API. The API rejects these values because they fall outside the allowed ranges (e.g., `timeout` must be in [2, 30]). This causes resource creation to fail with errors like `HealthCheck.Timeout is out of range: [2,30]`, even though the API documentation clearly states the defaults (timeout=2, interval_time=5, health_num=3, un_health_num=3).

## What Changes

- Modify `resourceTencentCloudGwlbTargetGroupCreate` and `resourceTencentCloudGwlbTargetGroupUpdate` so that `timeout`, `interval_time`, `health_num`, and `un_health_num` are read via `d.GetOk("health_check.0.<field>")` instead of from the flattened `health_check` map
- When any of these four fields is omitted, the provider no longer sends `0` to the GWLB API; the field pointer stays `nil` and is omitted from the request (SDK `omitempty`), letting the API apply its own defaults
- The schema is unchanged: the four fields remain `Optional: true, Computed: true`

## Capabilities

### New Capabilities
<!-- No new capabilities are introduced; this is a bug fix for existing capability. -->

### Modified Capabilities
<!-- No spec-level requirement changes; this is a bug fix that aligns behavior with API documentation. -->

## Impact

- **Code**: `tencentcloud/services/gwlb/resource_tc_gwlb_target_group.go` — `Create`/`Update` health check field handling
- **Tests**: `tencentcloud/services/gwlb/resource_tc_gwlb_target_group_test.go` — test that omits the four fields and verifies API defaults are reflected in state
- **Docs**: no doc change needed (schema fields are unchanged)
