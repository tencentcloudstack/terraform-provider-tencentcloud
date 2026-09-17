Provides a resource to create a BDRC instance copy pair

Example Usage

```hcl
resource "tencentcloud_bdrc_instance_copy_pair" "example" {
  protect_group_id          = "ProtectGroup-xxxxx"
  instance_copy_pair_name   = "tf-example-copy-pair"
  recovery_point_objective  = 15

  create_target_instance_parameters {
    source_instance_id   = "ins-xxxxxxxx"
    instance_charge_type = "POSTPAID_BY_HOUR"
    instance_type        = "S5.MEDIUM4"
    image_id             = "img-xxxxxxxx"
    instance_name        = "tf-example-instance"

    placement {
      zone       = "ap-guangzhou-3"
      project_id = 0
    }

    system_disk {
      disk_type            = "CLOUD_BSSD"
      disk_size            = 50
      delete_with_instance = true
    }

    virtual_private_cloud {
      vpc_id       = "vpc-xxxxxxxx"
      subnet_id    = "subnet-xxxxxxxx"
      as_vpc_gateway = false
    }

    internet_accessible {
      internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
      internet_max_bandwidth_out = 10
      public_ip_assigned         = false
    }

    login_settings {
      password = "TFexample123"
    }

    enhanced_service {
      security_service {
        enabled = true
      }
      monitor_service {
        enabled = true
      }
    }
  }
}
```

Import

BDRC instance copy pair can be imported using the copyPairId, e.g.

```
terraform import tencentcloud_bdrc_instance_copy_pair.example cvmcopypair-xxxxxxxx
```
