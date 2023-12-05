Provides a resource to create a cfw vpc_instance

Example Usage

If mode is 0

```hcl
resource "tencentcloud_cfw_vpc_instance" "example" {
  name = "tf_example"
  mode = 0

  vpc_fw_instances {
    name    = "fw_ins_example"
    vpc_ids = [
      "vpc-291vnoeu",
      "vpc-39ixq9ci"
    ]
    fw_deploy {
      deploy_region = "ap-guangzhou"
      width         = 1024
      cross_a_zone  = 1
      zone_set      = [
        "ap-guangzhou-6",
        "ap-guangzhou-7"
      ]
    }
  }

  switch_mode = 1
  fw_vpc_cidr = "auto"
}
```

If mode is 1

```hcl
resource "tencentcloud_cfw_vpc_instance" "example" {
  name = "tf_example"
  mode = 1

  vpc_fw_instances {
    name = "fw_ins_example"
    fw_deploy {
      deploy_region = "ap-guangzhou"
      width         = 1024
      cross_a_zone  = 0
      zone_set      = [
        "ap-guangzhou-6"
      ]
    }
  }

  ccn_id      = "ccn-peihfqo7"
  switch_mode = 1
  fw_vpc_cidr = "auto"
}
```

Import

cfw vpc_instance can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_vpc_instance.example cfwg-4ee69507
```