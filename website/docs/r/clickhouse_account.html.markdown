---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_account"
sidebar_current: "docs-tencentcloud-resource-clickhouse_account"
description: |-
  Provides a resource to create a clickhouse account
---

# tencentcloud_clickhouse_account

Provides a resource to create a clickhouse account

## Example Usage

```hcl
resource "tencentcloud_clickhouse_account" "account" {
  instance_id = "cdwch-xxxxxx"
  user_name   = "test"
  password    = "xxxxxx"
  describe    = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `password` - (Required, String) Password.
* `user_name` - (Required, String) User name.
* `describe` - (Optional, String) Description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

clickhouse account can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_account.account ${instance_id}#${user_name}
```

