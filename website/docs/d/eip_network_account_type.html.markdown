---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip_network_account_type"
sidebar_current: "docs-tencentcloud-datasource-eip_network_account_type"
description: |-
  Use this data source to query detailed information of eip network_account_type
---

# tencentcloud_eip_network_account_type

Use this data source to query detailed information of eip network_account_type

## Example Usage

```hcl
data "tencentcloud_eip_network_account_type" "network_account_type" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `network_account_type` - The network type of the user account, STANDARD is a standard user, LEGACY is a traditional user.


