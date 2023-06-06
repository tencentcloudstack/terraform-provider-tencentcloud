package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixVpcNetDetectStatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcNetDetectStatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_net_detect_states.net_detect_states")),
			},
		},
	})
}

const testAccVpcNetDetectStatesDataSource = `

data "tencentcloud_vpc_net_detect_states" "net_detect_states" {
  net_detect_ids = ["netd-12345678"]
}

`
