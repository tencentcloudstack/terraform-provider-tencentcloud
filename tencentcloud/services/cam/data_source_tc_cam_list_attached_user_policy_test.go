package cam_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamListAttachedUserPolicyDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamListAttachedUserPolicyDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cam_list_attached_user_policy.list_attached_user_policy"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_list_attached_user_policy.list_attached_user_policy", "policy_list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_list_attached_user_policy.list_attached_user_policy", "target_uin", "100032767426"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_list_attached_user_policy.list_attached_user_policy", "attach_type", "0"),
				),
			},
		},
	})
}

const testAccCamListAttachedUserPolicyDataSource = `

data "tencentcloud_cam_list_attached_user_policy" "list_attached_user_policy" {
  target_uin = 100032767426
  attach_type = 0
    }

`
