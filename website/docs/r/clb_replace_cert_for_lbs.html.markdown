---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_replace_cert_for_lbs"
sidebar_current: "docs-tencentcloud-resource-clb_replace_cert_for_lbs"
description: |-
  Provides a resource to create a clb replace_cert_for_lbs
---

# tencentcloud_clb_replace_cert_for_lbs

Provides a resource to create a clb replace_cert_for_lbs

## Example Usage

```hcl
resource "tencentcloud_clb_replace_cert_for_lbs" "replace_cert_for_lbs" {
  old_certificate_id = "zjUMifFK"
  certificate {
    cert_ca_name    = "test"
    cert_ca_content = "XXXXX"
  }
}
```
```hcl
terraform import tencentcloud_clb_replace_cert_for_lbs.replace_cert_for_lbs replace_cert_for_lbs_id
```

## Argument Reference

The following arguments are supported:

* `certificate` - (Required, List, ForceNew) Information such as the content of the new certificate.
* `old_certificate_id` - (Required, String, ForceNew) ID of the certificate to be replaced, which can be a server certificate or a client certificate.

The `certificate` object supports the following:

* `cert_ca_content` - (Optional, String) Content of the uploaded client certificate. When SSLMode = mutual, if there is no CertCaId, this parameter is required.
* `cert_ca_id` - (Optional, String) ID of a client certificate. When the listener adopts mutual authentication (i.e., SSLMode = mutual), if you leave this parameter empty, you must upload the client certificate, including CertCaContent and CertCaName.
* `cert_ca_name` - (Optional, String) Name of the uploaded client CA certificate. When SSLMode = mutual, if there is no CertCaId, this parameter is required.
* `cert_content` - (Optional, String) Content of the uploaded server certificate. If there is no CertId, this parameter is required.
* `cert_id` - (Optional, String) ID of a server certificate. If you leave this parameter empty, you must upload the certificate, including CertContent, CertKey, and CertName.
* `cert_key` - (Optional, String) Key of the uploaded server certificate. If there is no CertId, this parameter is required.
* `cert_name` - (Optional, String) Name of the uploaded server certificate. If there is no CertId, this parameter is required.
* `ssl_mode` - (Optional, String) Authentication type. Value range: UNIDIRECTIONAL (unidirectional authentication), MUTUAL (mutual authentication).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



