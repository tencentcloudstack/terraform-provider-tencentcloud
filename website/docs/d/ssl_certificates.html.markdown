---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_certificates"
sidebar_current: "docs-tencentcloud-datasource-ssl_certificates"
description: |-
  Use this data source to query SSL certificate.
---

# tencentcloud_ssl_certificates

Use this data source to query SSL certificate.

## Example Usage

```hcl
data "tencentcloud_ssl_certificates" "foo" {
  name = "certificate"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) ID of the SSL certificate to be queried.
* `name` - (Optional, String) Name of the SSL certificate to be queried.
* `result_output_file` - (Optional, String) Used to save results.
* `type` - (Optional, String) Type of the SSL certificate to be queried. Available values includes: `CA` and `SVR`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `certificates` - An information list of certificate. Each element contains the following attributes:
  * `begin_time` - Beginning time of the SSL certificate.
  * `cert` - Content of the SSL certificate.
  * `create_time` - Creation time of the SSL certificate.
  * `domain` - Primary domain of the SSL certificate.
  * `dv_auths` - DV certification information.
    * `dv_auth_key` - DV authentication key.
    * `dv_auth_value` - DV authentication value.
    * `dv_auth_verify_type` - DV authentication type.
  * `end_time` - Ending time of the SSL certificate.
  * `id` - ID of the SSL certificate.
  * `key` - Key of the SSL certificate.
  * `name` - Name of the SSL certificate.
  * `order_id` - Order ID returned.
  * `product_zh_name` - Certificate authority.
  * `project_id` - Project ID of the SSL certificate.
  * `status` - Status of the SSL certificate.
  * `subject_names` - ALL domains included in the SSL certificate. Including the primary domain name.
  * `type` - Type of the SSL certificate.


