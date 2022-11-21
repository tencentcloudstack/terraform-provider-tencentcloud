---
subcategory: "TDSQL for MySQL(dcdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_security_groups"
sidebar_current: "docs-tencentcloud-datasource-dcdb_security_groups"
description: |-
  Use this data source to query detailed information of dcdb securityGroups
---

# tencentcloud_dcdb_security_groups

Use this data source to query detailed information of dcdb securityGroups

## Example Usage

```hcl
data "tencentcloud_dcdb_security_groups" "securityGroups" {
  instance_id = "your_instance_id"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - security group list.
  * `create_time` - create time.
  * `inbound` - inbound rules.
    * `action` - policy action.
    * `cidr_ip` - cidr ip.
    * `ip_protocol` - internet protocol.
    * `port_range` - port range.
  * `outbound` - outbound rules.
    * `action` - policy action.
    * `cidr_ip` - cidr ip.
    * `ip_protocol` - internet protocol.
    * `port_range` - port range.
  * `project_id` - project id.
  * `security_group_id` - security group id.
  * `security_group_name` - security group name.


