---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_manager_detail"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_manager_detail"
description: |-
  Use this data source to query detailed information of ssl describe_manager_detail
---

# tencentcloud_ssl_describe_manager_detail

Use this data source to query detailed information of ssl describe_manager_detail

## Example Usage

```hcl
data "tencentcloud_ssl_describe_manager_detail" "describe_manager_detail" {
  manager_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `manager_id` - (Required, Int) Manager ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `company_id` - Manager Company ID.
* `company_info` - Manager&amp;#39;s company information.
  * `company_address` - Detailed address where the company is located.
  * `company_city` - The city where the company is.
  * `company_country` - Company country.
  * `company_id` - Company ID.
  * `company_name` - Company Name.
  * `company_phone` - company phone.
  * `company_province` - Province where the company is located.
  * `id_number` - ID numberNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `id_type` - typeNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `contact_first_name` - Contact name.
* `contact_last_name` - Contact name.
* `contact_mail` - Contact mailbox.
* `contact_phone` - contact number.
* `contact_position` - Contact position.
* `create_time` - Creation time.
* `expire_time` - Verify expiration timeNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `manager_department` - Administrator department.
* `manager_first_name` - Manager name.
* `manager_last_name` - Manager name.
* `manager_mail` - Manager mailbox.
* `manager_phone` - Manager phone call.
* `manager_position` - Manager position.
* `status` - Status: Audit: OK during the review: review passed inValid: expired expiRing: is about to expire Expired: expired.
* `verify_time` - Verify timeNote: This field may return NULL, indicating that the valid value cannot be obtained.


