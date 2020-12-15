---
subcategory: "Ckafka"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_user"
sidebar_current: "docs-tencentcloud-resource-ckafka_user"
description: |-
  Provides a resource to create a Ckafka user.
---

# tencentcloud_ckafka_user

Provides a resource to create a Ckafka user.

## Example Usage

Ckafka User

```hcl
resource "tencentcloud_ckafka_user" "foo" {
  instance_id  = "ckafka-f9ife4zz"
  account_name = "tf-test"
  password     = "test1234"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, ForceNew) Account name used to access to ckafka instance.
* `instance_id` - (Required, ForceNew) ID of the ckafka instance.
* `password` - (Required) Password of the account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the account.
* `update_time` - The last update time of the account.


## Import

Ckafka user can be imported using the instance_id#account_name, e.g.

```
$ terraform import tencentcloud_ckafka_user.foo ckafka-f9ife4zz#tf-test
```

