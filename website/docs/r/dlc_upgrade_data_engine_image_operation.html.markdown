---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_upgrade_data_engine_image_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_upgrade_data_engine_image_operation"
description: |-
  Provides a resource to create a DLC upgrade data engine image operation
---

# tencentcloud_dlc_upgrade_data_engine_image_operation

Provides a resource to create a DLC upgrade data engine image operation

## Example Usage

```hcl
resource "tencentcloud_dlc_upgrade_data_engine_image_operation" "example" {
  data_engine_id = "DataEngine-g5ds87d8"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String, ForceNew) Engine ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [DescribeDataEngine](https://www.tencentcloud.com/document/api/1155/62693)
  - [DescribeDataEngines](https://www.tencentcloud.com/search/?wd=DescribeDataEngines&type=document&lang=en)
  - [UpgradeDataEngineImage](https://www.tencentcloud.com/document/api/1155/62643)

