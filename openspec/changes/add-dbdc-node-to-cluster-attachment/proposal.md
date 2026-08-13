# Add tencentcloud_dbdc_node_to_db_custom_cluster_attachment Resource

## What

Add a new Terraform attachment resource
`tencentcloud_dbdc_node_to_db_custom_cluster_attachment` that binds a single DB
Custom node to a DB Custom cluster (and unbinds on delete). It is a
bind/unbind-style resource with only Create / Read / Delete (no Update — all
arguments are `ForceNew`).

Both add and remove are asynchronous: the resource polls
`DescribeDBCustomTaskStatus` (input `TaskId`, success when
`Status == "Succeeded"`).

## Why

The `dbdc` service has resources for clusters and nodes, plus a
`tencentcloud_dbdc_db_custom_cluster_nodes` data source, but no way to manage
the membership of a node within a cluster declaratively. This attachment closes
that gap.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `AddNodesToDBCustomCluster` | Async. Returns `TaskId` |
| Read | `DescribeDBCustomClusterNodes` | Query by `ClusterId`; `Limit` max 100; locate node by `NodeId` |
| Delete | `RemoveNodesFromDBCustomCluster` | Async. Returns `TaskId` |
| Poll (helper) | `DescribeDBCustomTaskStatus` | Input `TaskId`; `Status == "Succeeded"` means done |

## Resource ID

Composite `ClusterId#NodeId` (e.g. `dbcc-xxxxxxxx#dbcn-xxxxxxxx`). One resource
maps to exactly one node bound to one cluster.

## SDK Availability

`AddNodesToDBCustomCluster`, `DescribeDBCustomClusterNodes`,
`RemoveNodesFromDBCustomCluster` and `DescribeDBCustomTaskStatus` are all present
in the vendored SDK `tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`, and
`UseDbdcV20201029Client()` already exists. The shared helper
`waitDBCustomTaskSucceeded` lives in the same `dbdc` package and is reused. No
SDK upgrade is required.
