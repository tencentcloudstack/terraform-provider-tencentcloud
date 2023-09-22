---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_saas_domain"
sidebar_current: "docs-tencentcloud-resource-waf_saas_domain"
description: |-
  Provides a resource to create a waf saas_domain
---

# tencentcloud_waf_saas_domain

Provides a resource to create a waf saas_domain

## Example Usage

### If upstream_type is 0

Create a basic waf saas domain

```hcl
resource "tencentcloud_waf_saas_domain" "example" {
  instance_id = "waf_2kxtlbky01b3wceb"
  domain      = "tf.example.com"
  src_list = [
    "1.1.1.1"
  ]

  ports {
    port              = "80"
    protocol          = "http"
    upstream_port     = "80"
    upstream_protocol = "http"
  }
}
```

### Create a load balancing strategy is weighted polling saas domain

```hcl
resource "tencentcloud_waf_saas_domain" "example" {
  instance_id = "waf_2kxtlbky01b3wceb"
  domain      = "tf.example.com"
  src_list = [
    "1.1.1.1",
    "2.2.2.2"
  ]
  load_balance = "2"
  weights = [
    30,
    50
  ]

  ports {
    port              = "80"
    protocol          = "http"
    upstream_port     = "80"
    upstream_protocol = "http"
  }
}
```

### If upstream_type is 1

```hcl
resource "tencentcloud_waf_saas_domain" "example" {
  instance_id     = "waf_2kxtlbky01b3wceb"
  domain          = "tf.example.com"
  upstream_type   = 1
  upstream_domain = "test.com"

  ports {
    port              = "80"
    protocol          = "http"
    upstream_port     = "80"
    upstream_protocol = "http"
  }
}
```

### Create a waf saas domain with set Http&Https

```hcl
resource "tencentcloud_waf_saas_domain" "example" {
  instance_id     = "waf_2kxtlbky01b3wceb"
  domain          = "tf.example.com"
  is_cdn          = 3
  cert_type       = 2
  ssl_id          = "3a6B5y8v"
  load_balance    = "2"
  https_rewrite   = 1
  upstream_scheme = "https"
  src_list = [
    "1.1.1.1",
    "2.2.2.2"
  ]
  weights = [
    50,
    60
  ]

  ports {
    port              = "80"
    protocol          = "http"
    upstream_port     = "80"
    upstream_protocol = "http"
  }

  ports {
    port              = "443"
    protocol          = "https"
    upstream_port     = "443"
    upstream_protocol = "https"
  }

  ip_headers = [
    "headers_1",
    "headers_2",
    "headers_3",
  ]
}
```

### Create a complete waf saas domain

```hcl
resource "tencentcloud_waf_saas_domain" "example" {
  instance_id     = "waf_2kxtlbky01b3wceb"
  domain          = "tf.example.com"
  is_cdn          = 3
  cert_type       = 2
  ssl_id          = "3a6B5y8v"
  load_balance    = "2"
  https_rewrite   = 1
  is_http2        = 1
  upstream_scheme = "https"
  src_list = [
    "1.1.1.1",
    "2.2.2.2"
  ]
  weights = [
    50,
    60
  ]

  ports {
    port              = "80"
    protocol          = "http"
    upstream_port     = "80"
    upstream_protocol = "http"
  }

  ports {
    port              = "443"
    protocol          = "https"
    upstream_port     = "443"
    upstream_protocol = "https"
  }

  ip_headers = [
    "headers_1",
    "headers_2",
    "headers_3",
  ]

  is_keep_alive      = "1"
  active_check       = 1
  tls_version        = 3
  cipher_template    = 1
  proxy_read_timeout = 500
  proxy_send_timeout = 500
  sni_type           = 3
  sni_host           = "3.3.3.3"
  xff_reset          = 1
  bot_status         = 1
  api_safe_status    = 1
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain names that require defense.
* `instance_id` - (Required, String) Unique ID of Instance.
* `ports` - (Required, Set) This field needs to be set for multiple ports in the upstream server.
* `active_check` - (Optional, Int) Whether to enable active health detection, 0 represents disable and 1 represents enable.
* `api_safe_status` - (Optional, Int) Whether to enable api safe, 1 enable, 0 disable.
* `bot_status` - (Optional, Int) Whether to enable bot, 1 enable, 0 disable.
* `cert_type` - (Optional, Int) Certificate type, 0 represents no certificate, CertType=1 represents self owned certificate, and 2 represents managed certificate.
* `cert` - (Optional, String) Certificate content, When CertType=1, this parameter needs to be filled.
* `cipher_template` - (Optional, Int) Encryption Suite Template, 0:default  1:Universal template 2:Security template 3:Custom template.
* `ciphers` - (Optional, List: [`Int`]) Encryption Suite Information.
* `https_rewrite` - (Optional, Int) Whether redirect to https, 1 will redirect and 0 will not.
* `https_upstream_port` - (Optional, String) Upstream port for https, When listen ports has https port and UpstreamScheme is HTTP, the current field needs to be filled.
* `ip_headers` - (Optional, List: [`String`]) When is_cdn=3, this parameter needs to be filled in to indicate a custom header.
* `is_cdn` - (Optional, Int) Whether a proxy has been enabled before WAF, 0 no deployment, 1 deployment and use first IP in X-Forwarded-For as client IP, 2 deployment and use remote_addr as client IP, 3 deployment and use values of custom headers as client IP.
* `is_http2` - (Optional, Int) Whether enable HTTP2, Enabling HTTP2 requires HTTPS support, 1 means enabled, 0 does not.
* `is_keep_alive` - (Optional, String) Whether to enable keep-alive, 0 disable, 1 enable.
* `is_websocket` - (Optional, Int) Is WebSocket support enabled. 1 means enabled, 0 does not.
* `load_balance` - (Optional, String) Load balancing strategy, where 0 represents polling and 1 represents IP hash and 2 weighted round robin.
* `private_key` - (Optional, String) Certificate key, When CertType=1, this parameter needs to be filled.
* `proxy_read_timeout` - (Optional, Int) 300s.
* `proxy_send_timeout` - (Optional, Int) 300s.
* `sni_host` - (Optional, String) When SniType=3, this parameter needs to be filled in to represent a custom host.
* `sni_type` - (Optional, Int) Sni type fo upstream, 0:disable SNI; 1:enable SNI and SNI equal original request host; 2:and SNI equal upstream host 3:enable SNI and equal customize host.
* `src_list` - (Optional, List: [`String`]) Upstream IP List, When UpstreamType=0, this parameter needs to be filled.
* `ssl_id` - (Optional, String) Certificate ID, When CertType=2, this parameter needs to be filled.
* `tls_version` - (Optional, Int) Version of TLS Protocol.
* `upstream_domain` - (Optional, String) Upstream domain, When UpstreamType=1, this parameter needs to be filled.
* `upstream_scheme` - (Optional, String) Upstream scheme for https, http or https.
* `upstream_type` - (Optional, Int) Upstream type, 0 represents IP, 1 represents domain name.
* `weights` - (Optional, List: [`Int`]) Weight of each upstream.
* `xff_reset` - (Optional, Int) 0:disable xff reset; 1:ensable xff reset.

The `ports` object supports the following:

* `port` - (Required, String) Listening port.
* `protocol` - (Required, String) The listening protocol of listening port.
* `upstream_port` - (Required, String) The upstream port for listening port.
* `upstream_protocol` - (Required, String) The upstream protocol for listening port.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `domain_id` - Domain id.


## Import

waf saas_domain can be imported using the id, e.g.

```
terraform import tencentcloud_waf_saas_domain.example waf_2kxtlbky01b3wceb#tf.example.com#9647c91da0aa5f5aaa49d0ca40e2af24
```

