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
					resource.TestCheckResourceAttr("data.tencentcloud_security_group.foo", "create_time", "2017-07-31 20:03:00"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudSecurityGroupConfig = `
data "tencentcloud_security_group" "foo" {
	security_group_id = "sg-icy671l9"
}
`
