resource "tencentcloud_protocol_template" "example" {
  name      = "example"
  protocols = ["udp:all", "tcp:90,110", "icmp", "tcp:1110-1120"]
}

resource "tencentcloud_protocol_template_group" "example" {
  name         = "example"
  template_ids = [tencentcloud_protocol_template.example.id]
}

data "tencentcloud_protocol_templates" "example" {
  id = tencentcloud_protocol_template.example.id
}

data "tencentcloud_protocol_template_groups" "example" {
}
