---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_ssl"
sidebar_current: "docs-tencentcloud-resource-cynosdb_ssl"
description: |-
  Provides a resource to create a cynosdb ssl
---

# tencentcloud_cynosdb_ssl

Provides a resource to create a cynosdb ssl

## Example Usage

```hcl
resource "tencentcloud_cynosdb_ssl" "cynosdb_ssl" {
  cluster_id  = "cynosdbmysql-1e0nzayx"
  instance_id = "cynosdbmysql-ins-pfsv6q1e"
  status      = "ON"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster id.
* `instance_id` - (Required, String) instance id.
* `status` - (Required, String) Whether to enable SSL. `ON` means enabled, `OFF` means not enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `download_url` - Certificate download address.


## Import

cynosdb ssl can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_ssl.cynosdb_ssl ${cluster_id}#${instance_id}
```

