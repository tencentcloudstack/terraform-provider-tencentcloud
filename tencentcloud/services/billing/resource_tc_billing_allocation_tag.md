Provides a resource to create a Billing allocation tag

Example Usage

```hcl
resource "tencentcloud_tag" "example" {
  tag_key   = "tagKey"
  tag_value = "tagValue"
}

resource "tencentcloud_billing_allocation_tag" "example" {
  tag_key = tencentcloud_tag.example.tag_key
}
```

Import

Billing allocation tag can be imported using the id, e.g.

```
terraform import tencentcloud_billing_allocation_tag.example tagKey
```
