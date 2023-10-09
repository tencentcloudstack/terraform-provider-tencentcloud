package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapProxyGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyGroupsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxy_groups.proxy_groups"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxy_groups.proxy_groups", "proxy_group_list.0.group_id", "lg-5anbbou5"),
				),
			},
		},
	})
}

const testAccGaapProxyGroupsDataSource = `
data "tencentcloud_gaap_proxy_groups" "proxy_groups" {
	project_id = 0
	filters {
		  name = "GroupId"
		  values = ["lg-5anbbou5"]
	}
}
`
