---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_acl"
sidebar_current: "docs-tencentcloud-resource-ckafka_acl"
description: |-
  Provides a resource to create a Ckafka Acl.
---

# tencentcloud_ckafka_acl

Provides a resource to create a Ckafka Acl.

## Example Usage

Ckafka Acl

```hcl
resource "tencentcloud_ckafka_acl" "foo" {
  instance_id     = "ckafka-f9ife4zz"
  resource_type   = "TOPIC"
  resource_name   = "topic-tf-test"
  operation_type  = "WRITE"
  permission_type = "ALLOW"
  host            = "*"
  principal       = tencentcloud_ckafka_user.foo.account_name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Id of the ckafka instance.
* `operation_type` - (Required, ForceNew) ACL operation mode, valid values are `UNKNOWN`, `ANY`, `ALL`, `READ`, `WRITE`, `CREATE`, `DELETE`, `ALTER`, `DESCRIBE`, `CLUSTER_ACTION`, `DESCRIBE_CONFIGS` and `ALTER_CONFIGS`.
* `resource_name` - (Required, ForceNew) ACL resource name, which is related to `resource_type`. For example, if `resource_type` is `TOPIC`, this field indicates the topic name; if `resource_type` is `GROUP`, this field indicates the group name.
* `host` - (Optional, ForceNew) IP address allowed to access. The default value is `*`, which means that any host can access.
* `permission_type` - (Optional, ForceNew) ACL permission type, valid values are `UNKNOWN`, `ANY`, `DENY`, `ALLOW`, and `ALLOW` by default. Currently, CKafka supports `ALLOW` (equivalent to allow list), and other fields will be used for future ACLs compatible with open-source Kafka.
* `principal` - (Optional, ForceNew) User list. The default value is `*`, which means that any user can access. The current user can only be one included in the user list.
* `resource_type` - (Optional, ForceNew) ACL resource type. Valid values are `UNKNOWN`, `ANY`, `TOPIC`, `GROUP`, `CLUSTER`, `TRANSACTIONAL_ID`, and `TOPIC` by default. Currently, only `TOPIC` is available, and other fields will be used for future ACLs compatible with open-source Kafka.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Ckafka acl can be imported using the instance_id#permission_type#principal#host#operation_type#resource_type#resource_name, e.g.

```
$ terraform import tencentcloud_ckafka_acl.foo ckafka-f9ife4zz#ALLOW#test#*#WRITE#TOPIC#topic-tf-test
```

