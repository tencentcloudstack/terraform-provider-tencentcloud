package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudRumOfflineLogConfigAttachment_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumOfflineLogConfigAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_offline_log_config_attachment.offline_log_config_attachment", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_rum_offline_log_config_attachment.offlineLogConfigAttachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRumOfflineLogConfigAttachment = `

resource "tencentcloud_rum_offline_log_config_attachment" "offline_log_config_attachment" {
  project_key = ""
  unique_i_d = ""
}

`
