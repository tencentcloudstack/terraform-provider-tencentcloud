## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/cfw/resource_tc_cfw_cluster_nat_fw_switch.go` with schema definition and CRUD functions (Create calls `OpenClusterNatFwSwitch`, Read calls `DescribeNatCcnFwSwitch`, Update calls `ModifyClusterNatFwSwitch`, Delete calls `CloseClusterNatFwSwitch`). Use composite ID `nat_ins_id#ccn_id` with `tccommon.FILED_SP` separator. Include `Importer` for import support.
- [x] 1.2 Add `DescribeNatCcnFwSwitchById` helper method to `tencentcloud/services/cfw/service_tencentcloud_cfw.go` for the Read operation.

## 2. Documentation

- [x] 2.1 Create `tencentcloud/services/cfw/resource_tc_cfw_cluster_nat_fw_switch.md` with one-sentence description, Example Usage (using `nat_ccn_switch` block), and Import section showing composite ID format.

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_cfw_cluster_nat_fw_switch` in `tencentcloud/provider.go` under the CFW section (reference `tencentcloud_igtm_strategy` pattern).
- [x] 3.2 Add `tencentcloud_cfw_cluster_nat_fw_switch` entry to `tencentcloud/provider.md` under the CFW section.

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/cfw/resource_tc_cfw_cluster_nat_fw_switch_test.go` with unit tests using gomonkey to mock cloud API calls (`OpenClusterNatFwSwitch`, `DescribeNatCcnFwSwitch`, `ModifyClusterNatFwSwitch`, `CloseClusterNatFwSwitch`). Run with `go test -gcflags=all=-l`.
