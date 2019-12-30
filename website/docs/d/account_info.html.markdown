---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_account_info"
sidebar_current: "docs-tencentcloud-datasource-account_info"
description: |-
  Use this data source to query account info from tencentcloud
---

# tencentcloud_account_info

Use this data source to query account info from tencentcloud

## Example Usage

```hcl
data "tencentcloud_account_info" "info" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `appid` - Appid in tencentcloud.


