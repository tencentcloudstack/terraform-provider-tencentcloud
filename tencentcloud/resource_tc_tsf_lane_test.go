package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfLaneResource_basic -v
func TestAccTencentCloudTsfLaneResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfLaneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfLane,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfLaneExists("tencentcloud_tsf_lane.lane"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_lane.lane", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_name", "terraform-lane"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "remark", "lane desc"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_group_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_group_list.0.group_id", defaultTsfGroupId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_group_list.0.entrance", "true"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_tsf_lane.lane",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckTsfLaneDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_lane" {
			continue
		}

		res, err := service.DescribeTsfLaneById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf lane %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfLaneExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfLaneById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf lane %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfLaneVar = `
variable "group_id" {
	default = "` + defaultTsfGroupId + `"
}
`

const testAccTsfLane = testAccTsfLaneVar + `

resource "tencentcloud_tsf_lane" "lane" {
	lane_name = "terraform-lane"
	remark = "lane desc"
	lane_group_list {
		  group_id = var.group_id
		  entrance = true
	}
}

`
