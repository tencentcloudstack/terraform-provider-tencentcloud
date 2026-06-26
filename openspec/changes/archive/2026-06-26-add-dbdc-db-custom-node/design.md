# Design: tencentcloud_dbdc_db_custom_node Resource

## Architecture

Follows the `tencentcloud_igtm_monitor` style and reuses helpers already added
to the `dbdc` package by the cluster resource:

```
provider.go (register tencentcloud_dbdc_db_custom_node)
    └─ resource_tc_dbdc_db_custom_node.go (Create/Read/Update/Delete)
           └─ service_tencentcloud_dbdc.go
                  ├─ DescribeDBCustomNodeById(ctx, nodeId)          (new)
                  ├─ DescribeDBCustomTaskStatusById(ctx, taskId)    (existing)
                  └─ waitDBCustomTaskSucceeded(...)                 (existing, same package)
                         └─ dbdc SDK v20201029
```

Package `dbdc`, SDK alias `dbdcv20201029`, client
`...GetAPIV3Conn().UseDbdcV20201029Client()`.

## Schema

Arguments map 1:1 to `CreateDBCustomNodesRequest`. `client_token` is intentionally
omitted (idempotency-only, not a meaningful resource attribute — consistent with
the cluster resource decision). Only `tags` and `auto_renew` are mutable via an
API; everything else is `ForceNew`.

### Required (Create, ForceNew)

| Field | Type | Maps to | Description |
|---|---|---|---|
| `zone` | String | `Zone` | Availability zone (e.g. `ap-shanghai-5`) |
| `image_id` | String | `ImageId` | Image ID (`img-xxxxxxx`) |
| `vpc_id` | String | `VpcId` | VPC ID for the node SSH connection |
| `subnet_id` | String | `SubnetId` | Subnet ID for the node SSH connection |
| `node_type` | String | `NodeType` | Node spec (e.g. `DB.AT5.8XLARGE128`) |

### Optional

| Field | Type | ForceNew | Maps to | Description |
|---|---|---|---|---|
| `node_count` | Int (Default 1) | Yes | `NodeCount` | Number of nodes; fixed to 1 (one resource = one node) |
| `period` | Int (Default 1) | Yes | `Period` | Purchase duration in months (1-36) |
| `node_name` | String | Yes | `NodeName` | Node name (<=128 chars; no rename API) |
| `login_settings` | List (MaxItems 1) | Yes | `LoginSettings` | Login config (see below) |
| `auto_voucher` | Int | Yes | `AutoVoucher` | Whether to auto-deduct voucher (0/1) |
| `voucher_ids` | List of String | Yes | `VoucherIds` | Voucher ID list |
| `auto_renew` | Int | No (Computed) | `AutoRenew` | Auto-renew flag (1 renew / 0 not); mutable via RenewDBCustomNode |
| `tags` | Map[string]string | No | `Tags` ([]Tag) | Node tags; mutable via ModifyDBCustomNodeTags |

`login_settings` nested block (ForceNew):

| Field | Type | Maps to | Notes |
|---|---|---|---|
| `password` | String (Sensitive) | `LoginSettings.Password` | Login password |
| `key_ids` | List of String | `LoginSettings.KeyIds` | Key pair IDs (single ID supported) |
| `keep_image_login` | String | `LoginSettings.KeepImageLogin` | `true`/`false`, keep image login settings |

### Computed (read from `DescribeDBCustomNodes` -> `DBCustomNode`)

| Field | Type | Maps to |
|---|---|---|
| `node_id` | String | `NodeId` (same as resource ID) |
| `cluster_id` | String | `ClusterId` |
| `ssh_endpoint` | String | `SSHEndpoint` |
| `lan_ip` | String | `LanIP` |
| `cpu` | Int | `CPU` |
| `memory` | Int | `Memory` |
| `os_name` | String | `OsName` |
| `status` | String | `Status` (`Creating`/`Running`/`Isolating`/`Isolated`/`Activating`/`Destroying`) |
| `charge_type` | String | `ChargeType` |
| `expire_time` | String | `ExpireTime` |
| `created_time` | String | `CreatedTime` |
| `isolated_time` | String | `IsolatedTime` |
| `system_disk` | List (object) | `SystemDisk` {disk_type, disk_size} |
| `data_disks` | List (object) | `DataDisks` {disk_type, disk_size, disk_name} |

> `tags` modeled as Terraform `TypeMap` (idiomatic); converted to/from
> `[]*Tag{Key,Value}`.

## CRUD Logic

### Create
```
request = NewCreateDBCustomNodesRequest()
set Zone, ImageId, VpcId, SubnetId, NodeType (required)
set NodeCount (default 1), Period, NodeName, AutoRenew, AutoVoucher, VoucherIds, Tags
set LoginSettings{Password, KeyIds, KeepImageLogin} if block present
resource.Retry(WriteRetryTimeout): call CreateDBCustomNodes; guard result/Response
if len(NodeIds) == 0 || NodeIds[0] == nil -> error
nodeId = *NodeIds[0]
if TaskId != nil -> waitDBCustomTaskSucceeded(taskId)
d.SetId(nodeId)
return Read
```

### Read
```
respData = DescribeDBCustomNodeById(ctx, d.Id())   // DescribeDBCustomNodes with NodeIds=[id]
if respData == nil -> d.SetId(""); return nil
set every field with nil guards (zone, image_id, vpc_id, subnet_id, node_type,
    node_name, auto_renew, tags, and all computed fields incl. system_disk/data_disks)
```

### Update
```
if d.HasChange("tags"):
    compute AddTags / DeleteTagKeys from old/new map
    resource.Retry: call ModifyDBCustomNodeTags(NodeId, AddTags, DeleteTagKeys)

if d.HasChange("auto_renew"):
    request = NewRenewDBCustomNodeRequest()
    request.NodeId = d.Id()
    request.Period = period (from schema, default 1)
    request.AutoRenew = new auto_renew
    set AutoVoucher / VoucherIds if present
    resource.Retry: call RenewDBCustomNode; guard result/Response
return Read
```
(All other arguments are `ForceNew`.)

### Delete (two-stage, both async)
```
// Stage 1: isolate (no TaskId returned)
resource.Retry(WriteRetryTimeout): call IsolateDBCustomNode(NodeId); guard result/Response
waitDBCustomNodeStatus(ctx, service, nodeId, "Isolated")   // poll DescribeDBCustomNodes
    - status == "Isolated" -> done
    - node not found        -> treat as already gone, done
    - otherwise             -> retryable

// Stage 2: destroy (returns TaskId)
resource.Retry(WriteRetryTimeout): call DestroyDBCustomNode(NodeId); capture TaskId
if TaskId != nil -> waitDBCustomTaskSucceeded(taskId)
return nil
```

## Service Layer Additions

- `DescribeDBCustomNodeById(ctx, nodeId)` — wraps `DescribeDBCustomNodes` with
  `NodeIds=[nodeId]`, `Limit=100`, inside `resource.Retry` + `ratelimit.Check`;
  returns the single `*dbdcv20201029.DBCustomNode` (nil if not found), nil-safe.
- Reuse existing `DescribeDBCustomTaskStatusById` and `waitDBCustomTaskSucceeded`.
- `waitDBCustomNodeStatus` helper (resource-file local) for the isolate stage.

## Import

`ImportStatePassthrough` on `NodeId`:
```
terraform import tencentcloud_dbdc_db_custom_node.example dbcn-xxxxxxxx
```
`password` (write-only/sensitive) will be ignored on import verification.

## Conventions

- Every SDK call wrapped in `resource.Retry` (`WriteRetryTimeout` for writes,
  `ReadRetryTimeout` for reads) with `tccommon.RetryError`, plus
  `ratelimit.Check(request.GetAction())`.
- All response/pointer accesses nil-checked before dereference; list/slice
  index access guarded by length checks (`NodeIds`, `NodeSet`).
- Pagination `Limit` uses the documented maximum (100).
- English comments and code only.
- Doc `resource_tc_dbdc_db_custom_node.md` and test
  `resource_tc_dbdc_db_custom_node_test.go` follow the
  `resource_tc_config_compliance_pack.{md,_test.go}` naming convention.
