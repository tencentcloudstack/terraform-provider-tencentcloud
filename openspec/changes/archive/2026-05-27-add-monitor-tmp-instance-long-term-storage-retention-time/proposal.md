## Why

The `tencentcloud_monitor_tmp_instance` resource currently does not support configuring the long-term storage retention time (`LongTermStorageRetentionTime`) via `InstanceAttributes`. The cloud APIs `CreatePrometheusMultiTenantInstancePostPayMode` and `ModifyPrometheusInstanceAttributes` both support the `InstanceAttributes` parameter with key `LongTermStorageRetentionTime` (value: 60-730 days), and the `DescribePrometheusInstances` response returns the `InstanceAttributes` array including this key. Adding this parameter allows users to configure archive storage duration for Prometheus instances via Terraform.

## What Changes

- Add a new optional parameter `long_term_storage_retention_time` (type: int, range: 60-730) to the `tencentcloud_monitor_tmp_instance` resource schema.
- Wire the parameter to the `CreatePrometheusMultiTenantInstancePostPayMode` API's `InstanceAttributes` field (key: `LongTermStorageRetentionTime`) during resource creation.
- Wire the parameter to the `ModifyPrometheusInstanceAttributes` API's `InstanceAttributes` field during resource update.
- Read back the `LongTermStorageRetentionTime` from the `DescribePrometheusInstances` API response's `InstanceAttributes` array in the Read function.

## Capabilities

### New Capabilities
- `monitor-tmp-instance-long-term-storage-retention-time`: Add `long_term_storage_retention_time` parameter to `tencentcloud_monitor_tmp_instance` resource for configuring archive storage retention duration (60-730 days).

### Modified Capabilities

## Impact

- Resource file: `tencentcloud/services/tmp/resource_tc_monitor_tmp_instance.go` — schema addition, Create/Read/Update logic changes.
- Test file: `tencentcloud/services/tmp/resource_tc_monitor_tmp_instance_test.go` — updated unit test cases.
- Documentation: `tencentcloud/services/tmp/resource_tc_monitor_tmp_instance.md` — updated parameter documentation.
