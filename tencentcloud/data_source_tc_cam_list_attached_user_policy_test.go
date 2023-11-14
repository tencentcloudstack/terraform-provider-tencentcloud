package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamListAttachedUserPolicyDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamListAttachedUserPolicyDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cam_list_attached_user_policy.list_attached_user_policy")),
			},
		},
	})
}

const testAccCamListAttachedUserPolicyDataSource = `

data "tencentcloud_cam_list_attached_user_policy" "list_attached_user_policy" {
  target_uin = 
  rp = 
  attach_type = 
  strategy_type = 
  keyword = ""
    }

`
