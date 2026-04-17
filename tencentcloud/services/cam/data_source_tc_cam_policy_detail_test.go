package cam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCamPolicyDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCamPolicyDetailDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cam_policy_detail.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policy_detail.example", "policy_info.0.policy_name"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policy_detail.example", "policy_info.0.type"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_cam_policy_detail.example", "policy_info.0.policy_document"),
			),
		}},
	})
}

const testAccCamPolicyDetailDataSource = `
data "tencentcloud_cam_policy_detail" "example" {
  policy_id = 17698703
}
`
