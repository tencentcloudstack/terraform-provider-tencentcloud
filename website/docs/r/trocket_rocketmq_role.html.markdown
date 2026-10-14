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

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateRole](https://cloud.tencent.com/document/api/1493/98864) - This API is used to add a role.
* Read: [DescribeRoleList](https://cloud.tencent.com/document/api/1493/98862) - This API is used to query the list of roles. Filter parameter usage instructions are as follows:.
* Update: [ModifyRole](https://cloud.tencent.com/document/api/1493/98861) - This API is used to modify a role.
* Delete: [DeleteRole](https://cloud.tencent.com/document/api/1493/98863) - This API is used to delete a role. Make sure that the related information on this role is not used in the current code. After the role is deleted, the keys (AccessKey and SecretKey) used to produce or consume messages with this role become invalid immediately.

## Import

trocket rocketmq_role can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_role.rocketmq_role instanceId#role
```

