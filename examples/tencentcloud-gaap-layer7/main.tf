resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = "${var.server_cert}"
  key     = "${var.server_key}"
}

resource tencentcloud_gaap_certificate "bar" {
  type    = "CLIENT"
  content = "${var.client_ca}"
  key     = "${var.client_ca_key}"
}

resource tencentcloud_gaap_certificate "server" {
  type    = "SERVER"
  content = "${var.server_cert}"
  key     = "${var.server_key}"
}

resource tencentcloud_gaap_certificate "client" {
  type    = "CLIENT"
  content = "${var.client_ca}"
  key     = "${var.client_ca_key}"
}

resource tencentcloud_gaap_certificate "realserver" {
  type    = "REALSERVER"
  content = "${var.client_ca}"
  key     = "${var.client_ca_key}"
}

resource tencentcloud_gaap_certificate "basic" {
  type    = "BASIC"
  content = "test:tx2KGdo3zJg/."
}

resource tencentcloud_gaap_certificate "gaap" {
  type    = "PROXY"
  content = "${var.server_cert}"
  key     = "${var.server_key}"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 80
  proxy_id              = "${tencentcloud_gaap_proxy.foo.id}"
  certificate_id        = "${tencentcloud_gaap_certificate.foo.id}"
  client_certificate_id = "${tencentcloud_gaap_certificate.bar.id}"
  forward_protocol      = "HTTPS"
  auth_type             = 1
}

resource tencentcloud_gaap_realserver "foo" {
  domain = "www.qq.com"
  name   = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  domain = "qq.com"
  name   = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id           = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                = "www.qq.com"
  certificate_id        = "${tencentcloud_gaap_certificate.server.id}"
  client_certificate_id = "${tencentcloud_gaap_certificate.client.id}"

  realserver_auth               = true
  realserver_certificate_id     = "${tencentcloud_gaap_certificate.realserver.id}"
  realserver_certificate_domain = "qq.com"

  basic_auth    = true
  basic_auth_id = "${tencentcloud_gaap_certificate.basic.id}"

  gaap_auth    = true
  gaap_auth_id = "${tencentcloud_gaap_certificate.gaap.id}"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id     = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain          = "${tencentcloud_gaap_http_domain.foo.domain}"
  path            = "/"
  realserver_type = "DOMAIN"
  health_check    = false
  forward_host    = "www.qqq.com"

  realservers {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.domain}"
    port = 80
  }

  realservers {
    id   = "${tencentcloud_gaap_realserver.bar.id}"
    ip   = "${tencentcloud_gaap_realserver.bar.domain}"
    port = 80
  }
}

data "tencentcloud_gaap_http_domains" "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "${tencentcloud_gaap_http_domain.foo.domain}"
}

data tencentcloud_gaap_http_rules "foo" {
  listener_id  = "${tencentcloud_gaap_layer7_listener.foo.id}"
  path         = "${tencentcloud_gaap_http_rule.foo.path}"
  forward_host = "${tencentcloud_gaap_http_rule.foo.forward_host}"
}