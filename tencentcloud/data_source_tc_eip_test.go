package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudEipDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudEipDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_eip.my_eip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eip.my_eip", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eip.my_eip", "public_ip"),
				),
			},
			{
				Config: testAccTencentCloudEipDataSourceConfig_filter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_eip.my_eip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eip.my_eip", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eip.my_eip", "public_ip"),
					resource.TestCheckResourceAttr("data.tencentcloud_eip.my_eip", "status", "UNBIND"),
				),
			},
		},
	})
}

const testAccTencentCloudEipDataSourceConfig_basic = `
resource "tencentcloud_eip" "foo" {
}

data "tencentcloud_eip" "my_eip" {
}
`

const testAccTencentCloudEipDataSourceConfig_filter = `
resource "tencentcloud_eip" "foo" {
}

data "tencentcloud_eip" "my_eip" {
  filter {
    name   = "address-status"
    values = ["UNBIND"]
  }
}
`
