package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapProxyAndStatisticsListenersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyAndStatisticsListenersDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxy_and_statistics_listeners.proxy_and_statistics_listeners"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxy_and_statistics_listeners.proxy_and_statistics_listeners", "proxy_set.#"),
				),
			},
		},
	})
}

const testAccGaapProxyAndStatisticsListenersDataSource = `
data "tencentcloud_gaap_proxy_and_statistics_listeners" "proxy_and_statistics_listeners" {
	project_id = 0
}
`
