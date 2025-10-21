---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_publish_subscribe"
sidebar_current: "docs-tencentcloud-resource-sqlserver_publish_subscribe"
description: |-
  Provides a SQL Server PublishSubscribe resource belongs to SQL Server instance.
---

# tencentcloud_sqlserver_publish_subscribe

Provides a SQL Server PublishSubscribe resource belongs to SQL Server instance.

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_general_cloud_instance" "example_pub" {
  name                 = "tf-example-pub"
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_HSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  db_version           = "2008R2"
  security_group_list  = [tencentcloud_security_group.security_group.id]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
}

resource "tencentcloud_sqlserver_general_cloud_instance" "example_sub" {
  name                 = "tf-example-sub"
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_HSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  db_version           = "2008R2"
  security_group_list  = [tencentcloud_security_group.security_group.id]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
}

resource "tencentcloud_sqlserver_db" "example_pub" {
  instance_id = tencentcloud_sqlserver_general_cloud_instance.example_pub.id
  name        = "tf_example_db_pub"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_db" "example_sub" {
  instance_id = tencentcloud_sqlserver_general_cloud_instance.example_sub.id
  name        = "tf_example_db_sub"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_publish_subscribe" "example" {
  publish_instance_id    = tencentcloud_sqlserver_general_cloud_instance.example_pub.id
  subscribe_instance_id  = tencentcloud_sqlserver_general_cloud_instance.example_sub.id
  publish_subscribe_name = "example"
  delete_subscribe_db    = false
  database_tuples {
    publish_database   = tencentcloud_sqlserver_db.example_pub.name
    subscribe_database = tencentcloud_sqlserver_db.example_sub.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `database_tuples` - (Required, Set) Database Publish and Publish relationship list. The elements inside can be deleted and added individually, but modification is not allowed.
* `publish_instance_id` - (Required, String, ForceNew) ID of the SQL Server instance which publish.
* `subscribe_instance_id` - (Required, String, ForceNew) ID of the SQL Server instance which subscribe.
* `delete_subscribe_db` - (Optional, Bool) Whether to delete the subscriber database when deleting the Publish and Subscribe. `true` for deletes the subscribe database, `false` for does not delete the subscribe database. default is `false`.
* `publish_subscribe_name` - (Optional, String) The name of the Publish and Subscribe. Default is `default_name`.

The `database_tuples` object supports the following:

* `publish_database` - (Required, String) Publish the database.
* `subscribe_database` - (Required, String) Subscribe the database.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SQL Server PublishSubscribe can be imported using the publish_sqlserver_id#subscribe_sqlserver_id, e.g.

```
$ terraform import tencentcloud_sqlserver_publish_subscribe.example publish_sqlserver_id#subscribe_sqlserver_id
```

