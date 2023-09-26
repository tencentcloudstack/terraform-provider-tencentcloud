package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTrocketRocketmqRoleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTrocketRocketmqRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_role.rocketmq_role", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_role.rocketmq_role", "remark", "test for terraform"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_role.rocketmq_role", "perm_write", "false"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_role.rocketmq_role", "perm_read", "true"),
				),
			},
			{
				Config: testAccTrocketRocketmqRoleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_role.rocketmq_role", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_role.rocketmq_role", "remark", "test terraform"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_role.rocketmq_role", "perm_write", "true"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_role.rocketmq_role", "perm_read", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_trocket_rocketmq_role.rocketmq_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTrocketRocketmqRole = `

resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test_role"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-3a9fo1k9"
  subnet_id     = "subnet-8nby1yxg"
  tags          = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_role" "rocketmq_role" {
  instance_id = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  role        = "test_role"
  remark      = "test for terraform"
  perm_write  = false
  perm_read   = true
}

`

const testAccTrocketRocketmqRoleUpdate = `

resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test_role"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-3a9fo1k9"
  subnet_id     = "subnet-8nby1yxg"
  tags          = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_role" "rocketmq_role" {
  instance_id = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  role        = "test_role"
  remark      = "test terraform"
  perm_write  = true
  perm_read   = true
}

`
