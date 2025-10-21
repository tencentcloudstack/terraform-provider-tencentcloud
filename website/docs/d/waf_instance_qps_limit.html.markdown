---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_instance_qps_limit"
sidebar_current: "docs-tencentcloud-datasource-waf_instance_qps_limit"
description: |-
  Use this data source to query detailed information of waf instance_qps_limit
---

# tencentcloud_waf_instance_qps_limit

Use this data source to query detailed information of waf instance_qps_limit

## Example Usage

```hcl
data "tencentcloud_waf_instance_qps_limit" "example" {
  instance_id = "waf_2kxtlbky00b3b4qz"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Unique ID of Instance.
* `result_output_file` - (Optional, String) Used to save results.
* `type` - (Optional, String) Instance type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `qps_data` - Qps info.
  * `elastic_billing_default` - Elastic qps default value.
  * `elastic_billing_max` - Maximum elastic qps.
  * `elastic_billing_min` - Minimum elastic qps.
  * `qps_extend_intl_max` - Maximum qps of extend package for overseas.
  * `qps_extend_max` - Maximum qps of extend package.


