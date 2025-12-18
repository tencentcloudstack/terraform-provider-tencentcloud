package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhAssetSyncFlagConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhAssetSyncFlagConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_asset_sync_flag_config.example", "id"),
				),
			},
			{
				Config: testAccBhAssetSyncFlagConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_asset_sync_flag_config.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_bh_asset_sync_flag_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBhAssetSyncFlagConfig = `
resource "tencentcloud_bh_asset_sync_flag_config" "example" {
  auto_sync = true
}
`

const testAccBhAssetSyncFlagConfigUpdate = `
resource "tencentcloud_bh_asset_sync_flag_config" "example" {
  auto_sync = false
}
`
