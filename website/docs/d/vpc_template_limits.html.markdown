---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_template_limits"
sidebar_current: "docs-tencentcloud-datasource-vpc_template_limits"
description: |-
  Use this data source to query detailed information of vpc template_limits
---

# tencentcloud_vpc_template_limits

Use this data source to query detailed information of vpc template_limits

## Example Usage

```hcl
data "tencentcloud_vpc_template_limits" "template_limits" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `template_limit` - template limit.
  * `address_template_group_member_limit` - address template group member limit.
  * `address_template_member_limit` - address template member limit.
  * `service_template_group_member_limit` - service template group member limit.
  * `service_template_member_limit` - service template member limit.


