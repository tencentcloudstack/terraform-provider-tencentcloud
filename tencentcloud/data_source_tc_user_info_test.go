package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDataSourceUserInfoBasic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccDataUserInfoBasic,
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
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Need use subaccount aksk
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_SUB_ACCOUNT) },
				Config:    testAccDataUserInfoSubAccount,
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
