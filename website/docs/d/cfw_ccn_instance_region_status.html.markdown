---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_ccn_instance_region_status"
sidebar_current: "docs-tencentcloud-datasource-cfw_ccn_instance_region_status"
description: |-
  Use this data source to query detailed information of CFW ccn instance region status
---

# tencentcloud_cfw_ccn_instance_region_status

Use this data source to query detailed information of CFW ccn instance region status

## Example Usage

```hcl
data "tencentcloud_cfw_ccn_instance_region_status" "example" {
  ccn_id = "ccn-fkb9bo2v"
  instance_ids = [
    "vpc-axbsvrrg"
  ]
  routing_mode = 1
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String) CCN ID.
* `instance_ids` - (Optional, Set: [`String`]) List of instance IDs associated with CCN for querying traffic steering network deployment status.
* `result_output_file` - (Optional, String) Used to save results.
* `routing_mode` - (Optional, Int) Traffic steering routing method, 0: multi-route table, 1: policy routing.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region_fw_status` - List of regional firewall traffic steering network status.
  * `cidr` - CIDR of the traffic steering network, empty if traffic steering network is not deployed.
  * `region` - Region.
  * `status` - Traffic steering network deployment status.
1. `NotDeployed` Firewall cluster not deployed.
2. `Deployed` Firewall cluster deployed, but traffic steering network not created.
3. `Auto` Firewall cluster deployed, and traffic steering network created with automatically selected network segment.
4. `Custom` Firewall cluster deployed, and traffic steering network created with user-defined network segment.


