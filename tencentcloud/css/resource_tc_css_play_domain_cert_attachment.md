Provides a resource to create a css play_domain_cert_attachment. This resource is used for binding the play domain and specified certification together.

Example Usage

```hcl
data "tencentcloud_ssl_certificates" "foo" {
	name = "your_ssl_cert"
}

resource "tencentcloud_css_play_domain_cert_attachment" "play_domain_cert_attachment" {
  cloud_cert_id = data.tencentcloud_ssl_certificates.foo.certificates.0.id
  domain_info {
    domain_name = "your_domain_name"
    status = 1
  }
}
```

Import

css play_domain_cert_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_css_play_domain_cert_attachment.play_domain_cert_attachment domainName#cloudCertId
```