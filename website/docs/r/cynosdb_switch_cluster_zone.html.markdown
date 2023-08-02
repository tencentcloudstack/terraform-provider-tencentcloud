---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_switch_cluster_zone"
sidebar_current: "docs-tencentcloud-resource-cynosdb_switch_cluster_zone"
description: |-
  Provides a resource to create a cynosdb switch_cluster_zone
---

# tencentcloud_cynosdb_switch_cluster_zone

Provides a resource to create a cynosdb switch_cluster_zone

## Example Usage

```hcl
resource "tencentcloud_cynosdb_switch_cluster_zone" "switch_cluster_zone" {
  cluster_id = "cynosdbmysql-507j6phr"
  zone       = "ap-guangzhou-6"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster Id.
* `zone` - (Required, String) Availability zone to switch to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cynosdb switch_cluster_zone can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_switch_cluster_zone.switch_cluster_zone switch_cluster_zone_id
```

