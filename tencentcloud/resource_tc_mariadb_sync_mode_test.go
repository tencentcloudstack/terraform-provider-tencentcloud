package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbSyncModeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbSyncMode,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_sync_mode.sync_mode", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_sync_mode.sync_mode",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbSyncMode = `

resource "tencentcloud_mariadb_sync_mode" "sync_mode" {
  instance_id = ""
  sync_mode = 
}

`
