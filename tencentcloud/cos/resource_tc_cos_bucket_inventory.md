Provides a resource to create a cos bucket_inventory

Example Usage

```hcl
resource "tencentcloud_cos_bucket_inventory" "bucket_inventory" {
    name = "test123"
    bucket = "keep-test-xxxxxx"
    is_enabled = "true"
    included_object_versions = "Current"
    optional_fields {
        fields = ["Size", "ETag"]
    }
    filter {
        period {
            start_time = "1687276800"
        }
    }
    schedule {
        frequency = "Weekly"
    }
    destination {
        bucket = "qcs::cos:ap-guangzhou::keep-test-xxxxxx"
        account_id = ""
        format = "CSV"
        prefix = "cos_bucket_inventory"

    }
}
```

Import

cos bucket_inventory can be imported using the id, e.g.

```
terraform import tencentcloud_cos_bucket_inventory.bucket_inventory bucket_inventory_id
```