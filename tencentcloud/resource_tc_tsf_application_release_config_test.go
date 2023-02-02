package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTsfApplicationReleaseConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationReleaseConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_release_config.application_release_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_application_release_config.application_release_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApplicationReleaseConfig = `

resource "tencentcloud_tsf_application_release_config" "application_release_config" {
  config_id = ""
  group_id = ""
  release_desc = ""
                    }

`
