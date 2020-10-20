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
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB target group][Destroy] check: CLB target group still exists: %s", rs.Primary.ID)
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
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB target group][Exists] check: CLB target group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB target group][Exists] check: CLB target group id is not set")
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
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB target group][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbTargetGroup_basic = `
resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "qwe"
}
`
