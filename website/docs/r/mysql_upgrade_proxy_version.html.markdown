---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_upgrade_proxy_version"
sidebar_current: "docs-tencentcloud-resource-mysql_upgrade_proxy_version"
description: |-
  Provides a resource to create a mysql upgrade_proxy_version
---

# tencentcloud_mysql_upgrade_proxy_version

Provides a resource to create a mysql upgrade_proxy_version

## Example Usage

```hcl
resource "tencentcloud_mysql_upgrade_proxy_version" "upgrade_proxy_version" {
  instance_id       = "cdb-fitq5t9h"
  proxy_group_id    = "proxy-h1ub486b"
  src_proxy_version = "1.3.6"
  dst_proxy_version = "1.3.7"
  upgrade_time      = "nowTime"
}
```

## Argument Reference

The following arguments are supported:

* `dst_proxy_version` - (Required, String, ForceNew) Database agent upgrade version.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `proxy_group_id` - (Required, String, ForceNew) Database proxy ID.
* `src_proxy_version` - (Required, String, ForceNew) The current version of the database agent.
* `upgrade_time` - (Required, String, ForceNew) Upgrade time: nowTime (upgrade completed) timeWindow (instance maintenance time).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



