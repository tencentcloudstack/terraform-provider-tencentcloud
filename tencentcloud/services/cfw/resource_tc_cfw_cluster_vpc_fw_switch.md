Provides a resource to create a CFW cluster vpc fw switch

Example Usage

If switch_mode is 2

```hcl
resource "tencentcloud_cfw_cluster_vpc_fw_switch" "example" {
  ccn_id       = "ccn-8qv0ro89"
  switch_mode  = 2
  routing_mode = 0
  region_cidr_configs {
    region      = "ap-guangzhou"
    cidr_mode   = 1
    custom_cidr = ""
  }
}
```

If switch_mode is 1

```hcl
resource "tencentcloud_cfw_cluster_vpc_fw_switch" "example" {
  ccn_id       = "ccn-8qv0ro89"
  switch_mode  = 1
  routing_mode = 1
  region_cidr_configs {
    region      = "ap-guangzhou"
    cidr_mode   = 0
    custom_cidr = ""
  }
  
  region_cidr_configs {
    region      = "ap-chongqing"
    cidr_mode   = 0
    custom_cidr = ""
  }

   region_cidr_configs {
    region      = "ap-shanghai"
    cidr_mode   = 1
    custom_cidr = ""
  }

  interconnect_pairs {
    interconnect_mode = "FullMesh"
    group_a {
      instance_id      = "vpc-264i7uzy"
      instance_type    = "VPC"
      instance_region  = "ap-shanghai"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.124.0.0/16"
      ]
    }

    group_a {
      instance_id      = "vpc-h2i9m8xh"
      instance_type    = "VPC"
      instance_region  = "ap-chongqing"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.25.0.0/16"
      ]
    }

    group_b {
      instance_id      = "vpc-264i7uzy"
      instance_type    = "VPC"
      instance_region  = "ap-shanghai"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.124.0.0/16"
      ]
    }

    group_b {
      instance_id      = "vpc-h2i9m8xh"
      instance_type    = "VPC"
      instance_region  = "ap-chongqing"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.25.0.0/16"
      ]
    }
  }

  interconnect_pairs {
    interconnect_mode = "CrossConnect"
    group_a {
      instance_id      = "vpc-5l5uqrgx"
      instance_type    = "VPC"
      instance_region  = "ap-chongqing"
      access_cidr_mode = 1
      access_cidr_list = [
        "192.168.0.0/16"
      ]
    }

    group_b {
      instance_id      = "vpc-1yoh1nhh"
      instance_type    = "VPC"
      instance_region  = "ap-guangzhou"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.208.0.0/24",
        "172.16.0.0/16"
      ]
    }
  }
}
```
