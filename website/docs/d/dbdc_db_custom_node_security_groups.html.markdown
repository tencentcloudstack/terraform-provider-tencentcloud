---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_node_security_groups"
sidebar_current: "docs-tencentcloud-datasource-dbdc_db_custom_node_security_groups"
description: |-
  Use this data source to query DBDC custom node security groups information.
---

# tencentcloud_dbdc_db_custom_node_security_groups

Use this data source to query DBDC custom node security groups information.

## Example Usage

```hcl
data "tencentcloud_dbdc_db_custom_node_security_groups" "example" {
  node_id = "dbcn-abc12345"
}
```

## Argument Reference

The following arguments are supported:

* `node_id` - (Required, String) DB Custom node ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `groups` - Security group list bound to the DB Custom node.
  * `create_time` - Security group creation time.
  * `inbound` - Inbound rules.
    * `action` - Rule action, ACCEPT or DROP.
    * `address_module` - IP address template ID.
    * `cidr_ip` - Source/Destination IP or CIDR.
    * `desc` - Rule description.
    * `id` - Security group ID.
    * `ip_protocol` - Protocol type, e.g., tcp, udp, icmp, ALL.
    * `port_range` - Port range.
    * `service_module` - Protocol port template ID.
  * `outbound` - Outbound rules.
    * `action` - Rule action, ACCEPT or DROP.
    * `address_module` - IP address template ID.
    * `cidr_ip` - Source/Destination IP or CIDR.
    * `desc` - Rule description.
    * `id` - Security group ID.
    * `ip_protocol` - Protocol type, e.g., tcp, udp, icmp, ALL.
    * `port_range` - Port range.
    * `service_module` - Protocol port template ID.
  * `project_id` - Project ID.
  * `security_group_id` - Security group ID.
  * `security_group_name` - Security group name.
  * `security_group_remark` - Security group remark.


