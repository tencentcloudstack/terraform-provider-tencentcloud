package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseDiskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseDisk,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_disk.disk", "id")),
			},
		},
	})
}

const testAccLighthouseDisk = `

resource "tencentcloud_lighthouse_disk" "disk" {
  zone = "ap-guangzhou-3"
  disk_size = 20
  disk_type = "CLOUD_SSD"
  disk_charge_prepaid {
	period = 1
	renew_flag = "NOTIFY_AND_AUTO_RENEW"
	time_unit = "m"

  }
  disk_name = "test"
}

`
