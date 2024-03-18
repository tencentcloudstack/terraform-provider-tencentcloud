package bh_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbResourceResource_basic -v
func TestAccTencentCloudNeedFixDasbResourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_resource.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "deploy_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "vpc_id", "vpc-fmz6l9nz"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "subnet_id", "subnet-g7jhwhi2"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "resource_edition", "standard"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_resource.example", "resource_node"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "time_unit", "m"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_resource.example", "time_span"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_resource.example", "auto_renew_flag"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "deploy_zone", "ap-guangzhou-6"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_resource.example", "package_bandwidth"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "vpc_cidr_block", "10.35.0.0/16"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "cidr_block", "10.35.20.0/24"),
				),
			},
			{
				ResourceName:      "tencentcloud_dasb_resource.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDasbResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_resource.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "deploy_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "vpc_id", "vpc-fmz6l9nz"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "subnet_id", "subnet-g7jhwhi2"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "resource_edition", "pro"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_resource.example", "resource_node"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "time_unit", "m"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_resource.example", "time_span"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_resource.example", "auto_renew_flag"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "deploy_zone", "ap-guangzhou-6"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_resource.example", "package_bandwidth"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "vpc_cidr_block", "10.35.0.0/16"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_resource.example", "cidr_block", "10.35.20.0/24"),
				),
			},
		},
	})
}

const testAccDasbResource = `
resource "tencentcloud_dasb_resource" "example" {
  deploy_region     = "ap-guangzhou"
  deploy_zone       = "ap-guangzhou-6"
  vpc_id            = "vpc-fmz6l9nz"
  subnet_id         = "subnet-g7jhwhi2"
  vpc_cidr_block    = "10.35.0.0/16"
  cidr_block        = "10.35.20.0/24"
  resource_edition  = "standard"
  resource_node     = 50
  time_unit         = "m"
  time_span         = 1
  auto_renew_flag   = 1
  package_bandwidth = 1
}
`

const testAccDasbResourceUpdate = `
resource "tencentcloud_dasb_resource" "example" {
  deploy_region     = "ap-guangzhou"
  deploy_zone       = "ap-guangzhou-6"
  vpc_id            = "vpc-fmz6l9nz"
  subnet_id         = "subnet-g7jhwhi2"
  vpc_cidr_block    = "10.35.0.0/16"
  cidr_block        = "10.35.20.0/24"
  resource_edition  = "pro"
  resource_node     = 100
  time_unit         = "m"
  time_span         = 1
  auto_renew_flag   = 1
  package_bandwidth = 2
}
`
