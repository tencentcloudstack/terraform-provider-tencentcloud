---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_modify_data_engine_description_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_modify_data_engine_description_operation"
description: |-
  Provides a resource to create a DLC modify data engine description operation
---

# tencentcloud_dlc_modify_data_engine_description_operation

Provides a resource to create a DLC modify data engine description operation

## Example Usage

```hcl
resource "tencentcloud_dlc_modify_data_engine_description_operation" "example" {
  data_engine_name = "tf-example"
  message          = "message info."
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_name` - (Required, String, ForceNew) Engine description and its maximum length is 250 characters.
* `message` - (Required, String, ForceNew) Engine description and its maximum length is 250 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



