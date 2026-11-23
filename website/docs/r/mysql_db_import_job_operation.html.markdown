---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_db_import_job_operation"
sidebar_current: "docs-tencentcloud-resource-mysql_db_import_job_operation"
description: |-
  Provides a resource to create a mysql db_import_job_operation
---

# tencentcloud_mysql_db_import_job_operation

Provides a resource to create a mysql db_import_job_operation

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

resource "tencentcloud_mysql_account" "example" {
  mysql_id             = tencentcloud_mysql_instance.example.id
  name                 = "tf_example"
  password             = "Qwer@234"
  description          = "desc."
  max_user_connections = 10
}

resource "tencentcloud_mysql_db_import_job_operation" "example" {
  instance_id = tencentcloud_mysql_instance.example.id
  user        = tencentcloud_mysql_account.example.name
  password    = tencentcloud_mysql_account.example.password
  db_name     = "tf_example_db"
  file_name   = "tf_mysql.sql"
  cos_url     = "https://terraform-ci-1308919341.cos.ap-guangzhou.myqcloud.com/mysql/mysql.sql?q-sign-algorithm=sha1&q-ak=AKIDRnMWiUNr14F29GvCwOSHu9l_FdCdORqAxblrE10nDaO6mVI701oXTe-gL1QpClgW&q-sign-time=1684921483;1684925083&q-key-time=1684921483;1684925083&q-header-list=host&q-url-param-list=&q-signature=7410be4ef93075aebca459af4e617f8bcaa36f48&x-cos-security-token=EzDm9S6aRDwBLQcaxUNfb0TA30PqhOTa7d82a06a36e94b66bdbc6d09064a397bZypr0mD3oVkbJR9bRYix6BSDVYncX3Y2VCGYK6V2jFWZqIuEHoWJCe-2pDvJDNbMjF3ttWfLMqEouOkxNk28ay9NPHtMXrJgEEMb95BMAhGwi38oA2LjYfQRkk7AHesg2toSf11hiTAjVv-alf5uEidWGnFKe_6BgmnADYvtPptgXHNtsUZCxc33PF6tGBqX"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.
* `user` - (Required, String, ForceNew) The username of the cloud database.
* `cos_url` - (Optional, String, ForceNew) The name of the imported target database, if it is not passed, it means that no database is specified.
* `db_name` - (Optional, String, ForceNew) The name of the imported target database, if it is not passed, it means that no database is specified.
* `file_name` - (Optional, String, ForceNew) file name. This file refers to the file that the user has uploaded to Tencent Cloud, and only `.sql` files are supported.
* `password` - (Optional, String, ForceNew) The password of the user account of the cloud database instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `async_request_id` - The request ID of the asynchronous task.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateDBImportJob](https://www.tencentcloud.com/document/api/236/15858)
  - [DescribeAsyncRequestInfo](https://www.tencentcloud.com/document/api/236/34462)
* Read: [DescribeDBImportRecords](https://www.tencentcloud.com/document/api/236/15856)

