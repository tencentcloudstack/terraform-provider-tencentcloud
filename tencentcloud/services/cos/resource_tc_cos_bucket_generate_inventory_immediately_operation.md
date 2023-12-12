Provides a resource to generate a cos bucket inventory immediately

Example Usage

```hcl
resource "tencentcloud_cos_bucket_generate_inventory_immediately_operation" "generate_inventory_immediately" {
    inventory_id = "test"
    bucket = "keep-test-xxxxxx"
}
```