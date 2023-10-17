package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamSetPolicyVersionConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamSetPolicyVersionConfig,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_set_policy_version_config.set_policy_version_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_set_policy_version_config.set_policy_version_config", "policy_id", "171162811"),
					resource.TestCheckResourceAttr("tencentcloud_cam_set_policy_version_config.set_policy_version_config", "version_id", "2"),
				),
			},
			{
				Config: testAccCamSetPolicyVersionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_set_policy_version_config.set_policy_version_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_set_policy_version_config.set_policy_version_config", "policy_id", "171162811"),
					resource.TestCheckResourceAttr("tencentcloud_cam_set_policy_version_config.set_policy_version_config", "version_id", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_set_policy_version_config.set_policy_version_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamSetPolicyVersionConfig = `

resource "tencentcloud_cam_set_policy_version_config" "set_policy_version_config" {
  policy_id = 171162811
  version_id = 2
}

`
const testAccCamSetPolicyVersionConfigUpdate = `

resource "tencentcloud_cam_set_policy_version_config" "set_policy_version_config" {
  policy_id = 171162811
  version_id = 1
}

`
