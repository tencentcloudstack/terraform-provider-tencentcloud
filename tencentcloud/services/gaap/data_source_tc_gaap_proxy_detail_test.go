package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapProxyDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyDetailDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxy_detail.proxy_detail"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxy_detail.proxy_detail", "proxy_detail.#"),
				),
			},
		},
	})
}

const testAccGaapProxyDetailDataSource = `
data "tencentcloud_gaap_proxy_detail" "proxy_detail" {
	proxy_id = "link-8lpyo88p"
}
`
