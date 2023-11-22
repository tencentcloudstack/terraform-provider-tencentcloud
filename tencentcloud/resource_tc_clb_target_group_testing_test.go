package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTestingClbTargetGroup_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbTargetGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingClbTargetGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group.test", "target_group_name"),
				),
			},
		},
	})
}

func TestAccTencentCloudTestingClbInstanceTargetGroup(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingClbInstanceTargetGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.target_group"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_name", "tgt_grp_test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "port", "33"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "vpc_id", "vpc-humgpppd"),
				),
			},
			{
				Config: testAccTestingClbInstanceTargetGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.target_group"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_name", "tgt_grp_test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "vpc_id", "vpc-humgpppd"),
				),
			},
		},
	})
}

const testAccTestingClbTargetGroup_basic = `
resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "qwe"
}
`

const testAccTestingClbInstanceTargetGroup = `

resource "tencentcloud_clb_target_group" "target_group" {
    target_group_name = "tgt_grp_test"
    port              = 33
    vpc_id            = "vpc-humgpppd"
    target_group_instances {
      bind_ip = "172.16.0.17"
      port = 18800
    }
}
`

const testAccTestingClbInstanceTargetGroupUpdate = `
resource "tencentcloud_clb_target_group" "target_group" {
    target_group_name = "tgt_grp_test"
    port              = 44
    vpc_id            = "vpc-humgpppd"
     target_group_instances {
      bind_ip = "172.16.0.17"
      port = 18800
    }
}
`
