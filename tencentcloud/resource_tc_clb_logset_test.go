package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudClbLogset_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbLogsetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbLogset_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbLogsetExists("tencentcloud_clb_logset.test_logset"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_logset.test_logset", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_logset.test_logset", "topic_count"),
					resource.TestCheckResourceAttr("tencentcloud_clb_logset.test_logset", "name", "clb_logset_test1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_logset.test_logset", "period", "7"),
				),
			},
			{
				Config: testAccClbLogset_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbLogsetExists("tencentcloud_clb_logset.test_logset"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_logset.test_logset", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_logset.test_logset", "topic_count"),
					resource.TestCheckResourceAttr("tencentcloud_clb_logset.test_logset", "name", "clb_logset_test2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_logset.test_logset", "period", "7"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_logset.test_logset",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClbLogsetDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clsService := ClsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_logset" {
			continue
		}
		time.Sleep(5 * time.Second)
		resourceId := rs.Primary.ID
		info, err := clsService.DescribeClsLogSetById(ctx, resourceId)
		if info != nil && err == nil {
			return fmt.Errorf("[CHECK][CLB logset][Destroy] check: CLB logset still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbLogsetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB logset][Exists] check: CLB logset %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB logset][Exists] check: CLB logset id is not set")
		}
		service := ClsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		resourceId := rs.Primary.ID
		instance, err := service.DescribeClsLogSetById(ctx, resourceId)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CLB logset][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbLogset_basic = `
resource "tencentcloud_clb_logset" "test_logset" {
  name = "clb_logset_test1"
  period = 7
}
`

const testAccClbLogset_update = `
resource "tencentcloud_clb_logset" "test_logset" {
  name = "clb_logset_test2"
  period = 7
}
`