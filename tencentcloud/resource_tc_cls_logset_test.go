package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudClsLogset_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsLogset_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsLogsetExists("tencentcloud_cls_logset.logset"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_logset.logset", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_cls_logset.logset", "logset_name", "tf-logset-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_logset.logset",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClsLogsetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS logset][Exists] check: CLS logset %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS logset][Exists] check: CLS logset id is not set")
		}
		service := ClsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		resourceId := rs.Primary.ID
		instance, err := service.DescribeClsLogsetById(ctx, resourceId)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CLS logset][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsLogset_basic = `
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-logset-test"
  tags        = {
    "test" = "test"
  }
}
`
