Provides a resource to attach/detach the corresponding certificate for the domain name in specified cos bucket.

Example Usage

```hcl

resource "tencentcloud_cos_bucket_domain_certificate_attachment" "foo" {
  bucket = ""
  domain_certificate {
	domain = "domain_name"
    certificate {
      cert_type = "CustomCert"
      custom_cert {
        cert        = "===CERTIFICATE==="
        private_key = "===PRIVATE_KEY==="
      }
    }
  }
}

```