package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcInternetAddressStatisticsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcInternetAddressStatisticsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dc_internet_address_statistics.internet_address_statistics")),
			},
		},
	})
}

const testAccDcInternetAddressStatisticsDataSource = `

data "tencentcloud_dc_internet_address_statistics" "internet_address_statistics" {}

`
