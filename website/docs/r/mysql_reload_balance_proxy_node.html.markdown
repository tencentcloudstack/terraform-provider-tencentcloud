---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_reload_balance_proxy_node"
sidebar_current: "docs-tencentcloud-resource-mysql_reload_balance_proxy_node"
description: |-
  Provides a resource to create a mysql reload_balance_proxy_node
---

# tencentcloud_mysql_reload_balance_proxy_node

Provides a resource to create a mysql reload_balance_proxy_node

## Example Usage

```hcl
resource "tencentcloud_mysql_reload_balance_proxy_node" "reload_balance_proxy_node" {
  proxy_group_id   = "proxy-gmi1f78l"
  proxy_address_id = "proxyaddr-4wc4y1pq"
}
```

## Argument Reference

The following arguments are supported:

* `proxy_group_id` - (Required, String, ForceNew) Proxy id.
* `proxy_address_id` - (Optional, String, ForceNew) Proxy address id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



