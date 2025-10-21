---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_acl"
sidebar_current: "docs-tencentcloud-resource-ckafka_acl"
description: |-
  Provides a resource to create a Ckafka Acl.
---

# tencentcloud_ckafka_acl

Provides a resource to create a Ckafka Acl.

## Example Usage

### Ckafka Acl

```hcl
resource "tencentcloud_ckafka_user" "example" {
  instance_id  = "ckafka-7k5nbnem"
  account_name = "tf-example"
  password     = "Password@123"
}

resource "tencentcloud_ckafka_acl" "example" {
  instance_id     = "ckafka-7k5nbnem"
  resource_type   = "TOPIC"
  resource_name   = "tf-example-resource"
  operation_type  = "WRITE"
  permission_type = "ALLOW"
  host            = "*"
  principal       = tencentcloud_ckafka_user.example.account_name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the ckafka instance.
* `operation_type` - (Required, String, ForceNew) ACL operation mode. Valid values: `UNKNOWN`, `ANY`, `ALL`, `READ`, `WRITE`, `CREATE`, `DELETE`, `ALTER`, `DESCRIBE`, `CLUSTER_ACTION`, `DESCRIBE_CONFIGS` and `ALTER_CONFIGS`.
* `resource_name` - (Required, String, ForceNew) ACL resource name, which is related to `resource_type`. For example, if `resource_type` is `TOPIC`, this field indicates the topic name; if `resource_type` is `GROUP`, this field indicates the group name.
* `host` - (Optional, String, ForceNew) The default is *, which means that any host can access it. Support filling in IP or network segment, and support `;`separation.
* `permission_type` - (Optional, String, ForceNew) ACL permission type. Valid values: `UNKNOWN`, `ANY`, `DENY`, `ALLOW`. and `ALLOW` by default. Currently, CKafka supports `ALLOW` (equivalent to allow list), and other fields will be used for future ACLs compatible with open-source Kafka.
* `principal` - (Optional, String, ForceNew) User list. The default value is `*`, which means that any user can access. The current user can only be one included in the user list. For example: `root` meaning user root can access.
* `resource_type` - (Optional, String, ForceNew) ACL resource type. Valid values are `UNKNOWN`, `ANY`, `TOPIC`, `GROUP`, `CLUSTER`, `TRANSACTIONAL_ID`. and `TOPIC` by default. Currently, only `TOPIC` is available, and other fields will be used for future ACLs compatible with open-source Kafka.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Ckafka Acl can be imported using the instance_id#permission_type#principal#host#operation_type#resource_type#resource_name, e.g.

```
$ terraform import tencentcloud_ckafka_acl.example ckafka-7k5nbnem#ALLOW#tf-example#*#WRITE#TOPIC#tf-example-resource
```

