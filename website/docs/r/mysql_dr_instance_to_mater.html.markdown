---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_dr_instance_to_mater"
sidebar_current: "docs-tencentcloud-resource-mysql_dr_instance_to_mater"
description: |-
  Provides a resource to create a mysql dr_instance_to_mater
---

# tencentcloud_mysql_dr_instance_to_mater

Provides a resource to create a mysql dr_instance_to_mater

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-mysql"
  description = "mysql test"
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  slave_sync_mode   = 1
  instance_name     = "tf-example-mysql"
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [tencentcloud_security_group.security_group.id]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_mysql_dr_instance_to_mater" "example" {
  instance_id = tencentcloud_mysql_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Disaster recovery instance ID in the format of cdb-c1nl9rpv. It is the same as the instance ID displayed in the TencentDB console.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeAsyncRequestInfo](https://cloud.tencent.com/document/api/236/20410) - This API (DescribeAsyncRequestInfo) is used to query the async task execution result of a TencentDB instance.
* Create: [SwitchDrInstanceToMaster](https://cloud.tencent.com/document/api/236/52562) - This API is used to switch a cloud database disaster recovery instance to primary instance. Note that the request must be sent to the region where the disaster recovery instance is located.
* Read: [DescribeDBInstances](https://cloud.tencent.com/document/api/236/15872) - This API is used to query the list of TencentDB for MySQL instances. It supports filtering instances by conditions such as project ID, instance ID, access address, and instance status. It also supports querying the list of information about primary instances, disaster recovery instances, and read-only instances.
* Update: [DescribeAsyncRequestInfo](https://cloud.tencent.com/document/api/236/20410) - This API (DescribeAsyncRequestInfo) is used to query the async task execution result of a TencentDB instance.
* Update: [SwitchDrInstanceToMaster](https://cloud.tencent.com/document/api/236/52562) - This API is used to switch a cloud database disaster recovery instance to primary instance. Note that the request must be sent to the region where the disaster recovery instance is located.

## Import

mysql dr_instance_to_mater can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_dr_instance_to_mater.dr_instance_to_mater dr_instance_to_mater_id
```

