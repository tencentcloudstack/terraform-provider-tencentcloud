resource "tencentcloud_address_template" "example" {
  name      = "example"
  addresses = ["10.0.0.0/24", "1.1.1.1", "1.0.0.1-1.0.0.100"]
}

resource "tencentcloud_address_template_group" "example" {
  name         = "example"
  template_ids = [tencentcloud_address_template.example.id]
}

data "tencentcloud_address_template_groups" "example" {
}
