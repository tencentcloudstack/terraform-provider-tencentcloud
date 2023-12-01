data "tencentcloud_ssl_certificates" "svr" {
  name = tencentcloud_ssl_certificate.example.name
}

resource "tencentcloud_ssl_certificate" "example" {
  name = "ssl-svr"
  type = "SVR"
  cert = var.cert
  key  = var.key
}
