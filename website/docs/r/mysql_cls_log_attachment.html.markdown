---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_cls_log_attachment"
sidebar_current: "docs-tencentcloud-resource-mysql_cls_log_attachment"
description: |-
  Provides a resource to create a mysql log to cls
---

# tencentcloud_mysql_cls_log_attachment

Provides a resource to create a mysql log to cls

~> **NOTE:** The CLS resource bound to resource `tencentcloud_mysql_cls_log_attachment` needs to be manually deleted.

## Example Usage

### Create Error Log to ClS

```hcl
# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-6"
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

# create security group
resource "tencentcloud_security_group" "security_group" {
  name        = "sg-mysql"
  description = "mysql test"
}

# create mysql instance
resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = "ap-guangzhou-6"
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

# attachment cls log
resource "tencentcloud_mysql_cls_log_attachment" "example" {
  instance_id      = tencentcloud_mysql_instance.example.id
  log_type         = "error"
  create_log_set   = true
  create_log_topic = true
  log_set          = "tf_log_set"
  log_topic        = "tf_log_topic"
  period           = 30
  create_index     = true
  cls_region       = "ap-guangzhou"
}
```

### Create Slow Log to ClS

```hcl
resource "tencentcloud_mysql_cls_log_attachment" "example" {
  instance_id = tencentcloud_mysql_instance.example.id
  log_type    = "slowlog"
  log_set     = "50d499a8-c4c0-4442-aa04-e8aa8a02437d"
  log_topic   = "140d4d39-4307-45a8-9655-290f679b063d"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The id of instance.
* `log_set` - (Required, String, ForceNew) If `create_log_set` is `true`, use log set name, Else use log set Id.
* `log_topic` - (Required, String, ForceNew) If `create_log_topic` is `true`, use log topic name, Else use log topic Id.
* `log_type` - (Required, String, ForceNew) Log type. Support `error` or `slowlog`.
* `cls_region` - (Optional, String) Cls region.
* `create_index` - (Optional, Bool, ForceNew) Whether to create index.
* `create_log_set` - (Optional, Bool, ForceNew) Whether to create log set.
* `create_log_topic` - (Optional, Bool, ForceNew) Whether to create log topic.
* `period` - (Optional, Int, ForceNew) The validity period of the log theme is 30 days by default when not filled in.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `log_set_id` - Log set Id.
* `log_topic_id` - Log topic Id.
* `status` - Log Status.


## Import

mysql log to cls can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_cls_log_attachment.example cdb-8fk7id2l#slowlog
```

