# Add GS Android Instances and KeeWiDB Instances Data Sources

## What

Add two new Terraform data sources:

1. `tencentcloud_gs_android_instances` — query Android instances via `DescribeAndroidInstances` (GS product, SDK: gs/v20190118)
2. `tencentcloud_keewidb_instances` — query KeeWiDB instances via `DescribeInstances` (KeeWiDB product, SDK: keewidb/v20220308)

## Why

Users need to query existing GS Android instances and KeeWiDB instances for use in downstream Terraform configurations.

## APIs Used

| Data Source | API | Pagination |
|---|---|---|
| `tencentcloud_gs_android_instances` | `DescribeAndroidInstances` | Offset/Limit, max Limit=100 |
| `tencentcloud_keewidb_instances` | `DescribeInstances` | Offset/Limit, max Limit=1000 |

## SDK Dependencies (New)

Both SDK packages are not yet in vendor and must be added:
- `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gs/v20190118`
- `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/keewidb/v20220308`

New client accessor methods must be added to `tencentcloud/connectivity/client.go`.
