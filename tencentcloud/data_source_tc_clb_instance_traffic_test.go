package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClbInstanceTrafficDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstanceTrafficDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clb_instance_traffic.instance_traffic")),
			},
		},
	})
}

const testAccClbInstanceTrafficDataSource = `

data "tencentcloud_clb_instance_traffic" "instance_traffic" {
  load_balancer_region = "ap-guangzhou"
}
`
