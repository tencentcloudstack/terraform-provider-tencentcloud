---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_nat_fw_cluster_region_status"
sidebar_current: "docs-tencentcloud-datasource-cfw_nat_fw_cluster_region_status"
description: |-
  Use this data source to query detailed information of CFW NAT firewall cluster region status
---

# tencentcloud_cfw_nat_fw_cluster_region_status

Use this data source to query detailed information of CFW NAT firewall cluster region status

## Example Usage

```hcl
data "tencentcloud_cfw_nat_fw_cluster_region_status" "example" {
  nat_cluster_region_status_query_list {
    ccn_id       = "ccn-p3mlp0tj"
    nat_ins_id   = "nat-h1i1mf4n"
    asset_type   = "nat_ccn"
    routing_mode = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `nat_cluster_region_status_query_list` - (Required, List) List of query conditions for NAT firewall cluster region status.
* `result_output_file` - (Optional, String) Used to save results.

The `nat_cluster_region_status_query_list` object supports the following:

* `asset_type` - (Required, String) Asset type. Valid values: `nat_ccn` (CCN+NAT scenario), `nat` (standalone NAT scenario).
* `ccn_id` - (Required, String) CCN ID.
* `nat_ins_id` - (Required, String) NAT gateway ID.
* `routing_mode` - (Optional, Int) Traffic steering routing method. 0: multi-route table mode, 1: policy routing mode.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region_fw_status` - List of regional firewall cluster status.
  * `ccn_id` - CCN ID.
  * `cidr` - Traffic steering network CIDR. Only has value when Status is Auto or Custom.
  * `nat_ins_id` - NAT gateway ID.
  * `region` - Region, e.g. ap-guangzhou.
  * `routing_mode` - Traffic steering routing method. 0: multi-route table mode, 1: policy routing mode.
  * `status` - Region cluster status. Valid values: `NotDeployed` (cluster not deployed), `Deployed` (cluster deployed but traffic steering network not created), `DeployedCustomOnly` (cluster deployed but internal segment covered, need custom traffic steering segment), `Auto` (traffic steering network created with auto-assigned CIDR), `Custom` (traffic steering network created with custom CIDR).


