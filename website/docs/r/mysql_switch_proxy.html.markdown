---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_switch_proxy"
sidebar_current: "docs-tencentcloud-resource-mysql_switch_proxy"
description: |-
  Provides a resource to create a mysql switch_proxy
---

# tencentcloud_mysql_switch_proxy

Provides a resource to create a mysql switch_proxy

## Example Usage

```hcl
resource "tencentcloud_mysql_switch_proxy" "switch_proxy" {
  instance_id    = "cdb-fitq5t9h"
  proxy_group_id = "proxy-h1ub486b"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `proxy_group_id` - (Required, String, ForceNew) Proxy group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



