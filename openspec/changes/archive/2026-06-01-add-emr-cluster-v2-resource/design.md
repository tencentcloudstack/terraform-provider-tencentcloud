# Design: tencentcloud_emr_cluster_v2 Resource

## Context

Tencent Cloud EMR exposes two families of cluster-creation APIs under version `2019-01-03`:

1. **Legacy** `CreateInstance` (589/34261) — backs the existing `tencentcloud_emr_cluster` resource. Fields are flat and the resource has accumulated many product-specific special cases over the years.
2. **New** `CreateCluster` (589/83953) — the currently recommended path. Parameters are grouped into well-defined nested structs (`SceneSoftwareConfig`, `ZoneResourceConfiguration`, `CustomMetaDBInfo`, `DependService`, `NodeMark`, …), and new platform capabilities (multi-AZ, scene-based deployment, TF node marks, StarRocks `CosBucket`) are only surfaced through this API.

Modifying `tencentcloud_emr_cluster` in place would break backward compatibility (different required fields, different nesting). We therefore introduce a **new** resource `tencentcloud_emr_cluster_v2` that maps 1:1 to the new API.

The CRUD API surface for this first iteration:

| Op | API | Notes |
|---|---|---|
| Create | `CreateCluster` | Async. Returns `InstanceId`. |
| Read (cluster) | `DescribeInstances` | Filter by `InstanceIds=[InstanceId]`, read `Clusters[0]`. `Status == 2` means **running** (creation complete). |
| Read (nodes) | `DescribeClusterNodes` | `NodeFlag="all"`, paginated. Used to populate read-back node resource specs. |
| Delete | `TerminateInstance` | Async. Poll `DescribeInstances` until the cluster disappears or moves to a terminated status. |
| Update | _not in this change_ | All schema fields are declared without `ForceNew`; update handler is a no-op. Modify APIs (scale/tag/rename) will be added in a follow-up change per requirements. |

Project conventions to follow:
- Code layout mirrors `resource_tc_igtm_strategy.go` (CRUD handler file, thin service layer, request builder pattern).
- Resource ID = `InstanceId` (single token, no composite). Read/Delete parse `d.Id()` directly without `strings.Split`.
- `defer tccommon.LogElapsed(...)` and `defer tccommon.InconsistentCheck(...)` at the top of every CRUD function.
- `resource.Retry(...)` with `tccommon.WriteRetryTimeout` / `ReadRetryTimeout` around SDK calls.
- `schema.Timeouts{Create: 60 min, Read: 20 min, Delete: 30 min}` because all operations are long-running.

## Goals / Non-Goals

**Goals:**
- Expose every top-level `CreateCluster` request parameter as a Terraform schema field with a stable 1:1 mapping (no merged or renamed fields).
- Preserve all nested SDK struct boundaries (`LoginSettings`, `SceneSoftwareConfig`, `InstanceChargePrepaid`, `ScriptBootstrapActionConfig`, `Tag`, `CustomMetaDBInfo`, `DependService`, `ZoneResourceConfiguration` with full `AllNodeResourceSpec` → `NodeResourceSpec` → `DiskSpecInfo` tree, `NodeMark`) as nested Terraform blocks.
- Handle async Create by polling `DescribeInstances` until `Clusters[0].Status == 2`.
- Handle async Delete by polling `DescribeInstances` until the cluster is absent (empty `Clusters`) or moved to a terminated status.
- Code style matches `tencentcloud_igtm_strategy` (builder pattern inside Create, nil-safe read loops, `helper.String`/`helper.IntInt64`/`helper.BoolToPtr` etc.).
- Generate `.md` example file (naming convention `resource_tc_emr_cluster_v2.md`, mirroring `resource_tc_config_compliance_pack.md`).
- Generate acceptance test file (naming convention `resource_tc_emr_cluster_v2_test.go`).

**Non-Goals:**
- **No Update implementation** in this change. Schema fields MUST NOT carry `ForceNew: true` so that later updates can be added without forcing recreation.
- No companion data source (the existing `data_source_tc_emr` may be used; adding a v2 data source is out of scope).
- No sweeper registration modifications beyond what already exists for EMR.
- No changes to the existing `tencentcloud_emr_cluster` resource.

## Decisions

### D1. Resource ID is the plain `InstanceId`

`CreateCluster` returns `InstanceId` (e.g. `emr-f2da1cd`). This single token is a globally unique cluster identifier that is also what `DescribeInstances`, `DescribeClusterNodes`, and `TerminateInstance` all consume.

- ID set via `d.SetId(*response.Response.InstanceId)`.
- Read/Delete use `d.Id()` directly — no `FILED_SP` split needed.
- Import passthrough works out of the box.

**Alternative considered:** composite `region#instanceId`. Rejected — region is provider-level, not resource-level.

### D2. Schema layout strictly mirrors `CreateClusterRequestParams`

Every top-level field in `CreateClusterRequestParams` becomes a top-level schema attribute. Every nested SDK struct becomes a nested `TypeList` with `MaxItems: 1` (for singular structs) or without MaxItems (for arrays):

| SDK field | Terraform attribute | Schema type |
|---|---|---|
| `ProductVersion` | `product_version` | `TypeString, Required` |
| `EnableSupportHAFlag` | `enable_support_ha_flag` | `TypeBool, Required` |
| `InstanceName` | `instance_name` | `TypeString, Required` |
| `InstanceChargeType` | `instance_charge_type` | `TypeString, Required` |
| `LoginSettings` | `login_settings` | `TypeList, MaxItems:1, Required` — nested: `password` (sensitive), `public_key_id` |
| `SceneSoftwareConfig` | `scene_software_config` | `TypeList, MaxItems:1, Required` — nested: `software []string`, `scene_name` |
| `InstanceChargePrepaid` | `instance_charge_prepaid` | `TypeList, MaxItems:1, Optional` — nested: `period`, `renew_flag` |
| `SecurityGroupIds` | `security_group_ids` | `TypeList[String], Optional` |
| `ScriptBootstrapActionConfig` | `script_bootstrap_action_config` | `TypeList, Optional` — nested: `cos_file_uri`, `execution_moment`, `args []string`, `cos_file_name`, `remark` |
| `ClientToken` | `client_token` | `TypeString, Optional` |
| `NeedMasterWan` | `need_master_wan` | `TypeString, Optional` |
| `EnableRemoteLoginFlag` | `enable_remote_login_flag` | `TypeBool, Optional` |
| `EnableKerberosFlag` | `enable_kerberos_flag` | `TypeBool, Optional` |
| `CustomConf` | `custom_conf` | `TypeString, Optional` |
| `Tags` | `tags` | `TypeList, Optional` — nested: `tag_key`, `tag_value` |
| `DisasterRecoverGroupIds` | `disaster_recover_group_ids` | `TypeList[String], Optional` |
| `EnableCbsEncryptFlag` | `enable_cbs_encrypt_flag` | `TypeBool, Optional` |
| `MetaDBInfo` | `meta_db_info` | `TypeList, MaxItems:1, Optional` — nested: `meta_data_jdbc_url`, `meta_data_user`, `meta_data_pass` (sensitive), `meta_type`, `unify_meta_instance_id` |
| `DependService` | `depend_service` | `TypeList, Optional` — nested: `service_name`, `instance_id` |
| `ZoneResourceConfiguration` | `zone_resource_configuration` | `TypeList, Optional` — nested tree below |
| `CosBucket` | `cos_bucket` | `TypeString, Optional` |
| `NodeMarks` | `node_marks` | `TypeList, Optional` — nested: `node_type`, `node_names []string`, `zone` |
| `LoadBalancerId` | `load_balancer_id` | `TypeString, Optional` |
| `DefaultMetaVersion` | `default_meta_version` | `TypeString, Optional` |
| `NeedCdbAudit` | `need_cdb_audit` | `TypeInt, Optional` |
| `SgIP` | `sg_ip` | `TypeString, Optional` |
| `PartitionNumber` | `partition_number` | `TypeInt, Optional` |
| `WebUiVersion` | `web_ui_version` | `TypeInt, Optional` |

`zone_resource_configuration` expanded:

```
zone_resource_configuration (TypeList)
 ├─ virtual_private_cloud (TypeList, MaxItems:1) { vpc_id, subnet_id }
 ├─ placement            (TypeList, MaxItems:1) { zone, project_id }
 ├─ zone_tag             (TypeString)
 └─ all_node_resource_spec (TypeList, MaxItems:1)
     ├─ master_count, core_count, task_count, common_count (TypeInt)
     ├─ master_resource_spec (TypeList, MaxItems:1) → NodeResourceSpec
     ├─ core_resource_spec   (TypeList, MaxItems:1) → NodeResourceSpec
     ├─ task_resource_spec   (TypeList, MaxItems:1) → NodeResourceSpec
     └─ common_resource_spec (TypeList, MaxItems:1) → NodeResourceSpec

NodeResourceSpec
 ├─ instance_type (TypeString)
 ├─ system_disk    (TypeList) → DiskSpecInfo { count, disk_size, disk_type }
 ├─ data_disk      (TypeList) → DiskSpecInfo
 ├─ local_data_disk(TypeList) → DiskSpecInfo
 └─ tags           (TypeList) → { tag_key, tag_value }
```

**Rationale:** strict 1:1 mapping maximises predictability, simplifies the later Update implementation, and matches the `igtm_strategy` reference style. Any cross-field validation is deferred to the cloud API (which already validates) — this keeps the resource code mechanical and easy to review.

**Alternative considered:** flatten `login_settings` into top-level `password` / `public_key_id`. Rejected — diverges from the SDK shape and makes future Update mapping harder.

### D3. No `ForceNew` in this change

Per requirements, every field is `Optional` or `Required` without `ForceNew`. This keeps the door open for an Update lifecycle handler in a follow-up change (scale in/out, tag modification, password reset, etc.) without introducing a Terraform state-breaking migration.

**Trade-off:** `terraform plan` on an unchanged resource will never show "force replacement" for currently-immutable fields (e.g. `product_version`, `availability_zone`). In the current change the Update handler is a no-op that returns an explicit error-for-unsupported-field pattern is **not** added (keeps code minimal). Once the Update change lands, those fields that must be immutable will either be handled in Update (returning a friendly error) or marked `ForceNew` at that time.

### D4. Create waits for `Status == 2`

`CreateCluster` is asynchronous. The resource MUST NOT return from Create until the cluster is in state `Running`.

Implementation:

```go
err := resource.Retry(d.Timeout(schema.TimeoutCreate) - time.Minute, func() *resource.RetryError {
    cluster, e := service.DescribeEmrClusterV2ById(ctx, instanceId)
    if e != nil {
        return tccommon.RetryError(e)
    }
    if cluster == nil {
        return resource.RetryableError(fmt.Errorf("cluster %s not yet visible", instanceId))
    }
    switch *cluster.Status {
    case 2:
        return nil // running
    case 3, 4, 5, 6: // terminating / terminated / creating-failed etc.
        return resource.NonRetryableError(fmt.Errorf("cluster %s entered terminal status %d", instanceId, *cluster.Status))
    default: // 0 creating, 1 pending, ...
        return resource.RetryableError(fmt.Errorf("cluster %s status=%d, waiting", instanceId, *cluster.Status))
    }
})
```

Create timeout default: **60 minutes** (EMR cluster provisioning commonly takes 15–30 min).

### D5. Read splits into cluster-level + node-level calls

- `DescribeInstances` populates cluster-level fields (`product_version`, `instance_name`, HA flag, tags, charge type, master-wan, Kerberos, status, etc.). `DescribeInstances` response field names in `ClusterInstancesInfo` do not 1:1 match the request; read-back sets only fields that are reliably reported by the API. Fields that `DescribeInstances` does not return (e.g. `ClientToken`, `CustomConf`, secrets in `LoginSettings.Password`) are preserved from state (by **not** calling `d.Set` on them — they stay as the user declared).
- `DescribeClusterNodes` (paginated, `Limit=100`, `Offset` increment) populates `zone_resource_configuration` node specs. Iterate all pages and aggregate by `NodeType`. Best-effort mapping — if the API omits some specs, leave those sub-fields unset.

### D6. Delete polls for absence

`TerminateInstance` is also async. After calling the API, poll `DescribeInstances` until the returned `Clusters` slice is empty (or Status enters a terminated value). Delete timeout: **30 minutes**.

### D7. Password/secrets handling

`login_settings.password` and `meta_db_info.meta_data_pass` are marked `Sensitive: true`. Read never overwrites them from the API (they are not returned). They remain in state as the user last supplied.

### D8. File & symbol naming

- Resource file: `tencentcloud/services/emr/resource_tc_emr_cluster_v2.go`
- Example MD: `tencentcloud/services/emr/resource_tc_emr_cluster_v2.md`
- Test file: `tencentcloud/services/emr/resource_tc_emr_cluster_v2_test.go`
- Exported symbol: `ResourceTencentCloudEmrClusterV2`
- Service helpers (appended to `service_tencentcloud_emr.go`):
  - `(me *EMRService) DescribeEmrClusterV2ById(ctx, instanceId) (*emr.ClusterInstancesInfo, error)`
  - `(me *EMRService) DescribeEmrClusterV2Nodes(ctx, instanceId, nodeFlag) ([]*emr.NodeHardwareInfo, error)`

### D9. Provider registration

Register `tencentcloud_emr_cluster_v2` in `tencentcloud/provider.go` ResourcesMap, immediately after `tencentcloud_emr_cluster` to keep related resources grouped.

## Risks / Trade-offs

- **Read-back fidelity** → `DescribeInstances` does not round-trip every Create field. We accept per-field inconsistency rather than fake-setting state. State will reflect config for non-readable fields; config drift is detected only for fields the API actually returns. Mitigation: document the limitation in the `.md` example file.
- **Async timeouts in CI** → Acceptance tests can exceed CI budgets. Mitigation: mark the test with a minimal cluster spec (smallest instance type, single AZ) and rely on provider-level `TF_ACC=1` gating so it only runs when explicitly requested. Unit test focuses on schema shape.
- **Partial creation failures** → If `CreateCluster` returns `InstanceId` but the cluster later enters a failed status, the resource is left in state with `SetId(instanceId)`. Create returns a non-retryable error. Users must `terraform destroy` to clean up. This mirrors what `tencentcloud_emr_cluster` already does.
- **Future Update breaking change** → No-op Update today may silently ignore user changes. Mitigation: document in `.md` example that modifications currently require manual recreate until the v2-update change ships.

## Migration Plan

None — purely additive. Users continue to use `tencentcloud_emr_cluster` for existing state; new clusters may opt into `tencentcloud_emr_cluster_v2`. No state migration, no deprecation in this change.

## Open Questions

- Q1 (resolved): Whether `Status` semantics are `2 = Running`? → **Yes**, confirmed by existing `resource_tc_emr_cluster.go` wait logic and EMR docs 589/34266.
- Q2: Should Create surface `request_id` as computed? → Out of scope; not needed for reconciliation.
