package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseDiskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseDisk,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_disk.disk", "id"),
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_disk.disk", "disk_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_disk.disk", "disk_size", "20"),
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_disk.disk", "disk_type", "CLOUD_SSD"),
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_disk.disk", "zone", "ap-guangzhou-3"),
				),
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
