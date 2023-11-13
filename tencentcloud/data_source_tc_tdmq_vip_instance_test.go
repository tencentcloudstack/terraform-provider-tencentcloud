package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqVipInstanceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqVipInstanceDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_vip_instance.vip_instance")),
			},
		},
	})
}

const testAccTdmqVipInstanceDataSource = `

data "tencentcloud_tdmq_vip_instance" "vip_instance" {
  cluster_id = &lt;nil&gt;
    }

`
