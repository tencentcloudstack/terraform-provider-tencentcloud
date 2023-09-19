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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "sku_code", "experiment_500"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "remark", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "vpc_end_point"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "tags.tag_key", "rocketmq"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "tags.tag_value", "5.x"),
				),
			},
			{
				Config: testAccTrocketRocketmqInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "name", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "sku_code", "experiment_500"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "remark", "remark update"),
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "vpc_end_point"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "tags.tag_key", "rocketmq"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance", "tags.tag_value", "5.x.x"),
				),
			},
			{
				ResourceName:      "tencentcloud_trocket_rocketmq_instance.rocketmq_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudTrocketRocketmqInstanceResource_enablePublic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTrocketRocketmqInstancePublic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.rocketmq_instance_public", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance_public", "name", "test-public"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance_public", "enable_public", "true"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.rocketmq_instance_public", "bandwidth", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.rocketmq_instance_public", "public_end_point"),
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.rocketmq_instance_public", "vpc_end_point"),
				),
			},
			{
				ResourceName:      "tencentcloud_trocket_rocketmq_instance.rocketmq_instance_public",
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
  remark = "remark"
  vpc_id = "vpc-3a9fo1k9"
  subnet_id = "subnet-8nby1yxg"
  tags = {
    tag_key = "rocketmq"
    tag_value = "5.x"
  }
  enable_public = true
  bandwidth = 1
}
`

const testAccTrocketRocketmqInstanceUpdate = `
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name = "test-update"
  sku_code = "experiment_500"
  remark = "remark update"
  vpc_id = "vpc-3a9fo1k9"
  subnet_id = "subnet-8nby1yxg"
  tags = {
    tag_key = "rocketmq"
    tag_value = "5.x.x"
  }
  enable_public = true
  bandwidth = 1
}
`

const testAccTrocketRocketmqInstancePublic = `
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance_public" {
  instance_type = "EXPERIMENT"
  name = "test-public"
  sku_code = "experiment_500"
  remark = "remark"
  vpc_id = "vpc-3a9fo1k9"
  subnet_id = "subnet-8nby1yxg"
  tags = {
    tag_key = "rocketmq"
    tag_value = "5.x"
  }
  enable_public = true
  bandwidth = 1
}
`
