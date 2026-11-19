---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_account"
sidebar_current: "docs-tencentcloud-resource-dcdb_account"
description: |-
  Provides a resource to create a dcdb account
---

# tencentcloud_dcdb_account

Provides a resource to create a dcdb account

## Example Usage

```hcl
resource "tencentcloud_dcdb_hourdb_instance" "example" {
  instance_name     = "tf-example"
  zones             = ["ap-guangzhou-6", "ap-guangzhou-7"]
  shard_memory      = "4"
  shard_storage     = "100"
  shard_node_count  = "2"
  shard_count       = "2"
  vpc_id            = "vpc-i5yyodl9"
  subnet_id         = "subnet-hhi88a58"
  security_group_id = "sg-4z20n68d"
  db_version_id     = "8.0"
  resource_tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}

resource "tencentcloud_dcdb_account" "example" {
  instance_id          = tencentcloud_dcdb_hourdb_instance.example.id
  user_name            = "mysql"
  host                 = "127.0.0.1"
  password             = "Password@2026"
  read_only            = 0
  description          = "this is a test account"
  max_user_connections = 10
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String) db host.
* `instance_id` - (Required, String) instance id.
* `password` - (Required, String) password.
* `user_name` - (Required, String) account name.
* `description` - (Optional, String) description for account.
* `max_user_connections` - (Optional, Int) max user connections.
* `read_only` - (Optional, Int) whether the account is readonly. 0 means not a readonly account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAccount](https://www.tencentcloud.com/document/api/1042/34449)
* Read: [DescribeAccounts](https://www.tencentcloud.com/document/api/1042/34446)
* Update:
  - [ModifyAccountDescription](https://www.tencentcloud.com/document/api/1042/34435)
  - [ResetAccountPassword](https://www.tencentcloud.com/document/api/1042/34430)
* Delete:
  - [DeleteAccount](https://www.tencentcloud.com/document/api/1042/34448)
  - [DescribeAccounts](https://www.tencentcloud.com/document/api/1042/34446)

## Import

dcdb account can be imported using the instanceId#userName, e.g.
```
terraform import tencentcloud_dcdb_account.example tdsqlshard-0tyhkwkzj5#mysql
```

