---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_certificates"
sidebar_current: "docs-tencentcloud-datasource-gaap_certificates"
description: |-
  Use this data source to query GAAP certificate.
---

# tencentcloud_gaap_certificates

Use this data source to query GAAP certificate.

## Example Usage

```hcl
resource "tencentcloud_gaap_certificate" "foo" {
  type    = "BASIC"
  content = "test:tx2KGdo3zJg/."
  name    = "test_certificate"
}

data "tencentcloud_gaap_certificates" "foo" {
  id = tencentcloud_gaap_certificate.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) ID of the certificate to be queried.
* `name` - (Optional, String) Name of the certificate to be queried.
* `result_output_file` - (Optional, String) Used to save results.
* `type` - (Optional, String) Type of the certificate to be queried. Valid values: `BASIC`, `CLIENT`, `SERVER`, `REALSERVER` and `PROXY`. `BASIC` means basic certificate; `CLIENT` means client CA certificate; `SERVER` means server SSL certificate; `REALSERVER` means realserver CA certificate; `PROXY` means proxy SSL certificate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `certificates` - An information list of certificate. Each element contains the following attributes:
  * `begin_time` - Beginning time of the certificate.
  * `create_time` - Creation time of the certificate.
  * `end_time` - Ending time of the certificate.
  * `id` - ID of the certificate.
  * `issuer_cn` - Issuer name of the certificate.
  * `name` - Name of the certificate.
  * `subject_cn` - Subject name of the certificate.
  * `type` - Type of the certificate.


