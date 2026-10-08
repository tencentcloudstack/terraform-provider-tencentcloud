## Why

The CDB `tencentcloud_mysql_instance` resource does not expose the `DiskEncryption` parameter, even though the cloud API `CreateDBInstance` / `CreateDBInstanceHour` accept it as a request field and `DescribeDBInstances` returns it in the response (`InstanceInfo.DiskEncryption`). Users who need to enable disk encryption for cloud-disk edition CDB instances (a compliance/security requirement) cannot do so through Terraform, forcing them to use the console or API directly.

## What Changes

- Add `disk_encryption` (Optional, ForceNew, Computed, TypeString) parameter to the `tencentcloud_mysql_instance` resource schema. The value `on` means disk encryption is enabled; otherwise disk encryption is disabled. Only cloud-disk edition instances support this feature.
- Wire `disk_encryption` into the `CreateDBInstance` request (`request.DiskEncryption`) in the month-paid create path (`mysqlCreateInstancePayByMonth`).
- Wire `disk_encryption` into the `CreateDBInstanceHour` request (`request.DiskEncryption`) in the hourly-paid create path (`mysqlCreateInstancePayByUse`).
- Reuse the shared `mysqlAllInstanceRoleSet` helper (which handles both `CreateDBInstanceRequest` and `CreateDBInstanceHourRequest`) so the parameter is set consistently for all instance roles on both billing models.
- Read `disk_encryption` from the `DescribeDBInstances` response (`InstanceInfo.DiskEncryption`) in `tencentMsyqlBasicInfoRead` and set it into Terraform state (guard against nil).
- Mark `disk_encryption` as `ForceNew: true` because no Modify/Upgrade CDB API accepts `DiskEncryption`; it can only be set at creation time.
- Do NOT add `FourthZone` to the resource schema — although `CreateDBInstance`/`CreateDBInstanceHour` accept it, `DescribeDBInstances` does not return it, so there is no Read path to support it. This change is intentionally scoped to a single new parameter (`disk_encryption`) only.
- Update the resource example documentation `tencentcloud/services/cdb/resource_tc_mysql_instance.md` to add a `disk_encryption` usage example.

## Capabilities

### New Capabilities
- `mysql-instance-disk-encryption`: Enable the `disk_encryption` parameter on the `tencentcloud_mysql_instance` resource so users can enable disk encryption when creating cloud-disk edition CDB instances, and refresh the encryption status from the Describe API.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/cdb/resource_tc_mysql_instance.go` — add `disk_encryption` schema field; set it in the shared create helper (`mysqlAllInstanceRoleSet`) so both `CreateDBInstance` and `CreateDBInstanceHour` requests carry it; read it from the `DescribeDBInstances` response in `tencentMsyqlBasicInfoRead` (nil guard); the `Update` function requires no change because the field is `ForceNew`.
  - `tencentcloud/services/cdb/resource_tc_mysql_instance.md` — add `disk_encryption` usage example.
- **SDK dependency:** The vendored `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320` already includes `DiskEncryption` on `CreateDBInstanceRequest`, `CreateDBInstanceHourRequest`, and `InstanceInfo` (DescribeDBInstances response). No SDK upgrade is required.
- **Backward compatibility:** Fully backward compatible — the new parameter is `Optional` and `Computed`, so existing configurations and state are unaffected.
- **API constraints:** `DiskEncryption` is only accepted by `CreateDBInstance` / `CreateDBInstanceHour` (not by any Modify/Update API), therefore the parameter is immutable after creation (`ForceNew`). It is returned by `DescribeDBInstances`, so Read can refresh the value. Only cloud-disk edition instances support disk encryption (per API documentation); the API itself will validate this constraint.