package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudClsMachineGroup_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsMachineGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsMachineGroupExists("tencentcloud_cls_machine_group.group"),
					resource.TestCheckResourceAttr("tencentcloud_cls_machine_group.group", "group_name", "tf-basic-group"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_machine_group.group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClsMachineGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS machine group][Exists] check: CLS machine group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS machine group][Exists] check: CLS machine group id is not set")
		}
		clsService := ClsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := clsService.DescribeClsMachineGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("[CHECK][CLS machine group][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsMachineGroup = `
resource "tencentcloud_cls_machine_group" "group" {
  group_name        = "tf-basic-group"
  service_logging   = true
  auto_update       = true
  update_end_time   = "19:05:00"
  update_start_time = "17:05:00"

  machine_group_type {
    type   = "ip"
    values = [
      "192.168.1.1",
      "192.168.1.2",
    ]
  }
}
`
