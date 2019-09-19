resource "tencentcloud_security_group" "group" {
  name        = "${var.short_name}"
  description = "New security group"

  tags = {
    "test" = "test"
  }
}
