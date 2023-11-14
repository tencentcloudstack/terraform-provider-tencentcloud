package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveBackupStreamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveBackupStream,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_backup_stream.backup_stream", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_backup_stream.backup_stream",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveBackupStream = `

resource "tencentcloud_live_backup_stream" "backup_stream" {
  push_domain_name = ""
  app_name = ""
  stream_name = ""
  upstream_sequence = ""
}

`
