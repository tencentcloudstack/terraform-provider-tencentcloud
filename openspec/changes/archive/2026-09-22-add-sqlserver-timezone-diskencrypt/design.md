## Context

The `tencentcloud_sqlserver_basic_instance` resource manages SQL Server basic (单实例) instances via the TencentCloud SQL Server API. Currently the resource schema exposes core instance parameters (cpu, memory, storage, vpc, subnet, engine_version, etc.) but does NOT expose two create-time configuration knobs that the underlying cloud API `CreateBasicDBInstances` supports:

- `TimeZone` — system timezone string (default "China Standard Time")
- `DiskEncryptFlag` — disk encryption toggle (0=off, 1=on)

Both are accepted by the create API but neither is surfaced in Terraform, so users cannot provision instances with a custom timezone or with disk encryption enabled, nor can they see these values in state.

### Current state of implementation

The code change has already been applied to the working tree (resource + service files). This design document captures the vendor SDK field definitions and existing schema structure so the implementation can be verified without re-reading the large `models.go` file.

### Relevant cloud API SDK definitions (vendor)

Package: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328`
File: `vendor/.../sqlserver/v20180328/models.go`

**`CreateBasicDBInstancesRequest`** (models.go ~line 1019) — the two new request fields:

```go
// <p>系统时区，默认：China Standard Time</p>
TimeZone *string `json:"TimeZone,omitnil,omitempty" name:"TimeZone"`

// <p>磁盘加密标识，0-不加密，1-加密</p>
DiskEncryptFlag *int64 `json:"DiskEncryptFlag,omitnil,omitempty" name:"DiskEncryptFlag"`
```

Both fields are optional pointer fields on the request struct (`*string` and `*int64` respectively). They appear near the end of the request struct, after `Collation` and before `ThroughputPerformance`.

**`DBInstance`** (models.go ~line 2876) — the `DescribeDBInstances` list-element struct, contains `TimeZone` near the end (~line 3026):

```go
// <p>系统时区，默认：China Standard Time</p>
TimeZone *string `json:"TimeZone,omitnil,omitempty" name:"TimeZone"`
```

Note: `DBInstance` does NOT contain a disk-encryption field. The disk encryption status must be fetched separately.

**`DescribeDBInstancesAttributeResponseParams`** (models.go ~line 5584) — response of `DescribeDBInstancesAttribute`, contains `IsDiskEncryptFlag` (~line 5628):

```go
// 是否开启磁盘加密，1-开启，0-未开启
IsDiskEncryptFlag *int64 `json:"IsDiskEncryptFlag,omitnil,omitempty" name:"IsDiskEncryptFlag"`
```

This is a single-instance attribute query (request takes `InstanceId`), returns `IsDiskEncryptFlag *int64`.

### Existing resource schema structure (resource_tc_sqlserver_basic_instance.go)

The new fields are inserted into the `Schema` map of `ResourceTencentCloudSqlserverBasicInstance()`, immediately after the existing `engine_version` block and before `period`:

```go
"time_zone": {
    Type:        schema.TypeString,
    Optional:    true,
    Computed:    true,
    ForceNew:    true,
    Description: "System timezone for the SQL Server instance. Default is `China Standard Time`. This setting cannot be changed after creation.",
},
"disk_encrypt_flag": {
    Type:         schema.TypeInt,
    Optional:     true,
    Computed:     true,
    ForceNew:     true,
    ValidateFunc: tccommon.ValidateIntegerInRange(0, 1),
    Description:  "Disk encryption flag. `0` - Disabled (default), `1` - Enabled. Disk encryption cannot be changed after instance creation.",
},
```

The existing `collation` field (also ForceNew-via-immutableArgs, Optional with a Default) is the closest analog and sits later in the schema map.

## Goals / Non-Goals

**Goals:**
- Expose `time_zone` (string) and `disk_encrypt_flag` (int, 0|1) as optional, computed, ForceNew schema fields on `tencentcloud_sqlserver_basic_instance`.
- Pass both parameters into `CreateBasicDBInstances` during create.
- Read `time_zone` back from the existing `DescribeDBInstances` response (`DBInstance.TimeZone`).
- Read `disk_encrypt_flag` back via a new `DescribeDBInstancesAttribute` call (`IsDiskEncryptFlag`), since the list API does not return it.
- Enforce ForceNew by adding both fields to the `immutableArgs` list in the update function (the resource has no dedicated update API for these attributes).
- Keep the change fully backward compatible: both fields are Optional+Computed, so existing configs and state continue to work.

**Non-Goals:**
- No support for changing timezone or disk encryption post-creation (the cloud API does not offer update endpoints for these; they are inherently create-only, hence ForceNew).
- No modification of the `tencentcloud_sqlserver_basic_instances` data source (separate resource, out of scope).
- No new top-level schema block / no nested list wrapper around the instance attributes.
- No changes to provider registration (the resource already exists in `provider.go`/`provider.md`).

## Decisions

### Decision 1: Both fields are Optional + Computed + ForceNew

**Rationale.** The cloud API treats both as optional create inputs with server-side defaults ("China Standard Time" / 0). Marking them `Computed` means that when a user omits them, the next `read` populates state with whatever the API actually applied, preventing perpetual diffs. `ForceNew` is required because neither attribute has an update path in the API.

**Alternative considered:** `Optional + Default` (not Computed). Rejected — using `Default` for `disk_encrypt_flag=0` conflicts with `Computed` semantics in the SDK (Default and Computed are mutually exclusive in practice), and for `time_zone` hardcoding a default string would mask the real API-applied value. The Computed approach is consistent with the existing `availability_zone` field on this resource.

### Decision 2: `disk_encrypt_flag` uses `GetOkExists` in Create (not `GetOk`)

**Rationale.** `disk_encrypt_flag` is a `TypeInt` with valid value `0`. `GetOk` returns `false` for a zero value, which would silently drop an explicit `disk_encrypt_flag = 0` from the create request. `GetOkExists` correctly distinguishes "field absent" from "field set to 0". (`time_zone` is a string and uses `GetOk` since empty string is not a meaningful value.)

### Decision 3: Read path uses TWO API calls

**Rationale.** `DescribeDBInstances` returns `TimeZone` but not disk encryption. `DescribeDBInstancesAttribute` returns `IsDiskEncryptFlag` but not the full instance record. The read function therefore keeps the existing `DescribeSqlserverInstanceById` call and adds a second call to the new `DescribeSqlserverInstanceAttributeById` service method.

**Failure handling:** If `DescribeDBInstancesAttribute` fails, the read logs a `[WARN]` and continues rather than failing the whole read — this avoids a transient attribute-API error wiping out all other state. `disk_encrypt_flag` simply won't be populated in that pass. This is an intentional graceful-degradation choice.

**Alternative considered:** Failing the read on attribute-API error. Rejected — it would make the resource unreadable whenever the attribute endpoint has a hiccup, even though the core instance is healthy.

### Decision 4: New service method `DescribeSqlserverInstanceAttributeById`

Signature:
```go
func (me *SqlserverService) DescribeSqlserverInstanceAttributeById(ctx context.Context, instanceId string) (
    attribute *sqlserver.DescribeDBInstancesAttributeResponseParams, errRet error,
)
```
It wraps `DescribeDBInstancesAttribute`, sets `request.InstanceId`, applies `ratelimit.Check`, and returns `response.Response` (with a nil-safe extraction). This mirrors the existing `DescribeSqlserverInstanceById` style in the same service file.

### Decision 5: immutableArgs enforcement in Update

The update function already guards `collation` via an `immutableArgs` slice. Both new fields are appended:
```go
immutableArgs := []string{"collation", "time_zone", "disk_encrypt_flag"}
```
Returning an explicit error on `d.HasChange(v)` gives users a clear message instead of Terraform silently ignoring the change.

## Risks / Trade-offs

- **[Extra API call on every read]** → The read now issues `DescribeDBInstancesAttribute` in addition to `DescribeDBInstances`. Mitigation: the call is cheap and single-instance; if it fails the read degrades gracefully (warn + skip the field) rather than failing.
- **[Nil pointer dereference on `IsDiskEncryptFlag`]** → The attribute response field is `*int64` and may be nil. Mitigation: double nil-check `attribute != nil && attribute.IsDiskEncryptFlag != nil` before dereferencing; same pattern used for `instance.TimeZone`.
- **[ForceNew surprises existing users]** → Users who currently rely on the API default timezone/encryption and later upgrade the provider will see `time_zone`/`disk_encrypt_flag` become computed in state with no diff (since they didn't set them). Mitigation: both fields are Optional+Computed, so no plan-time diff is produced for configs that don't reference them.
- **[`disk_encrypt_flag = 0` silently dropped if `GetOk` used]** → Mitigated by using `GetOkExists` in the create path (Decision 2).
