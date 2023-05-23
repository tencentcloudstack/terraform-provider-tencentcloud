---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_project_security_groups"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_project_security_groups"
description: |-
  Use this data source to query detailed information of sqlserver project_security_groups
---

# tencentcloud_sqlserver_project_security_groups

Use this data source to query detailed information of sqlserver project_security_groups

## Example Usage

```hcl
data "tencentcloud_sqlserver_project_security_groups" "project_security_groups" {
  project_id = 0
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, Int) Project ID, which can be viewed through the console project management.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_group_set` - Security group details.
  * `create_time` - Creation time, time format: yyyy-mm-dd hh:mm:ss.
  * `inbound_set` - inbound rules.
    * `action` - Policy, ACCEPT or DROP.
    * `cidr_ip` - Destination IP or IP segment, such as 172.16.0.0/12.
    * `dir` - The direction defined by the rules, OUTPUT-outgoing rules INPUT-inbound rules.
    * `ip_protocol` - Network protocol, support UDP, TCP, etc.
    * `port_range` - port or port range.
  * `outbound_set` - outbound rules.
    * `action` - Policy, ACCEPT or DROP.
    * `cidr_ip` - Destination IP or IP segment, such as 172.16.0.0/12.
    * `dir` - The direction defined by the rules, OUTPUT-outgoing rules INPUT-inbound rules.
    * `ip_protocol` - Network protocol, support UDP, TCP, etc.
    * `port_range` - port or port range.
  * `project_id` - project ID.
  * `security_group_id` - Security group ID.
  * `security_group_name` - security group name.
  * `security_group_remark` - Security Group Remarks.


