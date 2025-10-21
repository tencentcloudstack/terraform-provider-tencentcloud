---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_project_security_groups"
sidebar_current: "docs-tencentcloud-datasource-dcdb_project_security_groups"
description: |-
  Use this data source to query detailed information of dcdb project_security_groups
---

# tencentcloud_dcdb_project_security_groups

Use this data source to query detailed information of dcdb project_security_groups

## Example Usage

```hcl
data "tencentcloud_dcdb_project_security_groups" "project_security_groups" {
  product    = "dcdb"
  project_id = 0
}
```

## Argument Reference

The following arguments are supported:

* `product` - (Required, String) Database engine name. Valid value: `dcdb`.
* `project_id` - (Optional, Int) Project ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `groups` - Security group details.
  * `create_time` - Creation time in the format of yyyy-mm-dd hh:mm:ss.
  * `inbound` - Inbound rule.
    * `action` - Policy, which can be `ACCEPT` or `DROP`.
    * `cidr_ip` - Source IP or source IP range, such as 192.168.0.0/16.
    * `ip_protocol` - Network protocol. UDP and TCP are supported.
    * `port_range` - Port.
  * `outbound` - Outbound rule.
    * `action` - Policy, which can be `ACCEPT` or `DROP`.
    * `cidr_ip` - Source IP or source IP range, such as 192.168.0.0/16.
    * `ip_protocol` - Network protocol. UDP and TCP are supported.
    * `port_range` - Port.
  * `project_id` - Project ID.
  * `security_group_id` - Security group ID.
  * `security_group_name` - Security group name.
  * `security_group_remark` - Security group remarks.


