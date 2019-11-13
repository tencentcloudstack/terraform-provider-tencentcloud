package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudSecurityGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_security_group.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_security_group.foo", "name", "tf-ci-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_security_group.foo", "description", "terraform-ci-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_security_group.foo", "be_associate_count", "2"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudSecurityGroupConfig = `
resource "tencentcloud_security_group" "foo" {
  name        = "tf-ci-test"
  description = "terraform-ci-test"
}

resource "tencentcloud_security_group_lite_rule" "foo" {
  security_group_id = "${tencentcloud_security_group.foo.id}"

  ingress = [
    "ACCEPT#192.168.1.0/24#80#TCP",
  ]

  egress = [
    "ACCEPT#192.168.0.0/16#ALL#TCP",
  ]
}

data "tencentcloud_security_group" "foo" {
  security_group_id = "${tencentcloud_security_group.foo.id}"
}
`
