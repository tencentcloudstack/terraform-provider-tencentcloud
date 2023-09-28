---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_trocket_rocketmq_role"
sidebar_current: "docs-tencentcloud-resource-trocket_rocketmq_role"
description: |-
  Provides a resource to create a trocket rocketmq_role
---

# tencentcloud_trocket_rocketmq_role

Provides a resource to create a trocket rocketmq_role

## Example Usage

```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test_role"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-xxxxx"
  subnet_id     = "subnet-xxxxx"
  tags = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_role" "rocketmq_role" {
  instance_id = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  role        = "test_role"
  remark      = "test for terraform"
  perm_write  = false
  perm_read   = true
}

output "access_key" {
  value = tencentcloud_trocket_rocketmq_role.rocketmq_role.access_key
}

output "secret_key" {
  value = tencentcloud_trocket_rocketmq_role.rocketmq_role.secret_key
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of instance.
* `perm_read` - (Required, Bool) Whether to enable consumption permission.
* `perm_write` - (Required, Bool) Whether to enable production permission.
* `remark` - (Required, String) remark.
* `role` - (Required, String, ForceNew) Name of role.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `access_key` - Access key.
* `created_time` - Created time.
* `modified_time` - Modified time.
* `secret_key` - Secret key.


## Import

trocket rocketmq_role can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_role.rocketmq_role instanceId#role
```

