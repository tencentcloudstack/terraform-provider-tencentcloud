Provides a resource to create a configuration for an AS (Auto scaling) instance.

~> **NOTE:**  In order to ensure the integrity of customer data, if the cvm instance was destroyed due to shrinking, it will keep the cbs associate with cvm by default. If you want to destroy together, please set `delete_with_instance` to `true`.

Example Usage

```hcl
data "tencentcloud_images" "example" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "example-launch-configuration"
  image_id           = data.tencentcloud_images.example.images.0.image_id
  instance_types     = ["SA1.SMALL1"]
  project_id         = 0
  system_disk_type   = "CLOUD_PREMIUM"
  system_disk_size   = "50"

  data_disk {
    disk_type = "CLOUD_PREMIUM"
    disk_size = 50
  }

  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out = 10
  public_ip_assigned         = true
  password                   = "Test@123#"
  enhanced_security_service  = false
  enhanced_monitor_service   = false
  user_data                  = "dGVzdA=="

  host_name_settings {
	host_name       = "host-name-test"
	host_name_style = "UNIQUE"
  }

  instance_tags = {
    tag = "example"
  }
}
```

Using `SPOTPAID` charge type

```
data "tencentcloud_images" "example" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name   = "launch-configuration"
  image_id             = data.tencentcloud_images.example.images.0.image_id
  instance_types       = ["SA1.SMALL1"]
  instance_charge_type = "SPOTPAID"
  spot_instance_type   = "one-time"
  spot_max_price       = "1000"
}
```

Import

AutoScaling Configuration can be imported using the id, e.g.

```
$ terraform import tencentcloud_as_scaling_config.example asc-n32ymck2
```