---
subcategory: "DBbrain"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_sql_filter"
sidebar_current: "docs-tencentcloud-resource-dbbrain_sql_filter"
description: |-
  Provides a resource to create a dbbrain sql_filter.
---

# tencentcloud_dbbrain_sql_filter

Provides a resource to create a dbbrain sql_filter.

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}
variable "region" {
  default = "ap-guangzhou"
}

data "tencentcloud_mysql_instance" "mysql" {
  instance_name = "instance_name"
}

locals {
  mysql_id = data.tencentcloud_mysql_instance.mysql.instance_list.0.mysql_id
}

resource "tencentcloud_dbbrain_sql_filter" "sql_filter" {
  instance_id = local.mysql_id
  session_token {
    user     = "test"
    password = "===password==="
  }
  sql_type        = "SELECT"
  filter_key      = "filter_key"
  max_concurrency = 10
  duration        = 3600
}
```

## Argument Reference

The following arguments are supported:

* `duration` - (Required, Int) filter duration.
* `filter_key` - (Required, String) filter key.
* `instance_id` - (Required, String) instance id.
* `max_concurrency` - (Required, Int) maximum concurreny.
* `session_token` - (Required, List) session token.
* `sql_type` - (Required, String) sql type, optional value is SELECT, UPDATE, DELETE, INSERT, REPLACE.
* `product` - (Optional, String) product, optional value is &amp;#39;mysql&amp;#39;, &amp;#39;cynosdb&amp;#39;.
* `status` - (Optional, String) filter status.

The `session_token` object supports the following:

* `password` - (Required, String) password.
* `user` - (Required, String) user name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `filter_id` - filter id.


