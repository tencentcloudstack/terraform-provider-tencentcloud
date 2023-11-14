package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamPolicyVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamPolicyVersion,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_policy_version.policy_version", "id")),
			},
			{
				ResourceName:      "tencentcloud_cam_policy_version.policy_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamPolicyVersion = `

resource "tencentcloud_cam_policy_version" "policy_version" {
  policy_id = 
  policy_document = ""
  set_as_default = 
}

`
