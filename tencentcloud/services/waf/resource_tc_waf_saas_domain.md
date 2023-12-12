Provides a resource to create a waf saas_domain

Example Usage

If upstream_type is 0

Create a basic waf saas domain

```hcl
resource "tencentcloud_waf_saas_domain" "example" {
  instance_id = "waf_2kxtlbky01b3wceb"
  domain      = "tf.example.com"
  src_list    = [
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

Create a load balancing strategy is weighted polling saas domain

```hcl
resource "tencentcloud_waf_saas_domain" "example" {
  instance_id = "waf_2kxtlbky01b3wceb"
  domain      = "tf.example.com"
  src_list    = [
    "1.1.1.1",
    "2.2.2.2"
  ]
  load_balance = "2"
  weights      = [
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

If upstream_type is 1

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

Create a waf saas domain with set Http&Https

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
  src_list        = [
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

Create a complete waf saas domain

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
  src_list        = [
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

Import

waf saas_domain can be imported using the id, e.g.

```
terraform import tencentcloud_waf_saas_domain.example waf_2kxtlbky01b3wceb#tf.example.com#9647c91da0aa5f5aaa49d0ca40e2af24
```