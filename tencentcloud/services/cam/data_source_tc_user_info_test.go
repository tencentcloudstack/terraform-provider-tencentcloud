package cam_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDataSourceUserInfoBasic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserInfoBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_user_info.info", "app_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_user_info.info", "uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_user_info.info", "owner_uin"),
				),
			},
		},
	})
}

func TestAccTencentCloudDataSourceUserInfoSubAccount(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				// Need use subaccount aksk
				Config: testAccDataUserInfoSubAccount,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_user_info.info_sub_account", "app_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_user_info.info_sub_account", "uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_user_info.info_sub_account", "owner_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_user_info.info_sub_account", "name"),
				),
			},
		},
	})
}

const testAccDataUserInfoBasic = `
data "tencentcloud_user_info" "info" {}
`
const testAccDataUserInfoSubAccount = `
data "tencentcloud_user_info" "info_sub_account" {}
`
