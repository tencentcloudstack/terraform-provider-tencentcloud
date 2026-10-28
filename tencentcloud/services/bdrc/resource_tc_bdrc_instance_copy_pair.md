Provides a resource to create a BDRC instance copy pair

Example Usage

```hcl
resource "tencentcloud_bdrc_disaster_recovery_site_pair" "example" {
  disaster_recovery_type = "CROSS_ZONE"
  source_region          = "ap-shanghai"
  source_zone            = "ap-shanghai-3"
  target_region          = "ap-shanghai"
  target_zone            = "ap-shanghai-4"
  source_vpc             = "vpc-lx6q09ji"
  target_vpc             = "vpc-jktad5e6"
  site_pair_product_type = "INSTANCE"
  site_pair_name         = "tf-example"
  copy_type              = "ASY"
}

resource "tencentcloud_bdrc_disaster_recovery_protect_group" "example" {
  site_pair_id             = tencentcloud_bdrc_disaster_recovery_site_pair.example.site_pair_id
  protect_group_name       = "tf-example"
  data_direction           = "POSITIVE"
  protect_group_type       = "INSTANCE"
  recovery_point_objective = 15
}

resource "tencentcloud_bdrc_disaster_recovery_vpc_mapping" "example" {
  site_pair_id     = tencentcloud_bdrc_disaster_recovery_site_pair.example.site_pair_id
  source_vpc_id    = "vpc-lx6q09ji"
  source_subnet_id = "subnet-nyrg9pkl"
  target_vpc_id    = "vpc-jktad5e6"
  target_subnet_id = "subnet-jdzuvvvb"
}

resource "tencentcloud_bdrc_security_group_mapping" "example" {
  site_pair_id             = tencentcloud_bdrc_disaster_recovery_site_pair.example.site_pair_id
  src_security_group_id    = "sg-ool7tmf8"
  target_security_group_id = "sg-jfy3gi92"
}

resource "tencentcloud_bdrc_instance_copy_pair" "example" {
  protect_group_id         = tencentcloud_bdrc_disaster_recovery_protect_group.example.protect_group_id
  instance_copy_pair_name  = "tf-example"
  recovery_point_objective = 15

  create_target_instance_parameters {
    source_instance_id   = "ins-0ryh0vvn"
    instance_charge_type = "POSTPAID_BY_HOUR"
    instance_type        = "S5.MEDIUM2"
    image_id             = "img-9qrfy1xt"
    instance_name        = "tf-example-instance"

    placement {
      zone = "ap-shanghai-4"
    }

    system_disk {
      disk_type = "CLOUD_BSSD"
      disk_size = 50
    }

    virtual_private_cloud {
      vpc_id    = "vpc-jktad5e6"
      subnet_id = "subnet-jdzuvvvb"
    }
  }
}
```
