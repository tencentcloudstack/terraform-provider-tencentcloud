---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_resource_packages_attachment"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cluster_resource_packages_attachment"
description: |-
  Provides a resource to create a cynosdb cluster_resource_packages_attachment
---

# tencentcloud_cynosdb_cluster_resource_packages_attachment

Provides a resource to create a cynosdb cluster_resource_packages_attachment

## Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_resource_packages_attachment" "cluster_resource_packages_attachment" {
  cluster_id  = "cynosdbmysql-q1d8151n"
  package_ids = ["package-hy4d2ppl"]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `package_ids` - (Required, Set: [`String`], ForceNew) Resource Package Unique ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cynosdb cluster_resource_packages_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_resource_packages_attachment.cluster_resource_packages_attachment cluster_resource_packages_attachment_id
```

