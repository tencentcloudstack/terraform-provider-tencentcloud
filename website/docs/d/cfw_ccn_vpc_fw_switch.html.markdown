---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_ccn_vpc_fw_switch"
sidebar_current: "docs-tencentcloud-datasource-cfw_ccn_vpc_fw_switch"
description: |-
  Use this data source to query detailed information of CFW ccn vpc fw switch
---

# tencentcloud_cfw_ccn_vpc_fw_switch

Use this data source to query detailed information of CFW ccn vpc fw switch

## Example Usage

```hcl
data "tencentcloud_cfw_ccn_vpc_fw_switch" "example" {
  ccn_id = "ccn-fkb9bo2v"
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String) CCN ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `interconnect_pairs` - Interconnect pair configuration.
  * `group_a` - Group A.
    * `access_cidr_list` - List of network segments for accessing firewall.
    * `access_cidr_mode` - Network segment mode for accessing firewall: 0-no access, 1-access all network segments associated with the instance, 2-access user-defined network segments.
    * `instance_id` - Instance ID.
    * `instance_region` - Region where the instance is located.
    * `instance_type` - Instance type such as VPC or DIRECTCONNECT.
  * `group_b` - Group B.
    * `access_cidr_list` - List of network segments for accessing firewall.
    * `access_cidr_mode` - Network segment mode for accessing firewall: 0-no access, 1-access all network segments associated with the instance, 2-access user-defined network segments.
    * `instance_id` - Instance ID.
    * `instance_region` - Region where the instance is located.
    * `instance_type` - Instance type such as VPC or DIRECTCONNECT.
  * `interconnect_mode` - Interconnect mode: "CrossConnect": cross interconnect (each instance in group A interconnects with each instance in group B), "FullMesh": full mesh (group A content is identical to group B, equivalent to pairwise interconnection within the group).


