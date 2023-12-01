---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_image_share_permission"
sidebar_current: "docs-tencentcloud-datasource-cvm_image_share_permission"
description: |-
  Use this data source to query detailed information of cvm image_share_permission
---

# tencentcloud_cvm_image_share_permission

Use this data source to query detailed information of cvm image_share_permission

## Example Usage

```hcl
data "tencentcloud_cvm_image_share_permission" "image_share_permission" {
  image_id = "img-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required, String) The ID of the image to be shared.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `share_permission_set` - Information on image sharing.
  * `account_id` - ID of the account with which the image is shared.
  * `created_time` - Time when an image was shared.


