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
					resource.TestCheckResourceAttr("data.tencentcloud_security_group.foo", "name", "test-foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_security_group.foo", "description", "test-foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_security_group.foo", "be_associate_count", "0"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudSecurityGroupConfig = `
resource "tencentcloud_security_group" "foo" {
    name        = "test-foo"
    description = "test-foo"
}

data "tencentcloud_security_group" "foo" {
	security_group_id = "${tencentcloud_security_group.foo.id}"
}
`
