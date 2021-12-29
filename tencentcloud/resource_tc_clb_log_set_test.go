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
					testAccCheckClbLogsetExists("tencentcloud_clb_log_set.test_logset"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_log_set.test_logset", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_clb_log_set.test_logset", "name", "clb_logset"),
					resource.TestCheckResourceAttr("tencentcloud_clb_log_set.test_logset", "period", "7"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_log_set.test_logset",
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
resource "tencentcloud_clb_log_set" "test_logset" {
  period = 7
}
`
