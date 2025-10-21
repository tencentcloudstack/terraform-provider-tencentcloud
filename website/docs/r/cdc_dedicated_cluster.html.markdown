---
subcategory: "CDC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdc_dedicated_cluster"
sidebar_current: "docs-tencentcloud-resource-cdc_dedicated_cluster"
description: |-
  Provides a resource to create a CDC dedicated cluster
---

# tencentcloud_cdc_dedicated_cluster

Provides a resource to create a CDC dedicated cluster

## Example Usage

```hcl
# create cdc site
resource "tencentcloud_cdc_site" "example" {
  name         = "tf-example"
  country      = "China"
  province     = "Guangdong Province"
  city         = "Guangzhou"
  address_line = "Tencent Building"
  description  = "desc."
}

# create cdc dedicated cluster
resource "tencentcloud_cdc_dedicated_cluster" "example" {
  site_id     = tencentcloud_cdc_site.example.id
  name        = "tf-example"
  zone        = "ap-guangzhou-6"
  description = "desc."
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Dedicated Cluster Name.
* `site_id` - (Required, String) Dedicated Cluster Site ID.
* `zone` - (Required, String) Dedicated Cluster Zone.
* `description` - (Optional, String) Dedicated Cluster Description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CDC dedicated cluster can be imported using the id, e.g.

```
terraform import tencentcloud_cdc_dedicated_cluster.example cluster-d574omhk
```

