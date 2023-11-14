package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamGroupUserAccountDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamGroupUserAccountDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cam_group_user_account.group_user_account")),
			},
		},
	})
}

const testAccCamGroupUserAccountDataSource = `

data "tencentcloud_cam_group_user_account" "group_user_account" {
  uid = &lt;nil&gt;
  rp = &lt;nil&gt;
  sub_uin = &lt;nil&gt;
    }

`
