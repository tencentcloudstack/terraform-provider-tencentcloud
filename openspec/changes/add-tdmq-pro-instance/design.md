## Context

The TencentCloud Terraform Provider currently has no resource for managing TDMQ professional cluster instances. The TDMQ SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217`) already provides the necessary APIs: `CreateProCluster`, `DescribeClusters`, `ModifyCluster`, and `DeleteCluster`. This resource follows the standard RESOURCE_KIND_GENERAL pattern used throughout the provider.

The resource will be placed at `tencentcloud/services/tdmq/resource_tc_tdmq_pro_instance.go` following the established file organization conventions.

## Goals / Non-Goals

**Goals:**
- Provide full CRUD lifecycle management for TDMQ professional cluster instances via Terraform.
- Support creation with zone selection, product specification, VPC binding, storage configuration, and billing parameters.
- Support updating cluster name, remark, and public access settings.
- Support importing existing clusters by cluster ID.
- Follow existing provider patterns (retry logic, error handling, logging conventions).

**Non-Goals:**
- Managing TDMQ topics, namespaces, or subscriptions within the cluster (separate resources).
- Supporting the deprecated `Tags` field in `CreateProCluster` (use standard tag management).
- Implementing the `InstanceVersion` parameter (not available in current vendor SDK).
- Scaling/resizing operations (not supported by current APIs).

## Decisions

### 1. Resource ID: Use `ClusterId` as the single resource ID

**Rationale**: The `CreateProCluster` response returns `ClusterId` which uniquely identifies the cluster. All subsequent operations (Read/Update/Delete) use this single ID. No composite ID is needed.

**Alternative considered**: Using `ClusterName` — rejected because names can be changed via `ModifyCluster` and are not guaranteed unique.

### 2. Schema Design: Separate create-only and updatable fields

**Rationale**: The Create and Update APIs accept different parameter sets:
- Create-only (ForceNew): `zone_ids`, `product_name`, `auto_renew_flag`, `time_span`, `auto_voucher`, `storage_size`, `vpc`
- Updatable: `cluster_name`, `remark`, `public_access_enabled`

Fields that are only in the Create API will be marked as `ForceNew: true` to trigger resource recreation if changed.

**Alternative considered**: Using immutableArgs pattern — rejected because ForceNew is the standard Terraform approach for create-only fields and provides better UX (plan shows destroy+create).

### 3. Read Implementation: Use `DescribeClusters` with `ClusterIdList` filter

**Rationale**: The `DescribeClusters` API supports filtering by `ClusterIdList`, which allows querying a specific cluster by ID. The response returns a `ClusterSet` array; we take the first matching element.

**Alternative considered**: A dedicated `DescribeCluster` (singular) API — not available in the SDK.

### 4. Error Handling: Standard retry with `tccommon.ReadRetryTimeout`

**Rationale**: All API calls will be wrapped in `resource.Retry()` with `tccommon.ReadRetryTimeout` for read operations and `tccommon.WriteRetryTimeout` for write operations, using `tccommon.RetryError()` for error wrapping. This follows the established provider pattern.

### 5. Tags: Exclude deprecated `Tags` field from Create

**Rationale**: The `Tags` field in `CreateProClusterRequest` is marked as deprecated ("集群的标签列表(已废弃)") in the SDK. Standard tag management should be handled separately.

## Risks / Trade-offs

- **[Risk] CreateProCluster is a billing API** → The resource creates paid instances. Tests must use gomonkey mocks to avoid real charges. Documentation should clearly indicate this creates billable resources.
- **[Risk] No dedicated describe-single-cluster API** → Must filter from list API. If the cluster is not found in the response, the resource will be marked as removed from state (standard Terraform pattern).
- **[Risk] InstanceVersion parameter missing from SDK** → The requirement lists this parameter but it does not exist in the vendor SDK. Mitigation: exclude it and document the limitation.
- **[Trade-off] ForceNew vs immutableArgs** → Using ForceNew means changing create-only fields will destroy and recreate the cluster (with potential data loss). This is the standard Terraform behavior and users should be aware.
