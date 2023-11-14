package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseRenewInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseRenewInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_renew_instance.renew_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_renew_instance.renew_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseRenewInstance = `

resource "tencentcloud_lighthouse_renew_instance" "renew_instance" {
  instance_ids = 
  instance_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_MANUAL_RENEW"

  }
  renew_data_disk = true
  auto_voucher = false
}

`
