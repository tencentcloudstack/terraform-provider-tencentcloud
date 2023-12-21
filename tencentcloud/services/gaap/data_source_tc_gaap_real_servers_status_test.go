package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapRealServersStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealServersStatusDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_real_servers_status.real_servers_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_real_servers_status.real_servers_status", "real_server_status_set.#"),
				),
			},
		},
	})
}

const testAccGaapRealServersStatusDataSource = `
data "tencentcloud_gaap_real_servers_status" "real_servers_status" {
	real_server_ids = ["rs-qcygnwpd"]
}
`
