---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_store_location_config"
sidebar_current: "docs-tencentcloud-resource-dlc_store_location_config"
description: |-
  Provides a resource to create a DLC store location config
---

# tencentcloud_dlc_store_location_config

Provides a resource to create a DLC store location config

## Example Usage

### Select user-defined COS path storage

```hcl
resource "tencentcloud_dlc_store_location_config" "example" {
  store_location = "cosn://tf-example-1308135196/demo"
  enable         = 1
}
```

### Select DLC internal storage

```hcl
resource "tencentcloud_dlc_store_location_config" "example" {
  store_location = ""
  enable         = 0
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Int) Whether to enable advanced settings. 0 means no while 1 means yes.
* `store_location` - (Required, String) The calculation results are stored in the cos path, such as: cosn://bucketname/.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



