Provides a resource to manage address template.

~> **NOTE:** It can be replaced by `tencentcloud_address_extra_template`.

Example Usage

```hcl
resource "tencentcloud_address_template" "foo" {
  name      = "cam-user-test"
  addresses = ["10.0.0.1", "10.0.1.0/24", "10.0.0.1-10.0.0.100"]
}
```

Import

Address template can be imported using the id, e.g.

```
$ terraform import tencentcloud_address_template.foo ipm-makf7k9e"
```