## Context

The `tencentcloud_mysql_instance` resource manages CDB MySQL instances. The resource uses a dual-path creation pattern: `mysqlCreateInstancePayByMonth` uses `CreateDBInstance` (prepaid), and `mysqlCreateInstancePayByUse` uses `CreateDBInstanceHour` (postpaid). Both paths share parameter-setting logic in `mysqlAllInstanceRoleSet` and `mysqlMasterInstanceRoleSet`, which handle a shared `requestInter interface{}` that is type-asserted to either `*cdb.CreateDBInstanceRequest` or `*cdb.CreateDBInstanceHourRequest`.

**Current state:**
- Resource file: `tencentcloud/services/cdb/resource_tc_mysql_instance.go`
- Schema is split across `TencentMsyqlBasicInfo()` and `ResourceTencentCloudMysqlInstance()`
- Create logic: `mysqlAllInstanceRoleSet` sets common params; `mysqlMasterInstanceRoleSet` sets master-specific params; both handle dual-path via `okByMonth` boolean
- Read logic: `tencentMsyqlBasicInfoRead` fetches `*cdb.InstanceInfo` via `DescribeDBInstanceById`, sets state fields
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320`

**API behavior analysis:**

| API | DestroyProtect in Request | DestroyProtect in Response |
|-----|--------------------------|----------------------------|
| `CreateDBInstance` | Yes (`DestroyProtect *string`, `on`/`off`) | N/A |
| `CreateDBInstanceHour` | Yes (`DestroyProtect *string`, `on`/`off`) | N/A |
| `DescribeDBInstances` | N/A | Yes (`InstanceInfo.DestroyProtect *string`, `on` or other) |

**Key constraint:** `DestroyProtect` is a string parameter with values `on` (enable) and `off` (disable). The `DescribeDBInstances` response returns the current status in `InstanceInfo.DestroyProtect`.

## Goals / Non-Goals

**Goals:**
- Add `destroy_protect` (Optional, Computed, TypeString) parameter to `tencentcloud_mysql_instance` with valid values `on` and `off`
- Pass `DestroyProtect` to both `CreateDBInstance` and `CreateDBInstanceHour` API requests when specified by user
- Read `DestroyProtect` from `DescribeDBInstances` API response (`InstanceInfo.DestroyProtect`) to support state refresh and import
- Modify `DestroyProtect` via `ModifyInstanceDestroyProtect` API during Update when `destroy_protect` changes
- Maintain full backward compatibility — existing configurations continue to work unchanged

**Non-Goals:**
- Adding `destroy_protect` to the `tencentcloud_mysql_instance` datasource (out of scope)

## Decisions

### Decision 1: Add `destroy_protect` to `TencentMsyqlBasicInfo()` schema map

**Rationale:** The `destroy_protect` parameter is a basic instance attribute applicable to all instance roles (master, readonly, dr), so it belongs in `TencentMsyqlBasicInfo()` alongside other common fields like `disk_type`, `device_type`. This is consistent with how similar create-only parameters are organized.

### Decision 2: Set `DestroyProtect` in `mysqlAllInstanceRoleSet` using the dual-path pattern

**Rationale:** The `mysqlAllInstanceRoleSet` function already handles dual-path parameter setting via `requestInter interface{}` with type assertion to `*cdb.CreateDBInstanceRequest` (okByMonth) and `*cdb.CreateDBInstanceHourRequest`. Adding `DestroyProtect` here follows the existing pattern used by `disk_type`, `device_type`, etc. The code reads the value with `d.GetOk("destroy_protect")` and sets it on whichever request struct is active.

### Decision 3: Read `DestroyProtect` in `tencentMsyqlBasicInfoRead` with nil check

**Rationale:** The `tencentMsyqlBasicInfoRead` function already sets state from `mysqlInfo *cdb.InstanceInfo`. Adding a nil-checked read of `mysqlInfo.DestroyProtect` follows the same pattern used for `mysqlInfo.DiskType` (line 967-969). The nil check is required because the API may return `nil` for the pointer field.

### Decision 4: `destroy_protect` is NOT ForceNew and NOT in immutableArgs

**Rationale:** The parameter is a standard Optional+Computed field, so it does not trigger recreation when changed. When a user changes `destroy_protect` after creation, the Update flow calls the `ModifyInstanceDestroyProtect` API via the `mysqlAllInstanceRoleUpdate` function to toggle the protection status without recreating the instance. The field is therefore NOT added to the `immutableArgs` list.

## Risks / Trade-offs

- **[Risk] API returns empty string instead of nil for `DestroyProtect`**: The `DescribeDBInstances` response may return an empty string for instances where destroy protection was never set.
  - **Mitigation:** Use `d.Set("destroy_protect", ...)` which handles empty strings gracefully; only skip setting if the pointer itself is nil.

- **[Risk] `ModifyInstanceDestroyProtect` API may fail on transitional instance states**: The API may reject the request when the instance is in a transitional state (e.g., creating, isolating).
  - **Mitigation:** The service method returns the raw API error to the user, who can retry the apply. This is consistent with other update operations in the resource.
