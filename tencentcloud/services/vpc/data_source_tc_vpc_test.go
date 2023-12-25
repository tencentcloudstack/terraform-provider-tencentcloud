package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudVpc_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudVpcConfig_id,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpc.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc.id", "name", "tf-ci-test"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudVpcConfig_id = `
resource "tencentcloud_vpc" "foo" {
  name       = "tf-ci-test"
  cidr_block = "10.0.0.0/16"
}

data "tencentcloud_vpc" "id" {
  id = tencentcloud_vpc.foo.id
}
`
