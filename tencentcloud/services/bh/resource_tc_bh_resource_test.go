package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhResourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_resource.example", "id"),
				),
			},
			{
				Config: testAccBhResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_resource.example", "id"),
				),
			},
		},
	})
}

const testAccBhResource = `
resource "tencentcloud_bh_resource" "example" {
  deploy_region    = "ap-guangzhou"
  vpc_id           = "vpc-q1of50wz"
  subnet_id        = "subnet-7uhvm46o"
  resource_edition = "standard"
  resource_node    = 20
  time_unit        = "m"
  time_span        = "1"
  pay_mode         = 1
  auto_renew_flag  = 1
  deploy_zone      = "ap-guangzhou-6"
  cidr_block       = "192.168.11.0/24"
  vpc_cidr_block   = "192.168.0.0/16"
}
`

const testAccBhResourceUpdate = `
resource "tencentcloud_bh_resource" "example" {
  deploy_region    = "ap-guangzhou"
  vpc_id           = "vpc-q1of50wz"
  subnet_id        = "subnet-7uhvm46o"
  resource_edition = "standard"
  resource_node    = 10
  time_unit        = "m"
  time_span        = "1"
  pay_mode         = 1
  auto_renew_flag  = 0
  deploy_zone      = "ap-guangzhou-6"
  cidr_block       = "192.168.11.0/24"
  vpc_cidr_block   = "192.168.0.0/16"
}
`
