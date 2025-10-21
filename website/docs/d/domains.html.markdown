---
subcategory: "Domain"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_domains"
sidebar_current: "docs-tencentcloud-datasource-domains"
description: |-
  Provide a datasource to query Domains.
---

# tencentcloud_domains

Provide a datasource to query Domains.

## Example Usage

```hcl
data "tencentcloud_domains" "foo" {
}
```

## Argument Reference

The following arguments are supported:

* `limit` - (Optional, Int) Specify data limit in range [1, 100]. Default: 20.
* `offset` - (Optional, Int) Specify data offset. Default: 0.
* `result_output_file` - (Optional, String) Used for save response as file locally.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Domain result list.
  * `auto_renew` - Whether the domain auto renew, 0 - manual renew, 1 - auto renew.
  * `buy_status` - Domain buy status.
  * `code_tld` - Domain code ltd.
  * `creation_date` - Domain create time.
  * `domain_id` - Domain ID.
  * `domain_name` - Domain name.
  * `expiration_date` - Domain expiration date.
  * `is_premium` - Whether the domain is premium.
  * `tld` - Domain ltd.


