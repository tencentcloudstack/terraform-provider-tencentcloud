Provides a resource to manage address extra template.

~> **NOTE:** Compare to `tencentcloud_address_template`, It contains remarks.


Example Usage

```hcl
resource "tencentcloud_address_extra_template" "foo" {
  name = "demo"

  addresses_extra {
    address     = "10.0.0.1"
    description = "create by terraform"
  }

  addresses_extra {
    address     = "10.0.1.0/24"
    description = "delete by terraform"
  }

  addresses_extra {
    address     = "10.0.0.1-10.0.0.100"
    description = "modify by terraform"
  }

  tags = {
    createBy = "terraform"
    deleteBy = "terraform"
  }

}
```

Import

Address template can be imported using the id, e.g.

```
$ terraform import tencentcloud_address_extra_template.foo ipm-makf7k9e
```