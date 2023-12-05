Provides a resource to create a kms white_box_key

Example Usage

```hcl
resource "tencentcloud_kms_white_box_key" "example" {
  alias       = "tf_example"
  description = "test desc."
  algorithm   = "SM4"
  status      = "Enabled"
  tags        = {
    "createdBy" = "terraform"
  }
}
```

Import

kms white_box_key can be imported using the id, e.g.

```
terraform import tencentcloud_kms_white_box_key.example 244dab8c-6dad-11ea-80c6-5254006d0810
```