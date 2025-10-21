---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_project_security_groups"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_project_security_groups"
description: |-
  Use this data source to query detailed information of cynosdb project_security_groups
---

# tencentcloud_cynosdb_project_security_groups

Use this data source to query detailed information of cynosdb project_security_groups

## Example Usage

```hcl
data "tencentcloud_cynosdb_project_security_groups" "project_security_groups" {
  project_id = 1250480
  search_key = "自定义模版"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Optional, Int) Project ID.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Search Keywords.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `groups` - Security Group Details.
  * `create_time` - Creation time, time format: yyyy mm dd hh: mm: ss.
  * `inbound` - Inbound Rules.
    * `action` - Action.
    * `address_module` - AddressModule.
    * `cidr_ip` - CidrIp.
    * `desc` - Description.
    * `id` - id.
    * `ip_protocol` - Ip protocol.
    * `port_range` - PortRange.
    * `service_module` - Service Module.
  * `outbound` - Outbound rules.
    * `action` - Action.
    * `address_module` - Address module.
    * `cidr_ip` - Cidr Ip.
    * `desc` - Description.
    * `id` - id.
    * `ip_protocol` - Ip protocol.
    * `port_range` - Port range.
    * `service_module` - Service module.
  * `project_id` - Project ID.
  * `security_group_id` - Security Group ID.
  * `security_group_name` - Security Group Name.
  * `security_group_remark` - Security Group Notes.


