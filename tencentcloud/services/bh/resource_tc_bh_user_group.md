Provides a resource to create a BH user group

Example Usage

```hcl
resource "tencentcloud_bh_user_group" "example" {
  name = "tf-example"
}
```

Import

BH user group can be imported using the id, e.g.

```
terraform import tencentcloud_bh_user_group.example 92
```
