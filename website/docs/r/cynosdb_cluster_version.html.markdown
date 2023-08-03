---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_version"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cluster_version"
description: |-
  Provides a resource to create a cynosdb cluster_version
---

# tencentcloud_cynosdb_cluster_version

Provides a resource to create a cynosdb cluster_version

## Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_version" "cluster_version" {
  cluster_id    = "cynosdbmysql-bws8h88b"
  cynos_version = "2.1.10"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `cynos_version` - (Required, String, ForceNew) Kernel version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



