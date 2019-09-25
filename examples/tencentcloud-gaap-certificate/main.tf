resource tencentcloud_gaap_certificate "server" {
  type    = "SERVER"
  name    = "test"
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

data "tencentcloud_gaap_certificates" "foo" {
  id = "${tencentcloud_gaap_certificate.server.id}"
}

data "tencentcloud_gaap_certificates" "bar" {
  name = "${tencentcloud_gaap_certificate.server.name}"
}