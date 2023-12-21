package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapProxyGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyGroupsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxy_groups.proxy_groups"),
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
