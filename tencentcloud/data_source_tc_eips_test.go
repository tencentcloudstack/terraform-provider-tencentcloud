package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudEipsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.eip"),
					resource.TestCheckResourceAttr("data.tencentcloud_eips.data_eips", "eip_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.data_eips", "eip_list.0.eip_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_eips.data_eips", "eip_list.0.eip_name", "tf-test-eip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.data_eips", "eip_list.0.eip_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.data_eips", "eip_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.data_eips", "eip_list.0.public_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eips.data_eips", "eip_list.0.create_time"),
				),
			},
		},
	})
}

const testAccEipsDataSource = `
resource "tencentcloud_eip" "eip" {
	name = "tf-test-eip"
}

data "tencentcloud_eips" "data_eips" {
	eip_id = "${tencentcloud_eip.eip.id}"
}
`
