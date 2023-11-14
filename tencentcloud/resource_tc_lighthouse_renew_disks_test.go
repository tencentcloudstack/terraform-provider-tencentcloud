package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseRenewDisksResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseRenewDisks,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_renew_disks.renew_disks", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_renew_disks.renew_disks",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseRenewDisks = `

resource "tencentcloud_lighthouse_renew_disks" "renew_disks" {
  disk_ids = 
  renew_disk_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_AUTO_RENEW"
		time_unit = "m"
		cur_instance_deadline = "2018-01-01 00:00:00"

  }
  auto_voucher = true
}

`
