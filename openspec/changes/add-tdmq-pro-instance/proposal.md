## Why

TencentCloud TDMQ (Tencent Distributed Message Queue) provides professional clusters for high-performance messaging workloads. Currently, there is no Terraform resource to manage the lifecycle of TDMQ professional cluster instances. Users need a `tencentcloud_tdmq_pro_instance` resource to create, read, update, and delete professional clusters via Infrastructure as Code.

## What Changes

- Add a new Terraform resource `tencentcloud_tdmq_pro_instance` (RESOURCE_KIND_GENERAL) that manages the full lifecycle of a TDMQ professional cluster.
- Implement Create via `CreateProCluster` API (parameters: `ZoneIds`, `ProductName`, `AutoRenewFlag`, `TimeSpan`, `ClusterName`, `AutoVoucher`, `StorageSize`, `Vpc`, `Tags`).
- Implement Read via `DescribeClusters` API (filter by `ClusterIdList`).
- Implement Update via `ModifyCluster` API (parameters: `ClusterName`, `Remark`, `PublicAccessEnabled`).
- Implement Delete via `DeleteCluster` API (parameter: `ClusterId`).
- Note: The `InstanceVersion` parameter listed in the requirement does NOT exist in the current vendor SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217`) and will be excluded.
- Note: The `Tags` field in `CreateProCluster` is marked as deprecated in the SDK; standard tag management via `tencentcloud_tag` resource or the provider's built-in tag handling should be used instead.
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.

## Capabilities

### New Capabilities
- `tdmq-pro-instance-crud`: Full CRUD lifecycle management for TDMQ professional cluster instances, including creation with zone/VPC/storage configuration, reading cluster status and attributes, updating cluster name/remark/public access, and deletion.

### Modified Capabilities

(none)

## Impact

- **New files**:
  - `tencentcloud/services/tdmq/resource_tc_tdmq_pro_instance.go` — resource implementation
  - `tencentcloud/services/tdmq/resource_tc_tdmq_pro_instance_test.go` — unit tests with gomonkey mocks
  - `tencentcloud/services/tdmq/resource_tc_tdmq_pro_instance.md` — example usage documentation
- **Modified files**:
  - `tencentcloud/provider.go` — register the new resource
  - `tencentcloud/provider.md` — add resource entry
- **Dependencies**: Uses existing vendor SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217`
- **APIs consumed**: `CreateProCluster`, `DescribeClusters`, `ModifyCluster`, `DeleteCluster`
