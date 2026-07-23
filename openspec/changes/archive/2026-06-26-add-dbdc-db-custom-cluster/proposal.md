# Add tencentcloud_dbdc_db_custom_cluster Resource

## What

Add a new Terraform resource `tencentcloud_dbdc_db_custom_cluster` for managing
Tencent Cloud DBDC (DB Custom) clusters. The resource supports the full CRUD
lifecycle. Both Create and Delete are asynchronous operations, so the resource
polls `DescribeDBCustomTaskStatus` until the returned task `Status` becomes
`Succeeded`. Update is limited to mutating cluster tags via
`ModifyDBCustomClusterTags` (the only mutable attribute exposed by the APIs).

## Why

The `dbdc` service currently only ships data sources
(`tencentcloud_dbdc_db_custom_clusters`, `..._nodes`, `..._cluster_nodes`,
`..._images`). There is no way to declaratively create or destroy a DB Custom
cluster through Terraform. This resource closes that gap so users can manage the
cluster lifecycle as infrastructure-as-code.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `CreateDBCustomCluster` | Async. Returns `ClusterId` (resource ID) and `TaskId` |
| Read | `DescribeDBCustomClusterDetail` | Query by `ClusterId` |
| Update (tags) | `ModifyDBCustomClusterTags` | `AddTags` / `DeleteTagKeys` diff |
| Delete | `DestroyDBCustomCluster` | Async. Returns `TaskId` |
| Poll (helper) | `DescribeDBCustomTaskStatus` | Input `TaskId`; `Status == "Succeeded"` means done |

## Resource ID

`ClusterId` returned by `CreateDBCustomCluster` (format `dbcc-xxxxxxxx`).

## SDK Availability

All five APIs are present in the vendored SDK
`tencentcloud-sdk-go/tencentcloud/dbdc/v20201029` (`client.go`), and the client
accessor `UseDbdcV20201029Client()` already exists in
`tencentcloud/connectivity/client.go`. No SDK upgrade is required.
