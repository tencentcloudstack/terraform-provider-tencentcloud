package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudLB_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_lb.classic"),
					resource.TestCheckResourceAttr("tencentcloud_lb.classic", "name", "tf-ci-test"),
					resource.TestCheckResourceAttr("tencentcloud_lb.classic", "type", "OPEN"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "forward"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "project_id"),
				),
			},
			{
				Config: testAccLbBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_lb.classic"),
					resource.TestCheckResourceAttr("tencentcloud_lb.classic", "name", "tf-ci-test-update"),
					resource.TestCheckResourceAttr("tencentcloud_lb.classic", "type", "OPEN"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "forward"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "project_id"),
				),
			},
		},
	})
}

func testAccCheckLBDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_lb" {
			continue
		}

		_, err := clbService.DescribeLoadBalancerById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("clb instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccLbBasic = `
resource "tencentcloud_lb" "classic" {
  type    = "OPEN"
  forward = "APPLICATION"
  name    = "tf-ci-test"
}
`

const testAccLbBasicUpdate = `
resource "tencentcloud_lb" "classic" {
  type    = "OPEN"
  forward = "APPLICATION"
  name    = "tf-ci-test-update"
}
`
