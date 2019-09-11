resource "tencentcloud_ssl_certificate" "ca" {
  name = "ssl-ca"
  type = "CA"
  cert = "${var.ca}"
}

resource "tencentcloud_ssl_certificate" "svr" {
  name = "ssl-svr"
  type = "SVR"
  cert = "${var.cert}"
  key  = "${var.key}"
}

data "tencentcloud_ssl_certificates" "ca" {
  name = "${tencentcloud_ssl_certificate.ca.name}"
}

data "tencentcloud_ssl_certificates" "svr" {
  type = "${tencentcloud_ssl_certificate.svr.type}"
}