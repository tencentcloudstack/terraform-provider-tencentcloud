package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamPolicyGrantingServiceAccessDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamPolicyGrantingServiceAccessDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cam_policy_granting_service_access.policy_granting_service_access")),
			},
		},
	})
}

const testAccCamPolicyGrantingServiceAccessDataSource = `

data "tencentcloud_cam_policy_granting_service_access" "policy_granting_service_access" {
  target_uin = &lt;nil&gt;
  role_id = &lt;nil&gt;
  group_id = &lt;nil&gt;
  service_type = "cvm"
  }

`
