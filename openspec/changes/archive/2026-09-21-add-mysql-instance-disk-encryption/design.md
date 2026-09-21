## Context

The `tencentcloud_mysql_instance` resource (`tencentcloud/services/cdb/resource_tc_mysql_instance.go`) already supports creating CDB instances via two billing paths:
- **Month-paid**: `mysqlCreateInstancePayByMonth` builds a `cdb.CreateDBInstanceRequest`.
- **Hourly-paid**: `mysqlCreateInstancePayByUse` builds a `cdb.CreateDBInstanceHourRequest`.

Both paths funnel common request fields through the shared helper `mysqlAllInstanceRoleSet(ctx, requestInter, d, meta)`, which uses a type switch on `requestInter` (`*cdb.CreateDBInstanceRequest` vs `*cdb.CreateDBInstanceHourRequest`) and assigns each field to the correct request type. Existing parameters like `disk_type`, `device_type`, and `destroy_protect` already follow this pattern.

The Read path (`tencentMsyqlBasicInfoRead`) reads the `DescribeDBInstances` response into a `*cdb.InstanceInfo` and sets fields into state, with nil guards for recently added fields (e.g. `DiskType`, `DestroyProtect`).

The vendored SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320`) already exposes:
- `CreateDBInstanceRequest.DiskEncryption` (*string)
- `CreateDBInstanceHourRequest.DiskEncryption` (*string)
- `InstanceInfo.DiskEncryption` (*string)  — returned by `DescribeDBInstances`

No SDK upgrade is required. No Modify/Update CDB API accepts `DiskEncryption`, so the parameter must be `ForceNew`.

## Goals / Non-Goals

**Goals:**
- Expose `disk_encryption` as an Optional + ForceNew + Computed TypeString on `tencentcloud_mysql_instance`.
- Set it on both create request types via the shared `mysqlAllInstanceRoleSet` helper.
- Read it back from `DescribeDBInstances` in `tencentMsyqlBasicInfoRead` with a nil guard.
- Keep the change fully backward compatible (Optional + Computed).

**Non-Goals:**
- Do NOT add `FourthZone`. `DescribeDBInstances` does not return it, so there is no Read path. The change is intentionally scoped to one new parameter.
- Do NOT support updating `disk_encryption` in place — it is `ForceNew` (no API supports it).
- Do NOT add Terraform-side validation that restricts `disk_encryption` to cloud-disk edition instances only — the cloud API validates this constraint and returns a clear error; duplicating it in the provider risks drift if the API's supported matrix changes.

## Decisions

### Decision 1: Set the field in `mysqlAllInstanceRoleSet` (not in each billing path separately)
The existing code already centralizes cross-billing-model fields in `mysqlAllInstanceRoleSet`. Adding `disk_encryption` there means a single code block handles both `CreateDBInstanceRequest` and `CreateDBInstanceHourRequest`, matching the established pattern for `disk_type` and `destroy_protect`.
- **Alternative considered**: Set the field separately in `mysqlCreateInstancePayByMonth` and `mysqlCreateInstancePayByUse`. **Rejected**: duplicates logic and diverges from the existing pattern.

### Decision 2: `ForceNew: true` (immutable after creation)
Vendor audit confirms `DiskEncryption` appears only on `CreateDBInstanceRequest` and `CreateDBInstanceHourRequest`; no `Modify*`/`Upgrade*`/`Adjust*` request includes it. Therefore it cannot be updated in place. `ForceNew` is the only correct option and aligns with how `disk_type` is handled.
- **Alternative considered**: Try to update via `ModifyDBInstanceConfig`/`UpgradeDBInstance`. **Rejected**: those APIs do not accept `DiskEncryption`.

### Decision 3: Mark `Computed: true`
The `DescribeDBInstances` response returns `DiskEncryption` (e.g. `on` or empty). Marking `Computed` allows imported/pre-existing instances to populate the value, consistent with `disk_type` and `destroy_protect`.

### Decision 4: Use `d.GetOk` and only set the request field when the user provides a value
Consistent with `destroy_protect`, the helper reads via `d.GetOk("disk_encryption")` and only assigns `DiskEncryption` when present. This preserves backward compatibility (existing configs without the field send no `DiskEncryption` and the API applies its default).

### Decision 5: Nil-guard the Read, mirroring `DestroyProtect`
In `tencentMsyqlBasicInfoRead`, set `disk_encryption` only when `mysqlInfo.DiskEncryption != nil`, exactly as done for `destroy_protect` and `disk_type`.

### Decision 6: No Update function change required
Since `disk_encryption` is `ForceNew`, the `Update` function does not need a `HasChange("disk_encryption")` branch. Terraform SDK will plan a replacement when the field changes.

## Risks / Trade-offs

- [Risk] Users enable `disk_encryption` on a non-cloud-disk instance and the API rejects the creation. → Mitigation: The cloud API returns a clear error; the provider surfaces it. No silent failure. Documented in the parameter description that only cloud-disk edition instances support it.
- [Risk] `DescribeDBInstances` returns an empty string for `DiskEncryption` on non-cloud-disk instances, causing diff noise with `Computed`. → Mitigation: `Computed: true` + nil guard means the value is read from the API; an empty-string response is stored as-is and is stable across reads, so no diff noise.
- [Trade-off] `ForceNew` means changing `disk_encryption` recreates the instance. This is acceptable/required because the API does not support changing encryption post-creation.