Use this resource to create TCR repository.

Example Usage

Create a tcr repository instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example"
  instance_type = "standard"
  delete_bucket = true
  tags = {
    "createdBy" = "Terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  name        = "tf_example"
  severity    = "medium"
}

resource "tencentcloud_tcr_repository" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  name           = "tf-example"
  brief_desc     = "desc."
  description    = "description."
  force_delete   = true
}
```

Import

TCR repository can be imported using the instanceId#nameSpaceName#name, e.g.

```
terraform import tencentcloud_tcr_repository.example tcr-s1jud21h#tf_example#tf-example
```