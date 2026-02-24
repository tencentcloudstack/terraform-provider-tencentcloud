Provides a resource to create a CVM launch template

Example Usage

```hcl
data "tencentcloud_images" "example" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "CentOS 8.2"
}

resource "tencentcloud_cvm_launch_template" "example" {
  launch_template_name = "tf-example"
  placement {
    zone       = "ap-guangzhou-6"
    project_id = 0
  }

  image_id                            = data.tencentcloud_images.example.images.0.image_id
  launch_template_version_description = "CentOS 8.2"
  instance_type                       = "S5.SMALL1"
  instance_charge_type                = "POSTPAID_BY_HOUR"
  system_disk {
    disk_size = 50
    disk_type = "CLOUD_PREMIUM"
  }

  data_disks {
    disk_size = 200
    disk_type = "CLOUD_PREMIUM"
  }

  virtual_private_cloud {
    subnet_id = "subnet-5l1ya4my"
    vpc_id    = "vpc-0m6078eb"
  }

  internet_accessible {
    internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
    public_ip_assigned   = false
  }

  instance_count = 1
  instance_name  = "instanceName"
  host_name      = "root"
  security_group_ids = [
    "sg-4z20n68d",
  ]

  enhanced_service {
    automation_service {
      enabled = true
    }

    monitor_service {
      enabled = true
    }

    security_service {
      enabled = true
    }
  }
}
```