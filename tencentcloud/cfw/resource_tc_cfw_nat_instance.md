Provides a resource to create a cfw nat_instance

Example Usage

If mode is 0

```hcl
resource "tencentcloud_cfw_nat_instance" "example" {
  name  = "tf_example"
  width = 20
  mode  = 0
  new_mode_items {
    vpc_list = [
      "vpc-5063ta4i"
    ]
    eips = [
      "152.136.168.192"
    ]
  }
  cross_a_zone = 0
  zone_set     = [
    "ap-guangzhou-7"
  ]
}
```

If mode is 1

```hcl
resource "tencentcloud_cfw_nat_instance" "example" {
  name        = "tf_example"
  width       = 20
  mode        = 1
  nat_gw_list = [
    "nat-9wwkz1kr"
  ]
  cross_a_zone = 1
  cross_a_zone = 0
  zone_set     = [
    "ap-guangzhou-6",
    "ap-guangzhou-7"
  ]
}
```

Import

cfw nat_instance can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_nat_instance.example cfwnat-54a21421
```