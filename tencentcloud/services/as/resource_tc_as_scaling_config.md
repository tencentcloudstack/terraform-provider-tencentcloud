Provides a resource to create a configuration for an AS (Auto scaling) instance.

~> **NOTE:**  In order to ensure the integrity of customer data, if the cvm instance was destroyed due to shrinking, it will keep the cbs associate with cvm by default. If you want to destroy together, please set `delete_with_instance` to `true`.

Example Usage

Create a normal configuration

```hcl
data "tencentcloud_images" "example" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 4 for x86_64"
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = data.tencentcloud_images.example.images.0.image_id
  instance_types     = ["SA5.MEDIUM4"]
  project_id         = 0
  system_disk_type   = "CLOUD_PREMIUM"
  system_disk_size   = "50"

  data_disk {
    disk_type = "CLOUD_PREMIUM"
    disk_size = 50
  }

  internet_charge_type              = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out        = 10
  public_ip_assigned                = true
  password                          = "Test@123#"
  enhanced_security_service         = false
  enhanced_monitor_service          = false
  enhanced_automation_tools_service = false
  user_data                         = "dGVzdA=="

  host_name_settings {
    host_name       = "host-name"
    host_name_style = "UNIQUE"
  }

  instance_tags = {
    tag = "example"
  }

  tags = {
    "createdBy" = "Terraform"
    "owner"     = "tf"
  }
}
```

Using SPOTPAID charge type

```hcl
data "tencentcloud_images" "example" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 4 for x86_64"
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name   = "tf-example"
  image_id             = data.tencentcloud_images.example.images.0.image_id
  instance_types       = ["SA5.MEDIUM4"]
  instance_charge_type = "SPOTPAID"
  spot_instance_type   = "one-time"
  spot_max_price       = "1000"

  tags = {
    "createdBy" = "Terraform"
    "owner"     = "tf"
  }
}
```

Using image family

```hcl
resource "tencentcloud_as_scaling_config" "example" {
  image_family                      = "business-daily-update"
  configuration_name                = "as-test-config"
  disk_type_policy                  = "ORIGINAL"
  enhanced_monitor_service          = false
  enhanced_security_service         = false
  enhanced_automation_tools_service = false
  instance_tags             = {}
  instance_types            = [
    "S5.SMALL2",
  ]
  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out = 0
  key_ids                    = []
  project_id                 = 0
  public_ip_assigned         = false
  security_group_ids         = [
    "sg-5275dorp",
  ]
  system_disk_size = 50
  system_disk_type = "CLOUD_BSSD"
}
```

Using DisasterRecoverGroupIds

```hcl
resource "tencentcloud_as_scaling_config" "example" {
  image_family                      = "business-daily-update"
  configuration_name                = "as-test-config"
  disk_type_policy                  = "ORIGINAL"
  enhanced_monitor_service          = false
  enhanced_security_service         = false
  enhanced_automation_tools_service = false
  disaster_recover_group_ids        = ["ps-e2u4ew"]
  instance_tags                     = {}
  instance_types                    = [
    "S5.SMALL2",
  ]
  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out = 0
  key_ids                    = []
  project_id                 = 0
  public_ip_assigned         = false
  security_group_ids         = [
    "sg-5275dorp",
  ]
  system_disk_size = 50
  system_disk_type = "CLOUD_BSSD"
}

```

Create a CDC configuration

```hcl
data "tencentcloud_images" "example" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 4 for x86_64"
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name   = "tf-example"
  image_id             = data.tencentcloud_images.example.images.0.image_id
  instance_types       = ["SA5.MEDIUM4"]
  project_id           = 0
  system_disk_type     = "CLOUD_PREMIUM"
  system_disk_size     = "50"
  instance_charge_type = "CDCPAID"
  dedicated_cluster_id = "cluster-262n63e8"

  data_disk {
    disk_type = "CLOUD_PREMIUM"
    disk_size = 50
  }

  internet_charge_type              = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out        = 10
  public_ip_assigned                = true
  password                          = "Test@123#"
  enhanced_security_service         = false
  enhanced_monitor_service          = false
  enhanced_automation_tools_service = false
  user_data                         = "dGVzdA=="

  host_name_settings {
    host_name       = "host-name"
    host_name_style = "UNIQUE"
  }

  instance_tags = {
    tag = "example"
  }

  tags = {
    "createdBy" = "Terraform"
    "owner"     = "tf"
  }
}
```

Create configuration with AntiDDos Eip

```hcl
data "tencentcloud_images" "example" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 4 for x86_64"
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = data.tencentcloud_images.example.images.0.image_id
  instance_types     = ["SA5.MEDIUM4"]
  project_id         = 0
  system_disk_type   = "CLOUD_HSSD"
  system_disk_size   = "50"
  security_group_ids = ["sg-l222vn6w"]

  data_disk {
    disk_type = "CLOUD_HSSD"
    disk_size = 50
  }

  internet_charge_type              = "BANDWIDTH_PACKAGE"
  internet_max_bandwidth_out        = 100
  public_ip_assigned                = true
  bandwidth_package_id              = "bwp-rp2nx3ab"
  ipv4_address_type                 = "AntiDDoSEIP"
  anti_ddos_package_id              = "bgp-31400fvq"
  is_keep_eip                       = true
  password                          = "Test@123#"
  enhanced_security_service         = false
  enhanced_monitor_service          = false
  enhanced_automation_tools_service = false
  user_data                         = "dGVzdA=="

  host_name_settings {
    host_name       = "host-name"
    host_name_style = "UNIQUE"
  }

  instance_tags = {
    tag = "example"
  }

  tags = {
    "createdBy" = "Terraform"
    "owner"     = "tf"
  }
}
```

Import

AutoScaling Configuration can be imported using the id, e.g.

```
$ terraform import tencentcloud_as_scaling_config.example asc-n32ymck2
```