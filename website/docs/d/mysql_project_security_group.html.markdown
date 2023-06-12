---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_project_security_group"
sidebar_current: "docs-tencentcloud-datasource-mysql_project_security_group"
description: |-
  Use this data source to query detailed information of mysql project_security_group
---

# tencentcloud_mysql_project_security_group

Use this data source to query detailed information of mysql project_security_group

## Example Usage

```hcl
data "tencentcloud_mysql_project_security_group" "project_security_group" {
  project_id = 1250480
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Optional, Int) project id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `groups` - Security group details.
  * `create_time` - Creation time, time format: yyyy-mm-dd hh:mm:sss.
  * `inbound` - inbound rules.
    * `action` - Policy, ACCEPT or DROPs.
    * `cidr_ip` - Source IP or IP range, such as 192.168.0.0/16.
    * `desc` - Rule description.
    * `dir` - The direction defined by the rule, the inbound rule is INPUT.
    * `ip_protocol` - Network protocol, support UDP, TCP, etc.
    * `port_range` - port.
  * `outbound` - outbound rules.
    * `action` - Policy, ACCEPT or DROP.
    * `cidr_ip` - Destination IP or IP segment, such as 172.16.0.0/12.
    * `desc` - Rule description.
    * `dir` - The direction defined by the rule, the inbound rule is OUTPUT.
    * `ip_protocol` - Network protocol, support UDP, TCP, etc.
    * `port_range` - port or port range.
  * `project_id` - project id.
  * `security_group_id` - Security group ID.
  * `security_group_name` - Security group name.
  * `security_group_remark` - Security group remark.


