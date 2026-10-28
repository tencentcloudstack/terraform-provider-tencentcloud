## Why

The CFW (Cloud Firewall) product supports NAT CCN cluster mode firewall switches, but there is currently no Terraform resource to manage the lifecycle of these switches. Users need a way to open, configure, and close NAT CCN cluster mode firewall switches via Terraform.

## What Changes

- Add new Terraform resource `tencentcloud_cfw_cluster_nat_fw_switch` (RESOURCE_KIND_GENERAL) to manage NAT CCN cluster mode firewall switches in CFW.
- The resource supports Create (open switch), Read (query switch config), Update (modify switch config), and Delete (close switch) operations.
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.

## Capabilities

### New Capabilities
- `cfw-cluster-nat-fw-switch`: Manages the NAT CCN cluster mode firewall switch lifecycle, including opening, reading, modifying, and closing the switch via CFW cloud APIs.

### Modified Capabilities

## Impact

- New file: `tencentcloud/services/cfw/resource_tc_cfw_cluster_nat_fw_switch.go`
- New file: `tencentcloud/services/cfw/resource_tc_cfw_cluster_nat_fw_switch_test.go`
- New file: `tencentcloud/services/cfw/resource_tc_cfw_cluster_nat_fw_switch.md`
- Modified: `tencentcloud/provider.go` (register new resource)
- Modified: `tencentcloud/provider.md` (document new resource)
- Cloud APIs used: `OpenClusterNatFwSwitch`, `DescribeNatCcnFwSwitch`, `ModifyClusterNatFwSwitch`, `CloseClusterNatFwSwitch` from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904`
