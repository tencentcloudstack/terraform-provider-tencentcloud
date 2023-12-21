package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapListenerRealServersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapListenerRealServersDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_listener_real_servers.listener_real_servers"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_listener_real_servers.listener_real_servers", "real_server_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_listener_real_servers.listener_real_servers", "bind_real_server_set.#"),
				),
			},
		},
	})
}

const testAccGaapListenerRealServersDataSource = `
data "tencentcloud_gaap_listener_real_servers" "listener_real_servers" {
	listener_id = "listener-4yrzte61"
}
`
