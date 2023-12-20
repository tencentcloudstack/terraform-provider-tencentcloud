package cam_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamGroupUserAccountDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamGroupUserAccountDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cam_group_user_account.group_user_account"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_group_user_account.group_user_account", "sub_uin", "100033690181")),
			},
		},
	})
}

const testAccCamGroupUserAccountDataSource = `

data "tencentcloud_cam_group_user_account" "group_user_account" {
  sub_uin = 100033690181
}
`
