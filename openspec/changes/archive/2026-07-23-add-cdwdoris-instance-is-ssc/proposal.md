## Why

The CDW Doris (cdwdoris) cloud API `CreateInstanceNew` supports the `IsSSC` parameter to control whether the instance uses storage-compute separation (存算分离) architecture. This parameter is currently not exposed in the Terraform resource `tencentcloud_cdwdoris_instance`, preventing users from creating storage-compute separation instances via Terraform. Adding this parameter enables users to leverage the storage-compute separation feature through Infrastructure as Code.

## What Changes

- Add a new optional `is_ssc` parameter (TypeBool, Optional, ForceNew) to the `tencentcloud_cdwdoris_instance` resource
- Pass the parameter to the `CreateInstanceNew` API during resource creation
- The parameter is write-only (not read back from the `DescribeInstance` API response as `InstanceInfo` does not include this field)

## Capabilities

### New Capabilities
- `cdwdoris-instance-is-ssc`: Add write-only `is_ssc` parameter to `tencentcloud_cdwdoris_instance` resource for storage-compute separation configuration

### Modified Capabilities
<!-- None - no existing specs are modified -->

## Impact

- **Affected code**: `tencentcloud/services/cdwdoris/resource_tc_cdwdoris_instance.go`
- **Affected tests**: `tencentcloud/services/cdwdoris/resource_tc_cdwdoris_instance_test.go`
- **Affected docs**: `tencentcloud/services/cdwdoris/resource_tc_cdwdoris_instance.md`
- **SDK dependency**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris/v20211228` (already in vendor, no update needed)
- **Backward compatibility**: Fully backward compatible — the new parameter is Optional, existing configurations continue to work unchanged