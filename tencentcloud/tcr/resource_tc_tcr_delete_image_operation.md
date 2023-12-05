Provides a resource to delete the specified tcr image.

Example Usage

To delete the specified image

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

resource "tencentcloud_tcr_delete_image_operation" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  repository_name = "repo"
  image_version = "v1" # the image want to delete
  namespace_name = tencentcloud_tcr_namespace.example.name
}
```