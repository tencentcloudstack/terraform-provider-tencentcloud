---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_store_location_config"
sidebar_current: "docs-tencentcloud-resource-dlc_store_location_config"
description: |-
  Provides a resource to create a dlc store_location_config
---

# tencentcloud_dlc_store_location_config

Provides a resource to create a dlc store_location_config

## Example Usage

```hcl
resource "tencentcloud_dlc_store_location_config" "store_location_config" {
  store_location = "cosn://cos-xxxxx-xxx/test/"
}
```

## Argument Reference

The following arguments are supported:

* `store_location` - (Required, String) Calculate the results of the COS path, such as: cosn: // bucketName/.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dlc store_location_config can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_store_location_config.store_location_config store_location_config_id
```

