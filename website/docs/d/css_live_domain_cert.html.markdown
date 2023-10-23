---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_live_domain_cert"
sidebar_current: "docs-tencentcloud-datasource-css_live_domain_cert"
description: |-
  Use this data source to query detailed information of css live_domain_cert
---

# tencentcloud_css_live_domain_cert

Use this data source to query detailed information of css live_domain_cert

## Example Usage

```hcl
data "tencentcloud_css_live_domain_cert" "live_domain_cert" {
  domain_name = "your_domain_name"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Playback domain name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domain_cert_info` - Certificate information.
  * `cert_domains` - List of domain names in the certificate.[*.x.com] for example.Note: this field may return `null`, indicating that no valid values can be obtained.
  * `cert_expire_time` - The certificate expiration time in UTC format.Note: Beijing time (UTC+8) is used.
  * `cert_id` - Certificate ID.
  * `cert_name` - Certificate name.
  * `cert_type` - Certificate type.0: user-added certificate1: Tencent Cloud-hosted certificate.
  * `cloud_cert_id` - Tencent Cloud SSL certificate ID.Note: this field may return `null`, indicating that no valid values can be obtained.
  * `create_time` - The creation time in UTC format.Note: Beijing time (UTC+8) is used.
  * `description` - Description.
  * `domain_name` - Domain name that uses this certificate.
  * `https_crt` - Certificate content.
  * `status` - Certificate status.


