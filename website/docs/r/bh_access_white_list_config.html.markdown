---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_access_white_list_config"
sidebar_current: "docs-tencentcloud-resource-bh_access_white_list_config"
description: |-
  Provides a resource to create a BH access white list config
---

# tencentcloud_bh_access_white_list_config

Provides a resource to create a BH access white list config

## Example Usage

```hcl
resource "tencentcloud_bh_access_white_list_config" "example" {
  allow_any  = false
  allow_auto = false
}
```

## Argument Reference

The following arguments are supported:

* `allow_any` - (Optional, Bool) true: allow all source IPs; false: do not allow all source IPs.
* `allow_auto` - (Optional, Bool) true: allow automatically added IPs; false: do not allow automatically added IPs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyAccessWhiteListAutoStatus](https://cloud.tencent.com/document/api/1025/125046) - 修改访问白名单自动添加IP状态
* Create: [ModifyAccessWhiteListStatus](https://cloud.tencent.com/document/api/1025/125044) - 修改访问白名单状态
* Read: [DescribeAccessWhiteListRules](https://cloud.tencent.com/document/api/1025/117903) - 查询访问白名单规则列表
* Update: [ModifyAccessWhiteListAutoStatus](https://cloud.tencent.com/document/api/1025/125046) - 修改访问白名单自动添加IP状态
* Update: [ModifyAccessWhiteListStatus](https://cloud.tencent.com/document/api/1025/125044) - 修改访问白名单状态

## Import

BH access white list config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_bh_access_white_list_config.example zDxkr768TFYadnFdX1fusQ==
```

