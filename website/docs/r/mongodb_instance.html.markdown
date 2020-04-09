---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance"
description: |-
  Provide a resource to create a Mongodb instance.
---

# tencentcloud_mongodb_instance

Provide a resource to create a Mongodb instance.

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "mongodb"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_3_WT"
  machine_type   = "GIO"
  available_zone = "ap-guangzhou-2"
  vpc_id         = "vpc-mz3efvbw"
  subnet_id      = "subnet-lk0svi3p"
  project_id     = 0
  password       = "mypassword"
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Required, ForceNew) The available zone of the Mongodb.
* `engine_version` - (Required, ForceNew) Version of the Mongodb, and available values include MONGO_3_WT, MONGO_3_ROCKS and MONGO_36_WT.
* `instance_name` - (Required) Name of the Mongodb instance.
* `machine_type` - (Required, ForceNew) Type of Mongodb instance, and available values include GIO and TGIO.
* `memory` - (Required) Memory size. The minimum value is 2, and unit is GB.
* `password` - (Required) Password of this Mongodb account.
* `volume` - (Required) Disk size. The minimum value is 25, and unit is GB.
* `project_id` - (Optional) ID of the project which the instance belongs.
* `security_groups` - (Optional) ID of the security group.
* `subnet_id` - (Optional, ForceNew) ID of the subnet within this VPC. The vaule is required if VpcId is set.
* `tags` - (Optional) The tags of the Mongodb.
* `vpc_id` - (Optional, ForceNew) ID of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the Mongodb instance.
* `status` - Status of the Mongodb instance, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2).
* `vip` - IP of the Mongodb instance.
* `vport` - IP port of the Mongodb instance.


## Import

Mongodb instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_instance.mongodb cmgo-41s6jwy4
```

