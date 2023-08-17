---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_publish_subscribes"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_publish_subscribes"
description: |-
  Use this data source to query Publish Subscribe resources for the specific SQL Server instance.
---

# tencentcloud_sqlserver_publish_subscribes

Use this data source to query Publish Subscribe resources for the specific SQL Server instance.

## Example Usage

```hcl
data "tencentcloud_sqlserver_publish_subscribes" "example" {
  instance_id = tencentcloud_sqlserver_publish_subscribe.example.publish_instance_id
}

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

* `instance_id` - (Required, String) ID of the SQL Server instance.
* `pub_or_sub_instance_id` - (Optional, String) The subscribe/publish instance ID. It is related to whether the `instance_id` is a publish instance or a subscribe instance. when `instance_id` is a publish instance, this field is filtered according to the subscribe instance ID; when `instance_id` is a subscribe instance, this field is filtering according to the publish instance ID.
* `pub_or_sub_instance_ip` - (Optional, String) The intranet IP of the subscribe/publish instance. It is related to whether the `instance_id` is a publish instance or a subscribe instance. when `instance_id` is a publish instance, this field is filtered according to the intranet IP of the subscribe instance; when `instance_id` is a subscribe instance, this field is based on the publish instance intranet IP filter.
* `publish_database` - (Optional, String) Name of publish database.
* `publish_subscribe_id` - (Optional, Int) The id of the Publish and Subscribe.
* `publish_subscribe_name` - (Optional, String) The name of the Publish and Subscribe.
* `result_output_file` - (Optional, String) Used to store results.
* `subscribe_database` - (Optional, String) Name of subscribe database.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `publish_subscribe_list` - Publish and subscribe list. Each element contains the following attributes.
  * `database_tuples` - Database Publish and Publish relationship list.
    * `last_sync_time` - Last sync time.
    * `publish_database` - Name of the publish SQL Server instance.
    * `status` - Publish and subscribe status between databases, valid values are `running`, `success`, `fail`, `unknow`.
    * `subscribe_database` - Name of the subscribe SQL Server instance.
  * `publish_instance_id` - ID of the SQL Server instance which publish.
  * `publish_instance_ip` - IP of the the SQL Server instance which publish.
  * `publish_instance_name` - Name of the SQL Server instance which publish.
  * `publish_subscribe_id` - The id of the Publish and Subscribe.
  * `publish_subscribe_name` - The name of the Publish and Subscribe.
  * `subscribe_instance_id` - ID of the SQL Server instance which subscribe.
  * `subscribe_instance_ip` - IP of the SQL Server instance which subscribe.
  * `subscribe_instance_name` - Name of the SQL Server instance which subscribe.


