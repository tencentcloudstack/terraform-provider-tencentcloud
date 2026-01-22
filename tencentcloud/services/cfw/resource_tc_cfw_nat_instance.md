Provides a resource to create a CFW nat instance

Example Usage

If mode is 0

```hcl
resource "tencentcloud_cfw_nat_instance" "example" {
  name         = "tf_example"
  cross_a_zone = 1
  width        = 20
  mode         = 0
  new_mode_items {
    vpc_list = [
      "vpc-40hif9or"
    ]
    eips = [
      "119.29.107.37"
    ]
  }
  zone_set = [
    "ap-guangzhou-6",
    "ap-guangzhou-7"
  ]
}
```

If mode is 1

```hcl
resource "tencentcloud_cfw_nat_instance" "example" {
  name         = "tf_example"
  cross_a_zone = 1
  width        = 20
  mode         = 1
  nat_gw_list = [
    "nat-9wwkz1kr"
  ]

  zone_set = [
    "ap-guangzhou-6",
    "ap-guangzhou-7"
  ]
}
```

Import

CFW nat instance can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_nat_instance.example cfwnat-54a21421
```
