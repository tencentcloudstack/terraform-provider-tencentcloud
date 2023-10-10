package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapRealServersStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealServersStatusDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_real_servers_status.real_servers_status"),
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
