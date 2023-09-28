---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_certificate_config"
sidebar_current: "docs-tencentcloud-resource-teo_certificate_config"
description: |-
  Provides a resource to create a teo certificate
---

# tencentcloud_teo_certificate_config

Provides a resource to create a teo certificate

## Example Usage

```hcl
resource "tencentcloud_teo_certificate_config" "certificate" {
  host    = "test.tencentcloud-terraform-provider.cn"
  mode    = "eofreecert"
  zone_id = "zone-2o1t24kgy362"
}
```

### Configure SSL certificate

```hcl
resource "tencentcloud_teo_certificate_config" "certificate" {
  host    = "test.tencentcloud-terraform-provider.cn"
  mode    = "sslcert"
  zone_id = "zone-2o1t24kgy362"

  server_cert_info {
    cert_id = "8xiUJIJd"
  }
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String, ForceNew) Acceleration domain name that needs to modify the certificate configuration.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `mode` - (Optional, String) Mode of configuring the certificate, the values are: `disable`: Do not configure the certificate; `eofreecert`: Configure EdgeOne free certificate; `sslcert`: Configure SSL certificate. If not filled in, the default value is `disable`.
* `server_cert_info` - (Optional, List) SSL certificate configuration, this parameter takes effect only when mode = sslcert, just enter the corresponding CertId. You can go to the SSL certificate list to view the CertId.

The `server_cert_info` object supports the following:

* `cert_id` - (Required, String) ID of the server certificate.Note: This field may return null, indicating that no valid values can be obtained.
* `alias` - (Optional, String) Alias of the certificate.Note: This field may return null, indicating that no valid values can be obtained.
* `common_name` - (Optional, String) Domain name of the certificate. Note: This field may return `null`, indicating that no valid value can be obtained.
* `deploy_time` - (Optional, String) Time when the certificate is deployed. Note: This field may return null, indicating that no valid values can be obtained.
* `expire_time` - (Optional, String) Time when the certificate expires. Note: This field may return null, indicating that no valid values can be obtained.
* `sign_algo` - (Optional, String) Signature algorithm. Note: This field may return null, indicating that no valid values can be obtained.
* `type` - (Optional, String) Type of the certificate. Values: `default`: Default certificate; `upload`: Specified certificate; `managed`: Tencent Cloud-managed certificate. Note: This field may return `null`, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo certificate can be imported using the id, e.g.

```
terraform import tencentcloud_teo_certificate_config.certificate zone_id#host#cert_id
```

