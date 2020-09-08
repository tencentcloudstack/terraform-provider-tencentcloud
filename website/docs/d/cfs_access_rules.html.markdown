---
subcategory: "Cloud File Storage(CFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfs_access_rules"
sidebar_current: "docs-tencentcloud-datasource-cfs_access_rules"
description: |-
  Use this data source to query the detail information of CFS access rule.
---

# tencentcloud_cfs_access_rules

Use this data source to query the detail information of CFS access rule.

## Example Usage

```hcl
data "tencentcloud_cfs_access_rules" "access_rules" {
  access_group_id = "pgroup-7nx89k7l"
  access_rule_id  = "rule-qcndbqzj"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_id` - (Required) A specified access group ID used to query.
* `access_rule_id` - (Optional) A specified access rule ID used to query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_rule_list` - An information list of CFS access rule. Each element contains the following attributes:
  * `access_rule_id` - ID of the access rule.
  * `auth_client_ip` - Allowed IP of the access rule.
  * `priority` - The priority level of access rule.
  * `rw_permission` - Read and write permissions.
  * `user_permission` - The permissions of accessing users.


