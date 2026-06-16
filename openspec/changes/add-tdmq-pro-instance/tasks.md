## 1. Resource Schema and CRUD Implementation

- [x] 1.1 Create `tencentcloud/services/tdmq/resource_tc_tdmq_pro_instance.go` with schema definition including: `zone_ids` (Required, ForceNew, List of int), `product_name` (Required, ForceNew, String), `storage_size` (Required, ForceNew, Int), `vpc` (Required, ForceNew, List MaxItems=1 with `vpc_id` and `subnet_id`), `cluster_name` (Optional, String), `time_span` (Optional, ForceNew, Int), `auto_renew_flag` (Optional, ForceNew, Int), `auto_voucher` (Optional, ForceNew, Int), `remark` (Optional, String), `public_access_enabled` (Optional, Bool). Computed fields: `cluster_id`, `status`, `deal_name`, `big_deal_id`.
- [x] 1.2 Implement Create function: call `CreateProCluster` API with retry, validate non-nil response and non-empty `ClusterId`, set resource ID, then call Read.
- [x] 1.3 Implement Read function: call `DescribeClusters` API with `ClusterIdList` filter, handle empty response (log ID then SetId("")), populate state attributes from `Cluster` struct fields.
- [x] 1.4 Implement Update function: call `ModifyCluster` API with changed fields (`cluster_name`, `remark`, `public_access_enabled`), then call Read to refresh state.
- [x] 1.5 Implement Delete function: call `DeleteCluster` API with `ClusterId`, handle not-found as success.

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_tdmq_pro_instance` resource in `tencentcloud/provider.go` resource map.
- [x] 2.2 Add `tencentcloud_tdmq_pro_instance` entry in `tencentcloud/provider.md`.

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/tdmq/resource_tc_tdmq_pro_instance.md` with Example Usage (showing required and optional parameters) and Import section.

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/tdmq/resource_tc_tdmq_pro_instance_test.go` with gomonkey-based unit tests covering Create, Read, Update, and Delete operations. Run tests with `go test -gcflags=all=-l`.

## 5. Verification

- [x] 5.1 Run `gofmt` on new Go files to ensure formatting compliance.
