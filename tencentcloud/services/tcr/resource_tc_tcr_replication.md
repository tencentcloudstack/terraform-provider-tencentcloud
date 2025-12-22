Provides a resource to create a TCR replication

Example Usage

```hcl
resource "tencentcloud_tcr_replication" "example" {
  source_registry_id      = "tcr-9q9h1nof"
  destination_registry_id = "tcr-jtih9ngc"
  rule {
    name           = "tf-example"
    dest_namespace = ""
    override       = true
    deletion       = true
    filters {
      type  = "name"
      value = "tf-example/**"
    }
  }

  destination_region_id = 1
  description           = "remark."
}
```
