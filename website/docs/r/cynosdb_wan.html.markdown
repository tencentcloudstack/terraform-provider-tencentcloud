---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_wan"
sidebar_current: "docs-tencentcloud-resource-cynosdb_wan"
description: |-
  Provides a resource to create a cynosdb wan
---

# tencentcloud_cynosdb_wan

Provides a resource to create a cynosdb wan

## Example Usage

```hcl
resource "tencentcloud_cynosdb_wan" "wan" {
  cluster_id      = "cynosdbmysql-bws8h88b"
  instance_grp_id = "cynosdbmysql-grp-lxav0p9z"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `instance_grp_id` - (Required, String) Instance Group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `wan_domain` - Domain name.
* `wan_ip` - Network ip.
* `wan_port` - Internet port.
* `wan_status` - Internet status.


## Import

cynosdb wan can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_wan.wan cynosdbmysql-bws8h88b#cynosdbmysql-grp-lxav0p9z
```

