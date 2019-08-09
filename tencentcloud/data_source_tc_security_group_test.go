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
					resource.TestCheckResourceAttr("data.tencentcloud_security_group.foo", "name", "ci-temp-security-groups-test"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudSecurityGroupConfig = `
resource "tencentcloud_security_group" "foo" {
  name        = "ci-temp-security-groups-test"
  description = "ci-temp-security-groups-test"
}

data "tencentcloud_security_group" "foo" {
	security_group_id = "${tencentcloud_security_group.foo.id}"
}
`
