Provides a resource to create a lighthouse key_pair

Example Usage

```hcl
resource "tencentcloud_lighthouse_key_pair" "key_pair" {
  key_name = "key_name_test"
}
```

Import

lighthouse key_pair can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_key_pair.key_pair key_pair_id
```