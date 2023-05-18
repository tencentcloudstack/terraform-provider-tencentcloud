---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_proxy_custom"
sidebar_current: "docs-tencentcloud-datasource-mysql_proxy_custom"
description: |-
  Use this data source to query detailed information of mysql proxy_custom
---

# tencentcloud_mysql_proxy_custom

Use this data source to query detailed information of mysql proxy_custom

## Example Usage

```hcl
data "tencentcloud_mysql_proxy_custom" "proxy_custom" {
  instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instanced id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `custom_conf` - proxy configuration.
  * `cpu` - number of cores.
  * `device_type` - Equipment type.
  * `device` - equipment.
  * `memory` - Memory.
  * `type` - type.
* `weight_rule` - weight limit.
  * `less_than` - division ceiling.
  * `weight` - weight limit.


