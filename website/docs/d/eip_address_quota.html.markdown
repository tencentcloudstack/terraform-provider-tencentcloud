---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip_address_quota"
sidebar_current: "docs-tencentcloud-datasource-eip_address_quota"
description: |-
  Use this data source to query detailed information of vpc address_quota
---

# tencentcloud_eip_address_quota

Use this data source to query detailed information of vpc address_quota

## Example Usage

```hcl
data "tencentcloud_eip_address_quota" "address_quota" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `quota_set` - The specified account EIP quota information.
  * `quota_current` - Current count.
  * `quota_id` - Quota name: TOTAL_EIP_QUOTA,DAILY_EIP_APPLY,DAILY_PUBLIC_IP_ASSIGN.
  * `quota_limit` - quota count.


