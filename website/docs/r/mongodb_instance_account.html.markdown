---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_account"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance_account"
description: |-
  Provides a resource to create a mongodb instance_account
---

# tencentcloud_mongodb_instance_account

Provides a resource to create a mongodb instance_account

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance_account" "instance_account" {
  instance_id         = "cmgo-lxaz2c9b"
  user_name           = "test_account"
  password            = "xxxxxxxx"
  mongo_user_password = "xxxxxxxxx"
  user_desc           = "test account"
  auth_role {
    mask      = 0
    namespace = "*"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.
* `mongo_user_password` - (Required, String, ForceNew) The password corresponding to the mongouser account. mongouser is the system default account, which is the password set when creating an instance.
* `password` - (Required, String) New account password. Password complexity requirements are as follows: character length range [8,32]. Contains at least letters, numbers and special characters (exclamation point!, at@, pound sign #, percent sign %, caret ^, asterisk *, parentheses (), underscore _).
* `user_name` - (Required, String, ForceNew) The new account name. Its format requirements are as follows: character range [1,32]. Characters in the range of [A,Z], [a,z], [1,9] as well as underscore _ and dash - can be input.
* `auth_role` - (Optional, List) The read and write permission information of the account.
* `user_desc` - (Optional, String) Account remarks.

The `auth_role` object supports the following:

* `mask` - (Required, Int) Permission information of the current account. 0: No permission. 1: read-only. 2: Write only. 3: Read and write.
* `namespace` - (Required, String) Refers to the name of the database with the current account permissions.*: Indicates all databases. db.name: Indicates the database of a specific name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



