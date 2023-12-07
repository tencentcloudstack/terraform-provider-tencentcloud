---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_companies"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_companies"
description: |-
  Use this data source to query detailed information of ssl describe_companies
---

# tencentcloud_ssl_describe_companies

Use this data source to query detailed information of ssl describe_companies

## Example Usage

```hcl
data "tencentcloud_ssl_describe_companies" "describe_companies" {
  company_id = 122
}
```

## Argument Reference

The following arguments are supported:

* `company_id` - (Optional, Int) Company ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `companies` - Company list.
  * `company_address` - Detailed address where the company is located.
  * `company_city` - The city where the company is.
  * `company_country` - Company country.
  * `company_id` - Company ID.
  * `company_name` - Company Name.
  * `company_phone` - company phone.
  * `company_province` - Province where the company is located.
  * `id_number` - ID numberNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `id_type` - typeNote: This field may return NULL, indicating that the valid value cannot be obtained.


