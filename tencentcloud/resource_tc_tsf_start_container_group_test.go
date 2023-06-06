package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfStartContainerGroupResource_basic -v
func TestAccTencentCloudTsfStartContainerGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfUnitNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfStartContainerGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfStartContainerGroupExists("tencentcloud_tsf_start_container_group.start_container_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_start_container_group.start_container_group", "id"),
				),
			},
			{
				Config: testAccTsfStartContainerGroupUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfStartContainerGroupExists("tencentcloud_tsf_start_container_group.start_container_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_start_container_group.start_container_group", "id"),
				),
			},
		},
	})
}

func testAccCheckTsfStartContainerGroupExists(r string) resource.TestCheckFunc {
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

		if *res.GroupStatus != "Running" && *res.GroupStatus != "Paused" {
			return fmt.Errorf("tsf container group %s start or stop operation failed", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfStartContainerGroup = `

resource "tencentcloud_tsf_start_container_group" "start_container_group" {
  group_id = "group-ynd95rea"
  operate = "stop"
}

`
const testAccTsfStartContainerGroupUp = `

resource "tencentcloud_tsf_start_container_group" "start_container_group" {
  group_id = "group-ynd95rea"
  operate = "start"
}

`
