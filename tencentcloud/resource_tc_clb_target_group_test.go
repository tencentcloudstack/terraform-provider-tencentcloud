package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudClbTargetGroup_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbTargetGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group.test", "target_group_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group.test", "vpc_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbInstanceTargetGroup(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstanceTargetGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.target_group"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_name", "tgt_grp_test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "port", "33"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.bind_ip", "10.0.0.4"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.port", "33"),
				),
			},
			{
				Config: testAccClbInstanceTargetGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.target_group"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_name", "tgt_grp_test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "port", "33"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.bind_ip", "10.0.0.4"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.port", "44"),
				),
			},
		},
	})
}

func testAccCheckClbTargetGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_target_group" {
			continue
		}
		time.Sleep(5 * time.Second)
		filters := map[string]string{}
		targetGroupInfos, err := clbService.DescribeTargetGroups(ctx, rs.Primary.ID, filters)
		if len(targetGroupInfos) > 0 && err == nil {
			return fmt.Errorf("[CHECK][CLB target group][Destroy] check: CLB target group still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbTargetGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB target group][Exists] check: CLB target group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB target group][Exists] check: CLB target group id is not set")
		}
		clbService := ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		filters := map[string]string{}
		targetGroupInfos, err := clbService.DescribeTargetGroups(ctx, rs.Primary.ID, filters)
		if err != nil {
			return err
		}
		if len(targetGroupInfos) == 0 {
			return fmt.Errorf("[CHECK][CLB target group][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbTargetGroup_basic = `
resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "qwe"
}
`

const testAccClbInstanceTargetGroup = `
resource "tencentcloud_clb_target_group" "target_group" {
    target_group_name = "tgt_grp_test"
    port              = 33
    vpc_id            = "vpc-4owdpnwr"
    target_group_instances {
      bind_ip = "172.16.16.95"
      port = 18800
    }
}
`

const testAccClbInstanceTargetGroupUpdate = `
resource "tencentcloud_clb_target_group" "target_group" {
    target_group_name = "tgt_grp_test"
    port              = 44
	vpc_id            = "vpc-4owdpnwr"
    target_group_instances {
      bind_ip = "172.16.16.95"
      port = 18800
    }
}
`
