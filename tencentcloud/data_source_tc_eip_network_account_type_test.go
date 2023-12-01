package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudEipNetworkAccountTypeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEipNetworkAccountTypeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_eip_network_account_type.network_account_type")),
			},
		},
	})
}

const testAccEipNetworkAccountTypeDataSource = `

data "tencentcloud_eip_network_account_type" "network_account_type" {}

`
