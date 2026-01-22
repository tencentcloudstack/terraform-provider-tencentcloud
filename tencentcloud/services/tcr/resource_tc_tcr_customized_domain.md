Provides a resource to create a tcr customized domain

Example Usage

Create a tcr customized domain

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example"
  instance_type = "premium"
  tags = {
    createdBy = "Terraform"
  }
}

resource "tencentcloud_tcr_customized_domain" "example" {
  registry_id    = tencentcloud_tcr_instance.example.id
  domain_name    = "www.demo.com"
  certificate_id = "your_cert_id"
  tags = {
    createdBy = "Terraform"
  }
}
```

Import

tcr customized domain can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_customized_domain.example tcr-fjvvsfdh#www.demo.com
```