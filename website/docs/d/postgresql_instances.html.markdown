---
subcategory: "PostgreSQL"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_instances"
sidebar_current: "docs-tencentcloud-datasource-postgresql_instances"
description: |-
  Use this data source to query postgresql instances
---

# tencentcloud_postgresql_instances

Use this data source to query postgresql instances

## Example Usage

```hcl
data "tencentcloud_postgresql_instances" "name" {
  name = "test"
}

data "tencentcloud_postgresql_instances" "project" {
  project_id = 0
}

data "tencentcloud_postgresql_instances" "id" {
  id = "postgres-h9t4fde1"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) ID of the postgresql instance to be query.
* `name` - (Optional) Name of the postgresql instance to be query.
* `project_id` - (Optional) Project ID of the postgresql instance to be query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - A list of postgresql instances. Each element contains the following attributes.
  * `auto_renew_flag` - Auto renew flag.
  * `availability_zone` - Availability zone.
  * `charge_type` - Pay type of the postgresql instance.
  * `charset` - Charset of the postgresql instance.
  * `create_time` - Create time of the postgresql instance.
  * `engine_version` - Version of the postgresql database engine.
  * `id` - ID of the postgresql instance.
  * `memory` - Memory size(in GB).
  * `name` - Name of the postgresql instance.
  * `private_access_ip` - IP address for private access.
  * `private_access_port` - Port for private access.
  * `project_id` - Project id, default value is 0.
  * `public_access_host` - Host for public access.
  * `public_access_port` - Port for public access.
  * `public_access_switch` - Indicates whether to enable the access to an instance from public network or not.
  * `root_user` - Instance root account name, default value is `root`.
  * `storage` - Volume size(in GB).
  * `subnet_id` - ID of subnet.
  * `tags` - The available tags within this postgresql.
  * `vpc_id` - ID of VPC.


