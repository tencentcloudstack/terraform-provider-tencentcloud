---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_certificate_config"
sidebar_current: "docs-tencentcloud-resource-teo_certificate_config"
description: |-
  Provides a resource to create a TEO certificate config
---

# tencentcloud_teo_certificate_config

Provides a resource to create a TEO certificate config

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

### Configure SSL certificate with edge mutual TLS

```hcl
resource "tencentcloud_teo_certificate_config" "certificate" {
  host    = "test.tencentcloud-terraform-provider.cn"
  mode    = "sslcert"
  zone_id = "zone-2o1t24kgy362"

  server_cert_info {
    cert_id = "8xiUJIJd"
  }

  client_cert_info {
    switch = "on"
    cert_infos {
      cert_id = "cert-client-001"
    }
  }
}
```

### Configure SSL certificate with upstream mutual TLS

```hcl
resource "tencentcloud_teo_certificate_config" "certificate" {
  host    = "test.tencentcloud-terraform-provider.cn"
  mode    = "sslcert"
  zone_id = "zone-2o1t24kgy362"

  server_cert_info {
    cert_id = "8xiUJIJd"
  }

  upstream_cert_info {
    upstream_mutual_tls {
      switch = "on"
      cert_infos {
        cert_id = "cert-upstream-001"
      }
    }
  }
}
```

### Configure SSL certificate with upstream certificate verification

```hcl
resource "tencentcloud_teo_certificate_config" "certificate" {
  host    = "test.tencentcloud-terraform-provider.cn"
  mode    = "sslcert"
  zone_id = "zone-2o1t24kgy362"

  server_cert_info {
    cert_id = "8xiUJIJd"
  }

  upstream_cert_info {
    upstream_certificate_verify {
      verification_mode = "custom_ca"
      custom_ca_certs {
        cert_id = "cert-ca-001"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String, ForceNew) Acceleration domain name that needs to modify the certificate configuration.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `client_cert_info` - (Optional, List) Edge mutual TLS authentication configuration, where client CA certificates are deployed on EO nodes for client-to-EO-node authentication. Disabled by default; leaving the field blank will retain the current configuration. This feature is currently in beta testing. please [contact us](https://cloud.tencent.com/online-service) to request access.
* `mode` - (Optional, String) Mode of configuring the certificate, the values are: `disable`: Do not configure the certificate; `eofreecert`: Configure EdgeOne free certificate; `eofreecert_manual`: Deploy a free certificate applied for through DNS delegation validation or file validation; `sslcert`: Configure SSL certificate. If not filled in, the default value is `disable`.
* `server_cert_info` - (Optional, List) SSL certificate configuration, this parameter takes effect only when mode = sslcert, just enter the corresponding CertId. You can go to the SSL certificate list to view the CertId.
* `upstream_cert_info` - (Optional, List) Configures the certificate presented by the EO node during origin-pull for mutual TLS authentication. Disabled by default; leaving the field blank will retain the current configuration. This feature is currently in beta testing. please [contact us](https://cloud.tencent.com/online-service) to request access.

The `cert_infos` object of `client_cert_info` supports the following:

* `cert_id` - (Optional, String) Certificate ID, which originates from the SSL side. You can check the CertId from the [SSL Certificate List](https://console.cloud.tencent.com/ssl).

The `cert_infos` object of `upstream_mutual_tls` supports the following:

* `cert_id` - (Optional, String) Certificate ID, which originates from the SSL side. You can check the CertId from the [SSL Certificate List](https://console.cloud.tencent.com/ssl).

The `client_cert_info` object supports the following:

* `switch` - (Required, String) Edge mutual TLS configuration switch, the values are: `on`: enable; `off`: disable.
* `cert_infos` - (Optional, List) Mutual TLS certificate list.
Note: When using ClientCertInfo as an input parameter in ModifyHostsCertificate, you only need to provide the CertId of the corresponding certificate. You can check the CertId from the [SSL Certificate List](https://console.cloud.tencent.com/ssl).

The `custom_ca_certs` object of `upstream_certificate_verify` supports the following:

* `cert_id` - (Optional, String) Certificate ID, which originates from the SSL side.

The `server_cert_info` object supports the following:

* `cert_id` - (Required, String) ID of the server certificate.Note: This field may return null, indicating that no valid values can be obtained.
* `alias` - (Optional, String) Alias of the certificate.Note: This field may return null, indicating that no valid values can be obtained.
* `common_name` - (Optional, String) Domain name of the certificate. Note: This field may return `null`, indicating that no valid value can be obtained.
* `deploy_time` - (Optional, String) Time when the certificate is deployed. Note: This field may return null, indicating that no valid values can be obtained.
* `expire_time` - (Optional, String) Time when the certificate expires. Note: This field may return null, indicating that no valid values can be obtained.
* `sign_algo` - (Optional, String) Signature algorithm. Note: This field may return null, indicating that no valid values can be obtained.
* `type` - (Optional, String) Type of the certificate. Values: `default`: Default certificate; `upload`: Specified certificate; `managed`: Tencent Cloud-managed certificate. Note: This field may return `null`, indicating that no valid value can be obtained.

The `upstream_cert_info` object supports the following:

* `upstream_certificate_verify` - (Optional, List) In the origin certificate verification scenario, this field is the CA certificate used by EO nodes during origin-pull for verifying the origin server's certificate. Deployed on EO nodes for EO to authenticate the server certificate. When used as an input parameter, leaving it blank means retaining the original configuration.
* `upstream_mutual_tls` - (Optional, List) In the origin-pull mutual authentication scenario, this field represents the certificate (including the public and private keys) carried during EO node origin-pull, which is deployed in the EO node for the origin server to authenticate the EO node. When used as an input parameter, it is left blank to indicate retaining the original configuration.

The `upstream_certificate_verify` object of `upstream_cert_info` supports the following:

* `custom_ca_certs` - (Optional, List) List of specified trusted CA certificates. The origin certificate must be signed by this CA to pass verification.
Note: Only required when VerificationMode is custom_ca. When used as input in ModifyHostsCertificate, you only need to provide the CertId.
* `verification_mode` - (Optional, String) Origin certificate verification mode. Values: `disable`: Disable origin certificate verification; `custom_ca`: Use specified trusted CA certificate for verification.

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

TEO certificate config can be imported using the id, e.g.

```
terraform import tencentcloud_teo_certificate_config.certificate zone_id#host
```

