resource "tencentcloud_security_group" "example" {
  name        = "tf-example-sg"
  description = "sg test"
  project_id  = 0

  tags = {
    "example" = "test"
  }
}