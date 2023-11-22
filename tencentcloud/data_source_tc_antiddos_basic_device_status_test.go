package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosBasicDeviceStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosBasicDeviceStatusDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_antiddos_basic_device_status.basic_device_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_basic_device_status.basic_device_status", "data.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_basic_device_status.basic_device_status", "clb_data.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_antiddos_basic_device_status.basic_device_status", "data.0.key", "127.0.0.1"),
				),
			},
		},
	})
}

const testAccAntiddosBasicDeviceStatusDataSource = `

data "tencentcloud_antiddos_basic_device_status" "basic_device_status" {
  ip_list = [
    "127.0.0.1"
  ]
  filter_region = 1
}

`
