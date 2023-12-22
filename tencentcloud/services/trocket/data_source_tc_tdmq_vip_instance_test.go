package trocket_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentNeedFixCloudTdmqVipInstanceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqVipInstanceDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_vip_instance.vip_instance"),
				),
			},
		},
	})
}

const testAccTdmqVipInstanceDataSource = `
data "tencentcloud_tdmq_vip_instance" "vip_instance" {
  cluster_id = "rocketmq-rd3545bkkj49"
}
`
