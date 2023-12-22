package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudSecurityGroup_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_security_group.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_security_group.foo", "name", "tf-ci-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_security_group.foo", "description", "terraform-ci-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_security_group.foo", "be_associate_count", "0"),
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
  security_group_id = tencentcloud_security_group.foo.id

  ingress = [
    "ACCEPT#192.168.1.0/24#80#TCP",
  ]

  egress = [
    "ACCEPT#192.168.0.0/16#ALL#TCP",
  ]
}

data "tencentcloud_security_group" "foo" {
  security_group_id = tencentcloud_security_group.foo.id
}
`
