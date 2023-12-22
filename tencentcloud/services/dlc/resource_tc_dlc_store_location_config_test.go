package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcStoreLocationConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcStoreLocationConfig,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_store_location_config.store_location_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_store_location_config.store_location_config", "store_location", "cosn://cos-lock-1308919341/test/"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_store_location_config.store_location_config", "enable", "1"),
				),
			},
			{
				Config: testAccDlcStoreLocationConfigUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_store_location_config.store_location_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_store_location_config.store_location_config", "store_location", "cosn://cos-lock-1308919341/"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_store_location_config.store_location_config", "enable", "1")),
			},
			{
				ResourceName:      "tencentcloud_dlc_store_location_config.store_location_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcStoreLocationConfig = `

resource "tencentcloud_dlc_store_location_config" "store_location_config" {
  store_location = "cosn://cos-lock-1308919341/test/"
  enable = 1
}

`
const testAccDlcStoreLocationConfigUpdate = `

resource "tencentcloud_dlc_store_location_config" "store_location_config" {
  store_location = "cosn://cos-lock-1308919341/"
  enable = 1
}

`
