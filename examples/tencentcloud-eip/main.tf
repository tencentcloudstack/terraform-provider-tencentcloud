resource "tencentcloud_eip" "foo" {
  name = "awesome_eip_example"

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_eips" "tags" {
  tags = tencentcloud_eip.foo.tags
}