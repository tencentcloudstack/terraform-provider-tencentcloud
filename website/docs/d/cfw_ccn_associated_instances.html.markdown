---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_ccn_associated_instances"
sidebar_current: "docs-tencentcloud-datasource-cfw_ccn_associated_instances"
description: |-
  Use this data source to query detailed information of CFW ccn associated instances
---

# tencentcloud_cfw_ccn_associated_instances

Use this data source to query detailed information of CFW ccn associated instances

## Example Usage

```hcl
data "tencentcloud_cfw_ccn_associated_instances" "example" {
  ccn_id = "ccn-fkb9bo2v"
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String) CCN ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ccn_associated_instances` - Information of instances associated with CCN.
  * `cidr_lst` - List of network segments for the instance.
  * `ins_type` - Instance type.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `instance_region` - Region where the instance belongs.


