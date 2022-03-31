package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCamPoliciesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamPoliciesDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.0.attachments"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.0.create_mode"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.0.policy_id"),
				),
			},
		},
	})
}

const testAccCamPoliciesDataSource_basic = defaultCamVariables + `
data "tencentcloud_cam_policies" "policies" {
  name = var.cam_policy_basic
}
`
