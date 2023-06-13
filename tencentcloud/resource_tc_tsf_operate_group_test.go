package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfOperateGroupResource_basic -v
func TestAccTencentCloudTsfOperateGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfUnitNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfOperateGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfOperateGroupExists("tencentcloud_tsf_operate_group.operate_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_operate_group.operate_group", "id"),
				),
			},
			{
				Config: testAccTsfOperateGroupUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfOperateGroupExists("tencentcloud_tsf_operate_group.operate_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_operate_group.operate_group", "id"),
				),
			},
		},
	})
}

func testAccCheckTsfOperateGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfStartGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf group %s is not found", rs.Primary.ID)
		}

		if *res.GroupStatus != "Running" && *res.GroupStatus != "Paused" {
			return fmt.Errorf("tsf group %s start or stop operation failed", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfOperateGroup = `

resource "tencentcloud_tsf_operate_group" "operate_group" {
	group_id = "group-yrjkln9v"
	operate  = "stop"
}

`

const testAccTsfOperateGroupUp = `

resource "tencentcloud_tsf_operate_group" "operate_group" {
	group_id = "group-yrjkln9v"
	operate  = "start"
}

`
