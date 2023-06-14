package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfOperateContainerGroupResource_basic -v
func TestAccTencentCloudTsfOperateContainerGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfOperateContainerGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfOperateContainerGroupExists("tencentcloud_tsf_operate_container_group.operate_container_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_operate_container_group.operate_container_group", "id"),
				),
			},
			{
				Config: testAccTsfOperateContainerGroupUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfOperateContainerGroupExists("tencentcloud_tsf_operate_container_group.operate_container_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_operate_container_group.operate_container_group", "id"),
				),
			},
		},
	})
}

func testAccCheckTsfOperateContainerGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfStartContainerGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf container group %s is not found", rs.Primary.ID)
		}

		if *res.Status != "Running" && *res.Status != "Paused" {
			return fmt.Errorf("tsf container group %s start or stop operation failed", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfOperateContainerGroup = `

resource "tencentcloud_tsf_operate_container_group" "operate_container_group" {
  group_id = "group-yqml6w3a"
  operate = "stop"
}

`
const testAccTsfOperateContainerGroupUp = `

resource "tencentcloud_tsf_operate_container_group" "operate_container_group" {
  group_id = "group-yqml6w3a"
  operate = "start"
}

`
