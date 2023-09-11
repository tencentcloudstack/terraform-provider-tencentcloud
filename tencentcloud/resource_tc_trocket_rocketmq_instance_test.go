package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTrocketRocketmqInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTrocketRocketmqInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_trocket_rocketmq_instance.rocketmq_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTrocketRocketmqInstance = `
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name = "test"
  sku_code = "experiment_500"
  remark = "test"
  vpc_id = "vpc-qmvl8z4f"
  subnet_id = "subnet-ncef9v74"
  tags = {
    tag_key = "rocketmq"
    tag_value = "5.x"
  }
}
`
