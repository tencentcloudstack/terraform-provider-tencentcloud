package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApplicationfileConfigReleaseResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationfileConfigRelease,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_applicationfile_config_release.applicationfile_config_release", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_applicationfile_config_release.applicationfile_config_release",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApplicationfileConfigRelease = `

resource "tencentcloud_tsf_applicationfile_config_release" "applicationfile_config_release" {
  config_id = "dcfg-f-123456"
  group_id = "group-123456"
  release_desc = "product release"
}

`
