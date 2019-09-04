resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-new"
  bandwidth         = 20
  concurrent        = 10
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
  enable            = false

  tags = {
    "test" = "test"
  }
}

data tencentcloud_gaap_proxies "foo" {
  ids = ["${tencentcloud_gaap_proxy.foo.id}"]
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"

  tags = {
    "test" = "test"
  }
}

data tencentcloud_gaap_realservers "foo" {
  ip = "${tencentcloud_gaap_realserver.foo.ip}"
}