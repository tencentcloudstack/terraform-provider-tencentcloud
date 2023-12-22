package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapGroupAndStatisticsProxyDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapGroupAndStatisticsProxyDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_group_and_statistics_proxy.group_and_statistics_proxy"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_group_and_statistics_proxy.group_and_statistics_proxy", "group_set.#"),
				),
			},
		},
	})
}

const testAccGaapGroupAndStatisticsProxyDataSource = `
data "tencentcloud_gaap_group_and_statistics_proxy" "group_and_statistics_proxy" {
  project_id = 0
}
`
