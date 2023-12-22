package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseRenewInstanceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseRenewInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_renew_instance.renew_instance", "id")),
			},
		},
	})
}

const testAccLighthouseRenewInstance = testAccLighthouseInstance + `

resource "tencentcloud_lighthouse_renew_instance" "renew_instance" {
  instance_id = tencentcloud_lighthouse_instance.instance.id
  instance_charge_prepaid {
	period = 1
	renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  }
  renew_data_disk = true
  auto_voucher = false
}

`
