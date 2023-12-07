Provides a resource to create a cfw address_template

Example Usage

If type is 1

```hcl
resource "tencentcloud_cfw_address_template" "example" {
  name      = "tf_example"
  detail    = "test template"
  ip_string = "1.1.1.1,2.2.2.2"
  type      = 1
}
```

If type is 5

```hcl
resource "tencentcloud_cfw_address_template" "example" {
  name      = "tf_example"
  detail    = "test template"
  ip_string = "www.qq.com,www.tencent.com"
  type      = 5
}
```
Import

cfw address_template can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_address_template.example mb_1300846651_1695611353900
```