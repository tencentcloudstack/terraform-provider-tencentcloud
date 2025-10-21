---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_rotation_history"
sidebar_current: "docs-tencentcloud-datasource-ssm_rotation_history"
description: |-
  Use this data source to query detailed information of ssm rotation_history
---

# tencentcloud_ssm_rotation_history

Use this data source to query detailed information of ssm rotation_history

## Example Usage

```hcl
data "tencentcloud_ssm_rotation_history" "example" {
  secret_name = "keep_terraform"
}
```

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, String) Secret name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version_ids` - The number of version numbers. The maximum number of version numbers that can be displayed to users is 10.


