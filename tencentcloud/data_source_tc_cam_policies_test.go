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
					testAccCheckCamPolicyExists("tencentcloud_cam_policy.policy"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_policies.policies", "policy_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_policies.policies", "policy_list.0.name", "cam-policy-test5"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_policies.policies", "policy_list.0.description", "test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.0.attachments"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.0.service_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.0.create_mode"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policies.policies", "policy_list.0.policy_id"),
				),
			},
		},
	})
}

const testAccCamPoliciesDataSource_basic = `
resource "tencentcloud_cam_policy" "policy" {
  name        = "cam-policy-test5"
  document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
  description = "test"
}
 
data "tencentcloud_cam_policies" "policies" {
  policy_id = tencentcloud_cam_policy.policy.id
}
`
