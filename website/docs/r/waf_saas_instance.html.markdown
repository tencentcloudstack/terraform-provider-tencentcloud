---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_saas_instance"
sidebar_current: "docs-tencentcloud-resource-waf_saas_instance"
description: |-
  Provides a resource to create a waf saas instance
---

# tencentcloud_waf_saas_instance

Provides a resource to create a waf saas instance

~> **NOTE:** Region only supports `ap-guangzhou` and `ap-seoul`.

## Example Usage

### Create a basic waf premium saas instance

```hcl
resource "tencentcloud_waf_saas_instance" "example" {
  goods_category = "premium_saas"
  instance_name  = "tf-example-saas-waf"
}
```

### Create a complete waf ultimate_saas instance

```hcl
resource "tencentcloud_waf_saas_instance" "example" {
  goods_category  = "ultimate_saas"
  instance_name   = "tf-example-saas-waf"
  time_span       = 1
  time_unit       = "m"
  auto_renew_flag = 1
  elastic_mode    = 1
  real_region     = "gz"
}
```

### Set waf ultimate_saas instance qps limit

```hcl
resource "tencentcloud_waf_saas_instance" "example" {
  goods_category  = "ultimate_saas"
  instance_name   = "tf-example-saas-waf"
  time_span       = 1
  time_unit       = "m"
  auto_renew_flag = 1
  elastic_mode    = 1
  real_region     = "gz"
  qps_limit       = 200000
}
```

## Argument Reference

The following arguments are supported:

* `goods_category` - (Required, String) Billing order parameters. support premium_saas, enterprise_saas, ultimate_saas.
* `auto_renew_flag` - (Optional, Int) Auto renew flag, 1: enable, 0: disable.
* `elastic_mode` - (Optional, Int) Is elastic billing enabled, 1: enable, 0: disable.
* `instance_name` - (Optional, String) Waf instance name.
* `qps_limit` - (Optional, Int) QPS Limit, Minimum setting 10000. Only `elastic_mode` is 1, can be set.
* `real_region` - (Optional, String) region. If Region is `ap-guangzhou`, support: gz, sh, bj, cd (Means: GuangZhou, ShangHai, BeiJing, ChengDu); If Region is `ap-seoul`, support: hk, sg, th, kr, in, de, ca, use, sao, usw, jkt (Means: HongKong, Singapore, Bandkok, Seoul, Mumbai, Frankfurt, Toronto, Virginia, SaoPaulo, SiliconValley, Jakarta).
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


