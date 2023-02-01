package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudNeedFixTsfLaneResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfLaneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfLane,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfLaneExists("tencentcloud_tsf_lane.lane"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_lane.lane", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_name", "lane-name"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "remark", "lane desc"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_group_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_group_list.0.group_id", "group-yn7j5l8a"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_group_list.0.entrance", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_lane.lane",
				ImportState:       true,
				ImportStateVerify: true,
			},
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

const testAccTsfLane = `

resource "tencentcloud_tsf_lane" "lane" {
	lane_name = "lane-name"
	remark = "lane desc"
	lane_group_list {
		  group_id = "group-yn7j5l8a"
		  entrance = true
		  # lane_group_id = ""
		  # lane_id = ""
		  # group_name = ""
		  # application_id = ""
		  # application_name = ""
		  # namespace_id = ""
		  # namespace_name = ""
		  # create_time =
		  # update_time =
		  # cluster_type = ""
	}
}

`
