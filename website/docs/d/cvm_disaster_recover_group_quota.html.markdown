---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_disaster_recover_group_quota"
sidebar_current: "docs-tencentcloud-datasource-cvm_disaster_recover_group_quota"
description: |-
  Use this data source to query detailed information of cvm disaster_recover_group_quota
---

# tencentcloud_cvm_disaster_recover_group_quota

Use this data source to query detailed information of cvm disaster_recover_group_quota

## Example Usage

```hcl
data "tencentcloud_cvm_disaster_recover_group_quota" "disaster_recover_group_quota" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `current_num` - The number of placement groups that have been created by the current user.
* `cvm_in_host_group_quota` - Quota on instances in a physical-machine-type disaster recovery group.
* `cvm_in_rack_group_quota` - Quota on instances in a rack-type disaster recovery group.
* `cvm_in_sw_group_quota` - Quota on instances in a switch-type disaster recovery group.
* `group_quota` - The maximum number of placement groups that can be created.


