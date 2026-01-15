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
* `upstream_cert_info` - (Optional, List) Configures the certificate presented by the EO node during origin-pull for mutual TLS authentication. Disabled by default; leaving the field blank will retain the current configuration. This feature is currently in beta testing. please [contact us](https://cloud.tencent.com/online-service) to request access.

The `cert_infos` object of `upstream_mutual_tls` supports the following:

* `cert_id` - (Required, String) Certificate ID, which originates from the SSL side. You can check the CertId from the [SSL Certificate List](https://console.cloud.tencent.com/ssl).

The `server_cert_info` object supports the following:

* `cert_id` - (Required, String) ID of the server certificate.Note: This field may return null, indicating that no valid values can be obtained.
* `alias` - (Optional, String) Alias of the certificate.Note: This field may return null, indicating that no valid values can be obtained.
* `common_name` - (Optional, String) Domain name of the certificate. Note: This field may return `null`, indicating that no valid value can be obtained.
* `deploy_time` - (Optional, String) Time when the certificate is deployed. Note: This field may return null, indicating that no valid values can be obtained.
* `expire_time` - (Optional, String) Time when the certificate expires. Note: This field may return null, indicating that no valid values can be obtained.
* `sign_algo` - (Optional, String) Signature algorithm. Note: This field may return null, indicating that no valid values can be obtained.
* `type` - (Optional, String) Type of the certificate. Values: `default`: Default certificate; `upload`: Specified certificate; `managed`: Tencent Cloud-managed certificate. Note: This field may return `null`, indicating that no valid value can be obtained.

The `upstream_cert_info` object supports the following:

* `upstream_mutual_tls` - (Optional, List) In the origin-pull mutual authentication scenario, this field represents the certificate (including the public and private keys) carried during EO node origin-pull, which is deployed in the EO node for the origin server to authenticate the EO node. When used as an input parameter, it is left blank to indicate retaining the original configuration.

The `upstream_mutual_tls` object of `upstream_cert_info` supports the following:

* `switch` - (Required, String) Mutual authentication configuration switch, the values are: `on`: enable; `off`: disable.
* `cert_infos` - (Optional, List) Mutual authentication certificate list.
Note: When using MutualTLS as an input parameter in ModifyHostsCertificate, you only need to provide the CertId of the corresponding certificate. You can check the CertId from the [SSL Certificate List](https://console.cloud.tencent.com/ssl).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.

## Import

teo certificate can be imported using the id, e.g.

```
terraform import tencentcloud_teo_certificate_config.certificate zone_id#host
```

