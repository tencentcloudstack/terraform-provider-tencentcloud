---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfs_access_rule"
sidebar_current: "docs-tencentcloud-resource-cfs_access_rule"
description: |-
  Provides a resource to create a CFS access rule.
---

# tencentcloud_cfs_access_rule

Provides a resource to create a CFS access rule.

## Example Usage

```hcl
resource "tencentcloud_cfs_access_rule" "foo" {
  access_group_id = "pgroup-7nx89k7l"
  auth_client_ip  = "10.10.1.0/24"
  priority        = 1
  rw_permission   = "RO"
  user_permission = "root_squash"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_id` - (Required, ForceNew) ID of a access group.
* `auth_client_ip` - (Required) A single IP or a single IP address range such as 10.1.10.11 or 10.10.1.0/24 indicates that all IPs are allowed. Please note that the IP entered should be CVM's private IP.
* `priority` - (Required) The priority level of rule. The range is 1-100, and 1 indicates the highest priority.
* `rw_permission` - (Optional) Read and write permissions. Valid values are `RO` and `RW`, and default is `RO`.
* `user_permission` - (Optional) The permissions of accessing users. Valid values are `all_squash`, `no_all_squash`, `root_squash` and `no_root_squash`, and default is `root_squash`. `all_squash` indicates that all access users are mapped as anonymous users or user groups; `no_all_squash` indicates that access users will match local users first and be mapped to anonymous users or user groups after matching failed; `root_squash` indicates that map access root users to anonymous users or user groups; `no_root_squash` indicates that access root users keep root account permission.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



