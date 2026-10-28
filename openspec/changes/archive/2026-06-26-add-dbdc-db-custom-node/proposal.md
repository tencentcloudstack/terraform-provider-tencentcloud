# Add tencentcloud_dbdc_db_custom_node Resource

## What

Add a new Terraform resource `tencentcloud_dbdc_db_custom_node` for managing a
single Tencent Cloud DBDC (DB Custom) node. One Terraform resource maps to
exactly one node: `NodeCount` defaults to `1` and the resource ID is taken from
`NodeIds[0]` returned by `CreateDBCustomNodes`.

Create and delete are asynchronous. The delete path is two-staged
(`IsolateDBCustomNode` then `DestroyDBCustomNode`). The resource polls
`DescribeDBCustomTaskStatus` (input `TaskId`, success when `Status == "Succeeded"`)
for the create and destroy tasks, and polls node status for the isolate stage
(since `IsolateDBCustomNode` does not return a task id).

## Why

The `dbdc` service currently exposes only data sources for nodes
(`tencentcloud_dbdc_db_custom_nodes`) and the cluster resource
(`tencentcloud_dbdc_db_custom_cluster`). There is no resource to declaratively
buy, renew, retag, isolate and destroy an individual DB Custom node. This
resource closes that gap.

> Note: the request originally referenced the name
> `tencentcloud_dbdc_db_custom_cluster`, but every API is node-scoped and the
> business rules describe a single node per resource. The confirmed resource
> name is `tencentcloud_dbdc_db_custom_node`.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `CreateDBCustomNodes` | Async. `NodeCount=1`; ID = `NodeIds[0]`; returns `TaskId` |
| Read | `DescribeDBCustomNodes` | Query by `NodeIds=[id]`; `Limit` max 100 |
| Update (renew) | `RenewDBCustomNode` | Renew / auto-renew settings |
| Update (tags) | `ModifyDBCustomNodeTags` | `AddTags` / `DeleteTagKeys` diff |
| Delete (stage 1) | `IsolateDBCustomNode` | No `TaskId`; poll node status until `Isolated` |
| Delete (stage 2) | `DestroyDBCustomNode` | Async. Returns `TaskId` |
| Poll (helper) | `DescribeDBCustomTaskStatus` | Input `TaskId`; `Status == "Succeeded"` means done |

## Resource ID

`NodeIds[0]` returned by `CreateDBCustomNodes` (format `dbcn-xxxxxxxx`).

## SDK Availability

All six node APIs plus `DescribeDBCustomTaskStatus` are present in the vendored
SDK `tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`, and the client accessor
`UseDbdcV20201029Client()` already exists. The shared helpers
`DescribeDBCustomTaskStatusById` and `waitDBCustomTaskSucceeded` were added with
the cluster resource and live in the same `dbdc` package, so they are reused. No
SDK upgrade is required.
