---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_update_aggregate_config_deliver"
sidebar_current: "docs-tencentcloud-resource-config_update_aggregate_config_deliver"
description: |-
  Provides a resource to manage the aggregate delivery settings (投递设置) of a Tencent Cloud Config account group (账号组).
---

# tencentcloud_config_update_aggregate_config_deliver

Provides a resource to manage the aggregate delivery settings (投递设置) of a Tencent Cloud Config account group (账号组).

## Example Usage

```hcl
resource "tencentcloud_config_update_aggregate_config_deliver" "example" {
  account_group_id     = "ca-ag-xxxxxx"
  status               = 1
  deliver_name         = "tf-example-aggregate-deliver"
  target_arn           = "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"
  deliver_prefix       = "config"
  deliver_type         = "COS"
  deliver_uin          = 0
  deliver_content_type = 3
}
```

## Argument Reference

The following arguments are supported:

* `account_group_id` - (Required, String, ForceNew) Account group ID.
* `status` - (Required, Int) Delivery switch. Valid values: 0 (disabled), 1 (enabled).
* `deliver_content_type` - (Optional, Int) Delivery content type. Valid values: 1 (configuration change), 2 (resource list), 3 (all).
* `deliver_name` - (Optional, String) Delivery service name.
* `deliver_prefix` - (Optional, String) Log prefix for stored delivery content.
* `deliver_type` - (Optional, String) Delivery type. Valid values: COS, CLS.
* `deliver_uin` - (Optional, Int) Member account uin that supports cross-account delivery. Only the delegated administrator can be used. The default value is 0, which means delivery to the administrator account.
* `target_arn` - (Optional, String) Resource ARN. COS format: qcs::cos:$region:$account:prefix/$appid/$BucketName. CLS format: qcs::cls:$region:$account:cls/topicId.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the delivery configuration.


## Import

Config aggregate deliver can be imported using the account_group_id, e.g.

```shell
terraform import tencentcloud_config_update_aggregate_config_deliver.example ca-ag-xxxxxx
```

