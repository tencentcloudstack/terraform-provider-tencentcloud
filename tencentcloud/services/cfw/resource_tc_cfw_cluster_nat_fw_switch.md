Provides a resource to create a CFW cluster NAT firewall switch

Example Usage

```hcl
resource "tencentcloud_cfw_cluster_nat_fw_switch" "example" {
  nat_ccn_switch {
    nat_ins_id   = "nat-h1i1mf4n"
    ccn_id       = "ccn-p3mlp0tj"
    switch_mode  = 1
    routing_mode = 1
    access_instance_list {
      instance_id      = "vpc-8za5vlpc"
      instance_type    = "VPC"
      instance_region  = "ap-shanghai"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.51.0.0/16"
      ]
    }

    access_instance_list {
      instance_id      = "vpc-ec6krnpp"
      instance_type    = "VPC"
      instance_region  = "ap-guangzhou"
      access_cidr_mode = 1
      access_cidr_list = [
        "172.17.0.0/16"
      ]
    }
  }
}
```

Import

CFW cluster NAT firewall switch can be imported using the composite id `nat_ins_id#ccn_id`, e.g.

```
terraform import tencentcloud_cfw_cluster_nat_fw_switch.example nat-h1i1mf4n#ccn-p3mlp0tj
```
