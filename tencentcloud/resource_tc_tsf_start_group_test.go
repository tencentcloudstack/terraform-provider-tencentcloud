package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfStartGroupResource_basic -v
func TestAccTencentCloudTsfStartGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfUnitNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfStartGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfStartGroupExists("tencentcloud_tsf_start_group.start_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_start_group.start_group", "id"),
				),
			},
			{
				Config: testAccTsfStartGroupUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfStartGroupExists("tencentcloud_tsf_start_group.start_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_start_group.start_group", "id"),
				),
			},
		},
	})
}

func testAccCheckTsfStartGroupExists(r string) resource.TestCheckFunc {
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

const testAccTsfStartGroup = `

resource "tencentcloud_tsf_start_group" "start_group" {
	group_id = "group-ynd95rea"
	operate  = "stop"
}

`

const testAccTsfStartGroupUp = `

resource "tencentcloud_tsf_start_group" "start_group" {
	group_id = "group-ynd95rea"
	operate  = "stop"
}

`
