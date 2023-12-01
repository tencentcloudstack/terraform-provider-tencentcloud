package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_policy_version.policy_version", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_policy_version.policy_version", "policy_id", "171173780"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_policy_version.policy_version", "policy_document"),
					resource.TestCheckResourceAttr("tencentcloud_cam_policy_version.policy_version", "set_as_default", "false"),
				),
			},
			{
				ResourceName:            "tencentcloud_cam_policy_version.policy_version",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"policy_document", "set_as_default"},
			},
		},
	})
}

const testAccCamPolicyVersion = `

resource "tencentcloud_cam_policy_version" "policy_version" {
  policy_id = 171173780
  policy_document = jsonencode({
    "version": "2.0",
    "statement": [
      {
        "effect": "allow",
        "action": [
          "sts:AssumeRole"
        ],
        "resource": [
          "*"
        ]
      },
      {
        "effect": "allow",
        "action": [
          "cos:PutObject"
        ],
        "resource": [
          "*"
        ]
      },
      {
        "effect": "deny",
        "action": [
          "aa:*"
        ],
        "resource": [
          "*"
        ]
      },
      {
        "effect": "deny",
        "action": [
          "aa:*"
        ],
        "resource": [
          "*"
        ]
      }
    ]
  })
  set_as_default = "false"
}
`
