package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapProxiesStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxiesStatusDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies_status.proxies_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies_status.proxies_status", "instance_status_set.#"),
				),
			},
		},
	})
}

const testAccGaapProxiesStatusDataSource = `
data "tencentcloud_gaap_proxies_status" "proxies_status" {
  proxy_ids = ["link-8lpyo88p"]
}
`
