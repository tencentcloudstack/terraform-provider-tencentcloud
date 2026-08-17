## Why

The CDB `CreateDBInstance` and `CreateDBInstanceHour` APIs support a `DestroyProtect` parameter to enable or disable instance destroy protection, and the `DescribeDBInstances` API returns the current `DestroyProtect` status in its response. However, the Terraform resource `tencentcloud_mysql_instance` does not expose this parameter. Users cannot configure or view the destroy protection setting through Terraform, forcing them to use the console or API directly.

## What Changes

- Add `destroy_protect` (Optional, Computed) parameter to `tencentcloud_mysql_instance` resource to support specifying the destroy protection status during instance creation. Valid values: `on` (enable destroy protection), `off` (disable destroy protection). The parameter is read back from the `DescribeDBInstances` API response (`InstanceInfo.DestroyProtect`).
- Pass `DestroyProtect` to both `CreateDBInstance` and `CreateDBInstanceHour` API requests when the user specifies `destroy_protect` in the configuration.
- Read `DestroyProtect` from the `DescribeDBInstances` API response (`InstanceInfo.DestroyProtect`) in the Read function to support state refresh and import.

## Capabilities

### New Capabilities
- `mysql-instance-destroy-protect`: Enable the `destroy_protect` parameter on the `tencentcloud_mysql_instance` resource to allow users to specify the destroy protection status when creating MySQL instances.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/cdb/resource_tc_mysql_instance.go` — add `destroy_protect` schema field, wire through Create flow (both `mysqlAllInstanceRoleSet` and `mysqlMasterInstanceRoleSet` paths), add Read support in `tencentMsyqlBasicInfoRead`
  - `tencentcloud/services/cdb/resource_tc_mysql_instance_test.go` — add test case for `destroy_protect` parameter
  - `tencentcloud/services/cdb/resource_tc_mysql_instance.md` — update documentation with usage example
- **SDK dependency:** `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320` — `CreateDBInstanceRequest`, `CreateDBInstanceHourRequest`, and `InstanceInfo` structs already include `DestroyProtect` field
- **Backward compatibility:** fully backward compatible — the new parameter is Optional and Computed; existing configurations continue to work unchanged
- **API constraints:** `DestroyProtect` is accepted by both `CreateDBInstance` and `CreateDBInstanceHour` create APIs. The `DescribeDBInstances` response includes `DestroyProtect` in the `InstanceInfo` struct, so Read can refresh this value. The value is a string (`on`/`off`).
