## Why

When users create a `tencentcloud_gwlb_target_group` resource with a `health_check` block but without explicitly setting `timeout`, `interval_time`, `health_num`, or `un_health_num`, Terraform sends the Go int zero value (`0`) for these fields to the GWLB API. The API rejects these values because they fall outside the allowed ranges (e.g., `timeout` must be in [2, 30]). This causes resource creation to fail with errors like `HealthCheck.Timeout is out of range: [2,30]`, even though the API documentation clearly states the defaults (timeout=2, interval_time=5, health_num=3, un_health_num=3).

## What Changes

- Add `Default: 2` to the `timeout` field in the `health_check` nested schema
- Add `Default: 5` to the `interval_time` field in the `health_check` nested schema
- Add `Default: 3` to the `health_num` field in the `health_check` nested schema
- Add `Default: 3` to the `un_health_num` field in the `health_check` nested schema

## Capabilities

### New Capabilities
<!-- No new capabilities are introduced; this is a bug fix for existing capability. -->

### Modified Capabilities
<!-- No spec-level requirement changes; this is a bug fix that aligns behavior with API documentation. -->

## Impact

- **Code**: `tencentcloud/services/gwlb/resource_tc_gwlb_target_group.go` — schema definition for `health_check` nested block
- **Tests**: `tencentcloud/services/gwlb/resource_tc_gwlb_target_group_test.go` — update test cases to verify default values
- **Docs**: `tencentcloud/services/gwlb/resource_tc_gwlb_target_group.md` — no change needed, `Default` values are auto-documented