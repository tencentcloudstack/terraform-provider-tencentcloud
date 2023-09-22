---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_clb_instance"
sidebar_current: "docs-tencentcloud-resource-waf_clb_instance"
description: |-
  Provides a resource to create a waf clb instance
---

# tencentcloud_waf_clb_instance

Provides a resource to create a waf clb instance

~> **NOTE:** Region only supports `ap-guangzhou` and `ap-seoul`.

## Example Usage

### Create a basic waf premium clb instance

```hcl
resource "tencentcloud_waf_clb_instance" "example" {
  goods_category = "premium_clb"
  instance_name  = "tf-example-clb-waf"
}
```

### Create a complete waf ultimate_clb instance

```hcl
resource "tencentcloud_waf_clb_instance" "example" {
  goods_category  = "ultimate_clb"
  instance_name   = "tf-example-clb-waf"
  time_span       = 1
  time_unit       = "m"
  auto_renew_flag = 1
  elastic_mode    = 1
}
```

## Argument Reference

The following arguments are supported:

* `goods_category` - (Required, String) Billing order parameters. support: premium_clb, enterprise_clb, ultimate_clb.
* `auto_renew_flag` - (Optional, Int) Auto renew flag, 1: enable, 0: disable.
* `elastic_mode` - (Optional, Int) Is elastic billing enabled, 1: enable, 0: disable.
* `instance_name` - (Optional, String) Waf instance name.
* `time_span` - (Optional, Int) Time interval.
* `time_unit` - (Optional, String) Time unit, support d, m, y. d: day, m: month, y: year.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_security` - waf instance api security status.
* `begin_time` - waf instance start time.
* `edition` - waf instance edition, clb or saas.
* `instance_id` - waf instance id.
* `status` - waf instance status.
* `valid_time` - waf instance valid time.


