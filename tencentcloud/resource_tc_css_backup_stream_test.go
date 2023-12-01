package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCssBackupStreamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssBackupStream,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_backup_stream.backup_stream", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_backup_stream.backup_stream", "app_name", "live"),
					resource.TestCheckResourceAttr("tencentcloud_css_backup_stream.backup_stream", "push_domain_name", "177154.push.tlivecloud.com"),
					resource.TestCheckResourceAttr("tencentcloud_css_backup_stream.backup_stream", "stream_name", "1308919341_test"),
					resource.TestCheckResourceAttr("tencentcloud_css_backup_stream.backup_stream", "upstream_sequence", "2210137729805899152"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_backup_stream.backup_stream",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssBackupStream = `

resource "tencentcloud_css_backup_stream" "backup_stream" {
  push_domain_name  = "177154.push.tlivecloud.com"
  app_name          = "live"
  stream_name       = "1308919341_test"
  upstream_sequence = "2210137729805899152"
}

`
