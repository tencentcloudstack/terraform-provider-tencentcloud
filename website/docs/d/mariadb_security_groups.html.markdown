---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_security_groups"
sidebar_current: "docs-tencentcloud-datasource-mariadb_security_groups"
description: |-
  Use this data source to query detailed information of mariadb securityGroups
---

# tencentcloud_mariadb_security_groups

Use this data source to query detailed information of mariadb securityGroups

## Example Usage

```hcl
data "tencentcloud_mariadb_security_groups" "securityGroups" {
  instance_id = "tdsql-4pzs5b67"
  product     = "mariadb"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `product` - (Required, String) product name, fixed to mariadb.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - SecurityGroup list.
  * `create_time` - Creation time, time format: `yyyy-mm-dd hh:mm:ss`.
  * `inbound` - Inbound rules.
    * `action` - Policy, ACCEPT or DROP.
    * `cidr_ip` - Source IP or IP range, such as 192.168.0.0/16.
    * `ip_protocol` - Network protocols, support `UDP`, `TCP`, etc.
    * `port_range` - port.
  * `outbound` - Outbound Rules.
    * `action` - Policy, ACCEPT or DROP.
    * `cidr_ip` - Source IP or IP range, such as 192.168.0.0/16.
    * `ip_protocol` - Network protocols, support `UDP`, `TCP`, etc.
    * `port_range` - port.
  * `project_id` - Project ID.
  * `security_group_id` - Security group ID.
  * `security_group_name` - security group name.
  * `security_group_remark` - Security Group Notes.


