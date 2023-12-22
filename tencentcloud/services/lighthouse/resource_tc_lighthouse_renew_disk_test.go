package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseRenewDiskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseRenewDisk,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_renew_disk.renew_disk", "id")),
			},
		},
	})
}

const testAccLighthouseRenewDisk = `
resource "tencentcloud_lighthouse_disk" "disk" {
	zone = "ap-guangzhou-3"
	disk_size = 20
	disk_type = "CLOUD_SSD"
	disk_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_AUTO_RENEW"
		time_unit = "m"
  
	}
	disk_name = "tmp"
}

resource "tencentcloud_lighthouse_renew_disk" "renew_disk" {
	disk_id = tencentcloud_lighthouse_disk.disk.id
	renew_disk_charge_prepaid {
		period = 2
		renew_flag = "NOTIFY_AND_MANUAL_RENEW"
		time_unit = "m"
	}
	auto_voucher = true
}

`
