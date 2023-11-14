package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamSetPolicyVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamSetPolicyVersion,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_set_policy_version.set_policy_version", "id")),
			},
			{
				ResourceName:      "tencentcloud_cam_set_policy_version.set_policy_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamSetPolicyVersion = `

resource "tencentcloud_cam_set_policy_version" "set_policy_version" {
  policy_id = 
  version_id = 
}

`
