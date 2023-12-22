package trabbit_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixTdmqRabbitmqVipInstanceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqVipInstanceDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rabbitmq_vip_instance.rabbitmq_vip_instance"),
				),
			},
		},
	})
}

const testAccTdmqRabbitmqVipInstanceDataSource = `

data "tencentcloud_tdmq_rabbitmq_vip_instance" "rabbitmq_vip_instance" {
  filters {
		name = ""
		values = 

  }
  }

`
