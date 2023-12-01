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

  instance_tags = {
    tag = "example"
  }
}