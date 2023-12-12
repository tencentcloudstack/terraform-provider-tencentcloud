Use this resource to create tcr repository.

Example Usage

Create a tcr repository instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "premium"
  delete_bucket = true
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
```

Import

tcr repository can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_repository.foo instance_id#namespace_name#repository_name
```