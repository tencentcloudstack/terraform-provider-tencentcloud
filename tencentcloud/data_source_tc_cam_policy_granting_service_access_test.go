package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
  role_id = 4611686018436805021
  service_type = "cam"
  }

`
