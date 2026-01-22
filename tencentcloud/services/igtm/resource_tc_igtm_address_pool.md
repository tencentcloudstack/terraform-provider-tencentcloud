Provides a resource to create a IGTM address pool

~> **NOTE:** Resource `tencentcloud_igtm_instance` needs to be created before using this resource.

Example Usage

```hcl
resource "tencentcloud_igtm_address_pool" "example" {
  pool_name        = "tf-example"
  traffic_strategy = "WEIGHT"
  address_set {
    addr      = "1.1.1.1"
    is_enable = "ENABLED"
    weight    = 90
  }

  address_set {
    addr      = "2.2.2.2"
    is_enable = "DISABLED"
    weight    = 50
  }
}
```

Import

IGTM address pool can be imported using the id, e.g.

```
terraform import tencentcloud_igtm_address_pool.example 1012132
```
