package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRumOfflineLogConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumOfflineLogConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_rum_offline_log_config.offline_log_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_rum_offline_log_config.offline_log_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRumOfflineLogConfig = `

resource "tencentcloud_rum_offline_log_config" "offline_log_config" {
  project_key = &lt;nil&gt;
  unique_i_d = &lt;nil&gt;
  }

`
