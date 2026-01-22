Provides a resource to create a CLB listener.

Example Usage

HTTP Listener

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id        = "lb-qck8thny"
  listener_name = "tf-example"
  port          = 80
  protocol      = "HTTP"
}
```

TCP/UDP Listener

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id                     = "lb-qck8thny"
  listener_name              = "tf-example"
  port                       = 80
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 2
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_port          = 200
  health_check_type          = "HTTP"
  health_check_http_path     = "/"
  health_check_http_code     = 2
  health_check_http_version  = "HTTP/1.0"
  health_check_http_method   = "GET"
  deregister_target_rst      = false
  idle_connect_timeout       = 900
}
```

TCP/UDP Listener with tcp health check

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id                     = "lb-qck8thny"
  listener_name              = "tf-example"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "TCP"
  health_check_port          = 200
  deregister_target_rst      = false
  idle_connect_timeout       = 900
}
```

TCP/UDP Listener with http health check

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id                     = "lb-qck8thny"
  listener_name              = "tf-example"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "HTTP"
  health_check_http_domain   = "www.tencent.com"
  health_check_http_code     = 16
  health_check_http_version  = "HTTP/1.1"
  health_check_http_method   = "HEAD"
  health_check_http_path     = "/"
  deregister_target_rst      = false
  idle_connect_timeout       = 900
}
```

TCP/UDP Listener with customer health check

```hcl
resource "tencentcloud_clb_listener" "example"{
  clb_id                     = "lb-qck8thny"
  listener_name              = "tf-example"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "CUSTOM"
  health_check_context_type  = "HEX"
  health_check_send_context  = "0123456789ABCDEF"
  health_check_recv_context  = "ABCD"
  target_type                = "TARGETGROUP"
}
```

HTTPS Listener with sigle certificate

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id               = "lb-0lh5au7v"
  listener_name        = "tf-example"
  port                 = "80"
  protocol             = "HTTPS"
  certificate_ssl_mode = "MUTUAL"
  certificate_id       = "VjANRdz8"
  certificate_ca_id    = "VfqO4zkB"
  sni_switch           = true
}
```

HTTPS Listener with multi certificates

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id               = "lb-l6cp6jt4"
  listener_name        = "tf-example"
  port                 = "80"
  protocol             = "HTTPS"
  sni_switch           = true

  multi_cert_info {
    ssl_mode = "UNIDIRECTIONAL"
    cert_id_list = [
      "LCYouprI",
      "JVO1alRN"
    ]
  }
}
```

TCP SSL Listener

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id                     = "lb-0lh5au7v"
  listener_name              = "tf-example"
  port                       = "80"
  protocol                   = "TCP_SSL"
  certificate_ssl_mode       = "MUTUAL"
  certificate_id             = "VjANRdz8"
  certificate_ca_id          = "VfqO4zkB"
  health_check_switch        = true
  health_check_time_out      = 2
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  scheduler                  = "WRR"
  target_type                = "TARGETGROUP"
}
```

Port Range Listener

```hcl
resource "tencentcloud_clb_instance" "example" {
  clb_name     = "tf-listener-test"
  network_type = "OPEN"
}

resource "tencentcloud_clb_listener" "example" {
  clb_id              = tencentcloud_clb_instance.example.id
  listener_name       = "tf-example"
  port                = 1
  end_port            = 6
  protocol            = "TCP"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "NODE"
}
```

Import

CLB listener can be imported using the clbId#listenerId (version >= 1.47.0), e.g.

```
$ terraform import tencentcloud_clb_listener.example lb-7a0t6zqb#lbl-hh141sn9
```
