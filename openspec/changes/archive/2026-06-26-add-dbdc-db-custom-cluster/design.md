# Design: tencentcloud_dbdc_db_custom_cluster Resource

## Architecture

Follows the `tencentcloud_igtm_monitor` style (CRUD handlers in the resource
file, a thin Read/poll helper on the service struct):

```
provider.go (register tencentcloud_dbdc_db_custom_cluster)
    └─ resource_tc_dbdc_db_custom_cluster.go (Create/Read/Update/Delete)
           └─ service_tencentcloud_dbdc.go
                  ├─ DescribeDBCustomClusterById(ctx, clusterId)
                  └─ DescribeDBCustomTaskStatusById(ctx, taskId)
                         └─ dbdc SDK v20201029
```

Package: `dbdc`. SDK alias: `dbdcv20201029`. Client:
`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client()`.

## Schema

The arguments map 1:1 to the `CreateDBCustomClusterRequest` parameters. Only
`Tags` is mutable through an API (`ModifyDBCustomClusterTags`); every other
create input is therefore `ForceNew`.

### Required (Create, ForceNew)

| Field | Type | Maps to | Description |
|---|---|---|---|
| `cluster_name` | String | `ClusterName` | Cluster name (<=128 chars, CN/EN/underscore) |
| `container_network` | List (MaxItems 1, object) | `ContainerNetwork` | Pod network for the cluster |
| `container_network.vpc_id` | String (Required) | `ContainerNetwork.VpcId` | Container network VPC ID |
| `container_network.subnet_ids` | List of String (Required) | `ContainerNetwork.SubnetIds` | Container network subnet ID list |
| `api_server_network` | List (MaxItems 1, object) | `ApiServerNetwork` | API Server network info |
| `api_server_network.vpc_id` | String (Required) | `ApiServerNetwork.VpcId` | API Server VPC ID |
| `api_server_network.subnet_id` | String (Required) | `ApiServerNetwork.SubnetId` | API Server subnet ID |

### Optional

| Field | Type | ForceNew | Maps to | Description |
|---|---|---|---|---|
| `cluster_description` | String | Yes | `ClusterDescription` | Cluster description |
| `client_token` | String | Yes | `ClientToken` | Idempotency token (not read back) |
| `tags` | Map[string]string | No | `Tags` ([]Tag Key/Value) | Cluster tags (mutable via ModifyDBCustomClusterTags) |

### Computed (read from `DescribeDBCustomClusterDetail`)

| Field | Type | Maps to |
|---|---|---|
| `region` | String | `Region` |
| `cluster_status` | String | `ClusterStatus` (`Creating`/`Running`/`Destroying`) |
| `cluster_version` | String | `ClusterVersion` |
| `cluster_node_num` | Int | `ClusterNodeNum` |
| `cluster_level` | String | `ClusterLevel` |
| `created_time` | String | `CreatedTime` |

> Note: `tags` is modeled as a Terraform `TypeMap` (idiomatic) while the SDK
> uses `[]*Tag{Key,Value}`. Create converts the map into `[]*Tag`; Read converts
> `[]*Tag` back to a map; Update diffs old/new keys into `AddTags` /
> `DeleteTagKeys`.

## Async Task Polling

Both Create and Delete return a `TaskId` (uint64). A shared helper polls task
status inside a `resource.Retry` loop:

```
waitDBCustomTaskSucceeded(ctx, service, taskId):
    resource.Retry(WriteRetryTimeout, func():
        status := DescribeDBCustomTaskStatusById(ctx, taskId)   // with retry inside
        switch status:
            "Succeeded" -> return nil
            "Failed"    -> return NonRetryableError
            default ("Running"/nil) -> return RetryableError("task still running")
    )
```

`DescribeDBCustomTaskStatusById` itself wraps `DescribeDBCustomTaskStatus` in a
`resource.Retry` for transient API errors, and guards every pointer
(`result`, `result.Response`, `result.Response.Status`) before dereferencing.

## CRUD Logic

### Create
```
request = NewCreateDBCustomClusterRequest()
set ClusterName, ContainerNetwork, ApiServerNetwork, ClusterDescription, ClientToken, Tags from schema
resource.Retry(WriteRetryTimeout): call CreateDBCustomCluster, guard result/Response, capture response
if response.ClusterId == nil -> error
if response.TaskId != nil -> waitDBCustomTaskSucceeded(taskId)
d.SetId(*ClusterId)
return Read
```

### Read
```
respData = DescribeDBCustomClusterById(ctx, d.Id())
if respData == nil -> d.SetId(""); return nil   // resource gone
set each field with nil guards (cluster_name, cluster_description, container_network,
    api_server_network, tags, region, cluster_status, cluster_version,
    cluster_node_num, cluster_level, created_time)
```

### Update
```
if d.HasChange("tags"):
    old, new := d.GetChange("tags")
    AddTags     = new entries (Key/Value)
    DeleteTagKeys = keys present in old but absent in new
    request.ClusterId = d.Id()
    resource.Retry(WriteRetryTimeout): call ModifyDBCustomClusterTags, guard result/Response
return Read
```
(All other arguments are `ForceNew`, so no other update path is needed.)

### Delete
```
request = NewDestroyDBCustomClusterRequest(); request.ClusterId = d.Id()
resource.Retry(WriteRetryTimeout): call DestroyDBCustomCluster, guard result/Response, capture TaskId
if TaskId != nil -> waitDBCustomTaskSucceeded(taskId)
return nil
```

## Import

`ImportStatePassthrough` on `ClusterId`:
```
terraform import tencentcloud_dbdc_db_custom_cluster.example dbcc-xxxxxxxx
```

## Conventions

- Every SDK call wrapped in `resource.Retry` with `WriteRetryTimeout` (write
  paths) or `ReadRetryTimeout` (read path), using `tccommon.RetryError`.
- `ratelimit.Check(request.GetAction())` before each call (matches existing
  dbdc service file).
- All response/pointer accesses nil-checked before dereference.
- English comments and code only.
- Doc file `resource_tc_dbdc_db_custom_cluster.md` and acceptance test
  `resource_tc_dbdc_db_custom_cluster_test.go` follow the naming convention of
  `resource_tc_config_compliance_pack.{md,_test.go}`.
