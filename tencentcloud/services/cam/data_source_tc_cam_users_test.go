package cam_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamUsersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCamUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUsersDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_cam_users.users", "user_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_users.users", "user_list.0.remark"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_users.users", "user_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_users.users", "user_list.0.user_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_users.users", "user_list.0.console_login"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_users.users", "user_list.0.phone_num"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_users.users", "user_list.0.country_code"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_users.users", "user_list.0.email"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_users.users", "user_list.0.uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_users.users", "user_list.0.uid"),
				),
			},
		},
	})
}

const testAccCamUsersDataSource_basic = tcacctest.DefaultCamVariables + `
data "tencentcloud_cam_users" "users" {
  name = var.cam_user_basic
}
`
