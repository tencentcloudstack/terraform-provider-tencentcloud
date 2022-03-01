package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDataSourceUserInfo(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

const testAccDataUserInfoBasic = `
data "tencentcloud_user_info" "info" {}
`
