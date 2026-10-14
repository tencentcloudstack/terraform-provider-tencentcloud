---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_user"
sidebar_current: "docs-tencentcloud-resource-ckafka_user"
description: |-
  Provides a resource to create a Ckafka user.
---

# tencentcloud_ckafka_user

Provides a resource to create a Ckafka user.

## Example Usage

### Ckafka User

```hcl
resource "tencentcloud_ckafka_user" "example" {
  instance_id  = "ckafka-7k5nbnem"
  account_name = "tf-example"
  password     = "Password@123"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, String, ForceNew) Account name used to access to ckafka instance.
* `instance_id` - (Required, String, ForceNew) ID of the ckafka instance.
* `password` - (Required, String) Password of the account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the account.
* `update_time` - The last update time of the account.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateUser](https://www.tencentcloud.com/document/product/597/35367) - This API is used to add a user.
* Read: [DescribeInstanceAttributes](https://www.tencentcloud.com/document/product/597/35358) - This API is used to obtain instance attributes.
* Read: [DescribeUser](https://www.tencentcloud.com/document/product/597/35363) - This API is used to query user information.
* Update: [ModifyPassword](https://www.tencentcloud.com/document/product/597/35362) - This API is used to change the password.
* Delete: [DeleteUser](https://www.tencentcloud.com/document/product/597/35365) - This API is used to delete a user.

## Import

Ckafka user can be imported using the instance_id#account_name, e.g.

```
$ terraform import tencentcloud_ckafka_user.example ckafka-7k5nbnem#tf-example
```

