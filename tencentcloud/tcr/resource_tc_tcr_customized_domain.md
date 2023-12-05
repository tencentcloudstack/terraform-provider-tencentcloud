Provides a resource to create a tcr customized domain

Example Usage

Create a tcr customized domain

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "premium"
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_customized_domain" "example" {
  registry_id    = tencentcloud_tcr_instance.example.id
  domain_name    = "www.test.com"
  certificate_id = "your_cert_id"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tcr customized_domain can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_customized_domain.customized_domain customized_domain_id
```