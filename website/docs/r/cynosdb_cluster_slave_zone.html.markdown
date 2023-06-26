---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_slave_zone"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cluster_slave_zone"
description: |-
  Provides a resource to create a cynosdb cluster_slave_zone
---

# tencentcloud_cynosdb_cluster_slave_zone

Provides a resource to create a cynosdb cluster_slave_zone

## Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_slave_zone" "cluster_slave_zone" {
  cluster_id = "cynosdbmysql-xxxxxxxx"
  slave_zone = "ap-guangzhou-3"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) The ID of cluster.
* `slave_zone` - (Required, String) Slave zone.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cynosdb cluster_slave_zone can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone cluster_slave_zone_id
```

