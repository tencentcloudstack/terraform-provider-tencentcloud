package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClbListenersByTargetsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListenersByTargetsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clb_listeners_by_targets.listeners_by_targets")),
			},
		},
	})
}

const testAccClbListenersByTargetsDataSource = `

data "tencentcloud_clb_listeners_by_targets" "listeners_by_targets" {
  backends {
		vpc_id = ""
		private_ip = ""

  }
  }

`
