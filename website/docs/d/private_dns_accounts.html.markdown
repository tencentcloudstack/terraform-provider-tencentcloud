---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_accounts"
sidebar_current: "docs-tencentcloud-datasource-private_dns_accounts"
description: |-
  Use this data source to query detailed information of privatedns accounts
---

# tencentcloud_private_dns_accounts

Use this data source to query detailed information of privatedns accounts

## Example Usage

### Query all accounts

```hcl
data "tencentcloud_private_dns_accounts" "example" {}
```

### Query accounts by filters

```hcl
data "tencentcloud_private_dns_accounts" "example" {
  filters {
    name   = "AccountUin"
    values = ["100022770160"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter parameters.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Parameter name.
* `values` - (Required, Set) Array of parameter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `account_set` - List of Private DNS accounts.


