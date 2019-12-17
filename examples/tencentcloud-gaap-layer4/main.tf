resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "119.29.29.29"
  name = "ci-test-gaap-realserver2"
}

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener-new"
  port            = 80
  scheduler       = "wrr"
  realserver_type = "IP"
  proxy_id        = tencentcloud_gaap_proxy.foo.id
  health_check    = true
  interval        = 11
  connect_timeout = 10

  realserver_bind_set {
    id     = tencentcloud_gaap_realserver.foo.id
    ip     = tencentcloud_gaap_realserver.foo.ip
    port   = 80
    weight = 1
  }

  realserver_bind_set {
    id     = tencentcloud_gaap_realserver.bar.id
    ip     = tencentcloud_gaap_realserver.bar.ip
    port   = 80
    weight = 2
  }
}

data tencentcloud_gaap_layer4_listeners "foo" {
  protocol    = "TCP"
  proxy_id    = tencentcloud_gaap_proxy.foo.id
  listener_id = tencentcloud_gaap_layer4_listener.foo.id
}