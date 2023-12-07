Provides a resource to operate a tcr image signature.

Example Usage

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "premium"
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf_example_ns"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_repository" "example" {
  instance_id	 = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  name 	         = "test"
  brief_desc 	 = "111"
  description	 = "111111111111111111111111111111111111"
}

resource "tencentcloud_tcr_create_image_signature_operation" "example" {
  registry_id     = tencentcloud_tcr_instance.example.id
  namespace_name  = tencentcloud_tcr_namespace.example.name
  repository_name = tencentcloud_tcr_repository.example.name
  image_version   = "v1"
}
```

Import

tcr image_signature_operation can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_create_image_signature_operation.image_signature_operation image_signature_operation_id
```