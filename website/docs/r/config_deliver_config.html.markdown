---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_deliver_config"
sidebar_current: "docs-tencentcloud-resource-config_deliver_config"
description: |-
  Provides a resource to manage Config delivery settings (global singleton configuration).
---

# tencentcloud_config_deliver_config

Provides a resource to manage Config delivery settings (global singleton configuration).

## Example Usage

```hcl
resource "tencentcloud_config_deliver_config" "example" {
  status               = 1
  deliver_name         = "tf-example-deliver"
  target_arn           = "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"
  deliver_prefix       = "config"
  deliver_type         = "COS"
  deliver_content_type = 3
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Required, Int) Delivery switch. Valid values: 0 (disabled), 1 (enabled).
* `deliver_content_type` - (Optional, Int) Delivery content type. Valid values: 1 (configuration change), 2 (resource list), 3 (all).
* `deliver_name` - (Optional, String) Delivery service name.
* `deliver_prefix` - (Optional, String) Log prefix for stored delivery content.
* `deliver_type` - (Optional, String) Delivery type. Valid values: COS, CLS.
* `target_arn` - (Optional, String) Resource ARN. COS format: qcs::cos:$region:$account:prefix/$appid/$BucketName. CLS format: qcs::cls:$region:$account:cls/topicId.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the delivery configuration.


