package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcStoreLocationConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcStoreLocationConfig,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_store_location_config.store_location_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_store_location_config.store_location_config", "store_location", "cosn://cos-lock-1308919341/test/")),
			},
			{
				Config: testAccDlcStoreLocationConfigUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_store_location_config.store_location_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_store_location_config.store_location_config", "store_location", "cosn://cos-lock-1308919341/")),
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
}

`
const testAccDlcStoreLocationConfigUpdate = `

resource "tencentcloud_dlc_store_location_config" "store_location_config" {
  store_location = "cosn://cos-lock-1308919341/"
}

`
