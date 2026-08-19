# Design: tencentcloud_dbdc_node_to_db_custom_cluster_attachment Resource

## Architecture

Follows the `tencentcloud_organization_member_auth_policy_attachment` style
(Create / Read / Delete only, composite ID split with a separator) and reuses
the async task helper already in the `dbdc` package:

```
provider.go (register tencentcloud_dbdc_node_to_db_custom_cluster_attachment)
    └─ resource_tc_dbdc_node_to_db_custom_cluster_attachment.go (Create/Read/Delete)
           └─ service_tencentcloud_dbdc.go
                  ├─ DescribeDBCustomClusterNodeById(ctx, clusterId, nodeId)  (new)
                  └─ waitDBCustomTaskSucceeded(...)                           (existing, same package)
                         └─ dbdc SDK v20201029
```

Package `dbdc`, SDK alias `dbdcv20201029`, client
`...GetAPIV3Conn().UseDbdcV20201029Client()`.

## Composite ID

`ClusterId#NodeId`, joined/split with the literal `#` separator (per business
rule). Read/Delete split the ID and validate it has exactly 2 parts.

## Schema

Arguments map to `AddNodesToDBCustomClusterRequest`. The API accepts a `NodeIds`
list, but since the resource ID is `ClusterId#NodeId` (one node per attachment),
a single `node_id` argument is exposed and sent as `NodeIds=[node_id]`. All
arguments are `ForceNew` (attachment resource has no in-place update).

### Required (ForceNew)

| Field | Type | Maps to | Description |
|---|---|---|---|
| `cluster_id` | String | `ClusterId` | DB Custom cluster ID |
| `node_id` | String | `NodeIds=[node_id]` | DB Custom node ID to add to the cluster |

### Optional (ForceNew)

| Field | Type | Maps to | Description |
|---|---|---|---|
| `image_id` | String | `ImageId` | OS image to reset the node to after it is added |
| `login_settings` | List (MaxItems 1) | `LoginSettings` | Login config (see below) |

`login_settings` nested block (ForceNew):

| Field | Type | Maps to | Notes |
|---|---|---|---|
| `password` | String (Sensitive) | `LoginSettings.Password` | Login password |
| `key_ids` | List of String | `LoginSettings.KeyIds` | Key pair IDs (single ID supported) |
| `keep_image_login` | String | `LoginSettings.KeepImageLogin` | `true`/`false` |

### Computed (read from `DescribeDBCustomClusterNodes` -> `DBCustomClusterNode`)

| Field | Type | Maps to |
|---|---|---|
| `node_name` | String | `NodeName` |
| `lan_ip` | String | `LanIP` |
| `ssh_endpoint` | String | `SSHEndpoint` |
| `status` | String | `Status` |
| `zone` | String | `Zone` |
| `node_type` | String | `NodeType` |

## CRUD Logic

### Create
```
request = NewAddNodesToDBCustomClusterRequest()
request.ClusterId = cluster_id
request.NodeIds = [node_id]
set ImageId if provided
set LoginSettings{Password, KeyIds, KeepImageLogin} if block present
resource.Retry(WriteRetryTimeout): call AddNodesToDBCustomCluster; guard result/Response; capture TaskId
if TaskId != nil -> waitDBCustomTaskSucceeded(taskId)
d.SetId(cluster_id + "#" + node_id)
return Read
```

### Read
```
idSplit = split(d.Id(), "#"); require len == 2
clusterId, nodeId = idSplit[0], idSplit[1]
node = DescribeDBCustomClusterNodeById(ctx, clusterId, nodeId)
if node == nil -> d.SetId(""); return nil
_ = d.Set("cluster_id", clusterId); _ = d.Set("node_id", nodeId)
set computed fields with nil guards (node_name, lan_ip, ssh_endpoint, status, zone, node_type)
```

### Delete
```
idSplit = split(d.Id(), "#"); require len == 2
request = NewRemoveNodesFromDBCustomClusterRequest()
request.ClusterId = clusterId; request.NodeIds = [nodeId]
resource.Retry(WriteRetryTimeout): call RemoveNodesFromDBCustomCluster; guard result/Response; capture TaskId
if TaskId != nil -> waitDBCustomTaskSucceeded(taskId)
return nil
```

(No Update — all arguments `ForceNew`.)

## Service Layer Additions

- `DescribeDBCustomClusterNodeById(ctx, clusterId, nodeId)` — wraps
  `DescribeDBCustomClusterNodes` with `ClusterId`, paginates with `Limit=100`
  inside `resource.Retry` + `ratelimit.Check`, returns the matching
  `*dbdcv20201029.DBCustomClusterNode` (nil if not found), nil/length-safe.
- Reuse existing `waitDBCustomTaskSucceeded`.

## Import

`ImportStatePassthrough` on the composite ID:
```
terraform import tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example dbcc-xxxxxxxx#dbcn-xxxxxxxx
```
`password`, `image_id`, `login_settings` are write-only/create-time and are
ignored on import verification.

## Conventions

- Every SDK call wrapped in `resource.Retry` (`WriteRetryTimeout` for writes,
  `ReadRetryTimeout` for reads) with `tccommon.RetryError` + `ratelimit.Check`.
- All response/pointer accesses nil-checked; slice index access guarded.
- Pagination `Limit` uses the documented maximum (100).
- English comments and code only.
- Doc `resource_tc_dbdc_node_to_db_custom_cluster_attachment.md` and test
  `resource_tc_dbdc_node_to_db_custom_cluster_attachment_test.go` follow the
  `resource_tc_config_compliance_pack.{md,_test.go}` naming convention.
