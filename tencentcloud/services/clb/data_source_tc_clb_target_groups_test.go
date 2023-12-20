package clb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	targetGroupById     = "data.tencentcloud_clb_target_groups.target_group_info_id"
	targetGroupByName   = "data.tencentcloud_clb_target_groups.target_group_info_name"
	targetGroupResource = "tencentcloud_clb_target_group.test"
)

func TestAccTencentCloudClbTargetGroupDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbTargetGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceClbTargetGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists(targetGroupResource),
					resource.TestCheckResourceAttrSet(targetGroupById, "list.#"),
					resource.TestCheckResourceAttrSet(targetGroupById, "list.0.target_group_id"),
					resource.TestCheckResourceAttrSet(targetGroupById, "list.0.vpc_id"),
					resource.TestCheckResourceAttrSet(targetGroupById, "list.0.target_group_name"),
					resource.TestCheckResourceAttrSet(targetGroupById, "list.0.port"),
					resource.TestCheckResourceAttrSet(targetGroupById, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(targetGroupById, "list.0.update_time"),
					resource.TestCheckResourceAttrSet(targetGroupById, "list.0.associated_rule_list.#"),
					resource.TestCheckResourceAttrSet(targetGroupById, "list.0.target_group_instance_list.#"),
				),
			},
			{
				Config: testAccTencentCloudDataSourceClbTargetGroupName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists(targetGroupResource),
					resource.TestCheckResourceAttrSet(targetGroupByName, "list.#"),
					resource.TestCheckResourceAttrSet(targetGroupByName, "list.0.target_group_id"),
					resource.TestCheckResourceAttrSet(targetGroupByName, "list.0.vpc_id"),
					resource.TestCheckResourceAttrSet(targetGroupByName, "list.0.target_group_name"),
					resource.TestCheckResourceAttrSet(targetGroupByName, "list.0.port"),
					resource.TestCheckResourceAttrSet(targetGroupByName, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(targetGroupByName, "list.0.update_time"),
					resource.TestCheckResourceAttrSet(targetGroupByName, "list.0.associated_rule_list.#"),
					resource.TestCheckResourceAttrSet(targetGroupByName, "list.0.target_group_instance_list.#"),
				),
			},
		},
	})
}

const tareGetGroupBase = tcacctest.DefaultVpcSubnets + `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-rule-data"
  vpc_id = local.vpc_id
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.listener_basic.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}

resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test-target-keep-1"
    vpc_id = local.vpc_id
}

resource "tencentcloud_clb_target_group_attachment" "group" {
    clb_id          = tencentcloud_clb_instance.clb_basic.id
    listener_id     = tencentcloud_clb_listener.listener_basic.listener_id
    rule_id         = tencentcloud_clb_listener_rule.rule_basic.rule_id
    target_group_id = tencentcloud_clb_target_group.test.id 
}
`

const testAccTencentCloudDataSourceClbTargetGroup = tareGetGroupBase + `
data "tencentcloud_clb_target_groups" "target_group_info_id" {
  target_group_id = tencentcloud_clb_target_group.test.id
}
`

const testAccTencentCloudDataSourceClbTargetGroupName = tareGetGroupBase + `
data "tencentcloud_clb_target_groups" "target_group_info_name" {
  target_group_name = tencentcloud_clb_target_group.test.target_group_name
}
`
