package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcCcn_bandwidthDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcCcn_bandwidthDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_ccn_bandwidth.ccn_bandwidth")),
			},
		},
	})
}

const testAccVpcCcn_bandwidthDataSource = `

data "tencentcloud_vpc_ccn_bandwidth" "ccn_bandwidth" {
  filters {
		name = "source-region"
		values = 

  }
}

`
