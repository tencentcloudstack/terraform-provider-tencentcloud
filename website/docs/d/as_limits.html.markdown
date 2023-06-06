---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_limits"
sidebar_current: "docs-tencentcloud-datasource-as_limits"
description: |-
  Use this data source to query detailed information of as limits
---

# tencentcloud_as_limits

Use this data source to query detailed information of as limits

## Example Usage

```hcl
data "tencentcloud_as_limits" "limits" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `max_number_of_auto_scaling_groups` - Maximum number of auto scaling groups allowed for creation by the user account.
* `max_number_of_launch_configurations` - Maximum number of launch configurations allowed for creation by the user account.
* `number_of_auto_scaling_groups` - Current number of auto scaling groups under the user account.
* `number_of_launch_configurations` - Current number of launch configurations under the user account.


