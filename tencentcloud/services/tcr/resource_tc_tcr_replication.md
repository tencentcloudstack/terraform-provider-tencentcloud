Provides a resource to create a TCR replication

Example Usage

Sync source tcr instance to target instance

Synchronize an existing tcr instance to the destination instance. This operation is often used in the cross-multiple region scenario.
Assume you have had two TCR instances before this operation. This example shows how to sync a tcr instance from ap-guangzhou(gz) to ap-shanghai(sh).

```hcl
# tcr instance on ap-guangzhou
resource "tencentcloud_tcr_instance" "example_gz" {
  name          = "tf-example-tcr-gz"
  instance_type = "premium"
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example_gz" {
  instance_id    = tencentcloud_tcr_instance.example_gz.id
  name           = "tf_example_ns_gz"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

# tcr instance on ap-shanghai
resource "tencentcloud_tcr_instance" "example_sh" {
  name          = "tf-example-tcr-sh"
  instance_type = "premium"
  delete_bucket = true
}

resource "tencentcloud_tcr_namespace" "example_sh" {
  instance_id    = tencentcloud_tcr_instance.example_sh.id
  name           = "tf_example_ns_sh"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_replication" "example" {
  source_registry_id      = tencentcloud_tcr_instance.example_gz.id
  destination_registry_id = tencentcloud_tcr_instance.example_sh.id
  rule {
    name           = "tf-example"
    dest_namespace = tencentcloud_tcr_namespace.example_sh.name
    override       = true
    deletion       = true
    filters {
      type  = "name"
      value = join("/", [tencentcloud_tcr_namespace.example_gz.name, "**"])
    }
  }

  destination_region_id = 1
  description           = "remark."
}
```
