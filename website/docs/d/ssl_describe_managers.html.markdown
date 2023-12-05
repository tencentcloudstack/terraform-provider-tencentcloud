---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_managers"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_managers"
description: |-
  Use this data source to query detailed information of ssl describe_managers
---

# tencentcloud_ssl_describe_managers

Use this data source to query detailed information of ssl describe_managers

## Example Usage

```hcl
data "tencentcloud_ssl_describe_managers" "describe_managers" {
  company_id = "11772"
}
```

## Argument Reference

The following arguments are supported:

* `company_id` - (Required, Int) Company ID.
* `manager_mail` - (Optional, String) Vague query manager email (will be abandoned), please use Searchkey.
* `manager_name` - (Optional, String) Manager&amp;#39;s name (will be abandoned), please use Searchkey.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Manager&amp;#39;s surname/Manager name/mailbox/department precise matching.
* `status` - (Optional, String) Filter according to the status of the manager, and the value is available&amp;#39;None&amp;#39; Unable to submit review&amp;#39;Audit&amp;#39;, Asian Credit Review&amp;#39;Caaudit&amp;#39; CA review&amp;#39;OK&amp;#39; has been reviewed&amp;#39;Invalid&amp;#39; review failed&amp;#39;Expiring&amp;#39; is about to expire&amp;#39;Expired&amp;#39; expired.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `managers` - Company Manager List.
  * `cert_count` - Number of administrative certificates.
  * `create_time` - Creation timeNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `domain_count` - Number of administrators.
  * `expire_time` - Examine the validity expiration timeNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `manager_department` - Administrator department.
  * `manager_first_name` - Manager name.
  * `manager_id` - Manager ID.
  * `manager_last_name` - Manager name.
  * `manager_mail` - Manager mailbox.
  * `manager_phone` - Manager phone call.
  * `manager_position` - Manager position.
  * `status` - Status: Audit: OK during the review: review passed inValid: expired expiRing: is about to expire Expired: expired.
  * `submit_audit_time` - The last time the review timeNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `verify_time` - Examination timeNote: This field may return NULL, indicating that the valid value cannot be obtained.


