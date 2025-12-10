---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_transparent_encrypt"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cluster_transparent_encrypt"
description: |-
  Provides a resource to create a Cynosdb cluster transparent encrypt
---

# tencentcloud_cynosdb_cluster_transparent_encrypt

Provides a resource to create a Cynosdb cluster transparent encrypt

~> **NOTE:** Once activated, it cannot be deactivated.

~> **NOTE:** If you have not enabled the KMS service or authorized the KMS key before, you will need to enable the KMS service and then authorize the KMS key in order to complete the corresponding enabling or authorization operations and unlock the subsequent settings for data encryption.

## Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_transparent_encrypt" "example" {
  cluster_id                = cynosdbmysql-bu6hlulf
  key_id                    = "f063c18b-xxxx-xxxx-xxxx-525400d3a886"
  key_region                = "ap-guangzhou"
  key_type                  = "custom"
  is_open_global_encryption = false
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `key_type` - (Required, String) Key type (cloud, custom).
* `is_open_global_encryption` - (Optional, Bool) Whether to enable global encryption.
* `key_id` - (Optional, String) Key Id.
* `key_region` - (Optional, String) Key region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Cynosdb cluster transparent encrypt can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_transparent_encrypt.example cynosdbmysql-bu6hlulf
```

