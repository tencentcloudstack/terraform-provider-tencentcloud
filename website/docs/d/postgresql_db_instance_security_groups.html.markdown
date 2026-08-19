---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_db_instance_security_groups"
sidebar_current: "docs-tencentcloud-datasource-postgresql_db_instance_security_groups"
description: |-
  Use this data source to query security groups associated with TencentCloud PostgreSQL instance or read-only group.
---

# tencentcloud_postgresql_db_instance_security_groups

Use this data source to query security groups associated with TencentCloud PostgreSQL instance or read-only group.

## Example Usage

### Query security groups associated with a PostgreSQL instance

```hcl
data "tencentcloud_postgresql_db_instance_security_groups" "example" {
  db_instance_id = "postgres-ckwcgdf1"
}
```

### Query security groups associated with a read-only group

```hcl
data "tencentcloud_postgresql_db_instance_security_groups" "example" {
  read_only_group_id = "pgrogrp-6fexzcmy"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Optional, String) Instance ID, which can be obtained via the DescribeDBInstances API. DBInstanceId, ReadOnlyGroupId at least pass one; if you want to query the security group associated with the instance, only pass the DBInstanceId field.
* `read_only_group_id` - (Optional, String) Read-only group ID, which can be obtained via the DescribeReadOnlyGroups API. DBInstanceId, ReadOnlyGroupId at least pass one; if you want to query the security group associated with the read-only group, only pass the ReadOnlyGroupId.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_group_set` - Security group information list of the instance.
  * `create_time` - Creation time.
  * `inbound` - Inbound rule.
    * `action` - Policy, ACCEPT or DROP.
    * `cidr_ip` - Source or destination IP or IP range, e.g. 172.16.0.0/12.
    * `description` - Rule description.
    * `ip_protocol` - Network protocol, supports UDP, TCP, etc.
    * `port_range` - Port.
  * `outbound` - Outbound rule.
    * `action` - Policy, ACCEPT or DROP.
    * `cidr_ip` - Source or destination IP or IP range, e.g. 172.16.0.0/12.
    * `description` - Rule description.
    * `ip_protocol` - Network protocol, supports UDP, TCP, etc.
    * `port_range` - Port.
  * `project_id` - Project ID.
  * `security_group_description` - Security group remark.
  * `security_group_id` - Security group ID.
  * `security_group_name` - Security group name.


