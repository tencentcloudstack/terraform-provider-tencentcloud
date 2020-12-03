resource "tencentcloud_service_template" "example" {
  name     = "example"
  services = ["udp:all", "tcp:90,110", "icmp", "tcp:1110-1120"]
}

resource "tencentcloud_service_template_group" "example" {
  name         = "example"
  template_ids = [tencentcloud_service_template.example.id]
}

data "tencentcloud_service_template_groups" "example" {
}
