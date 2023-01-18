package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfLaneResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfLane,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_lane.lane", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_lane.lane",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfLane = `

resource "tencentcloud_tsf_lane" "lane" {
  lane_name = ""
  remark = ""
  lane_group_list {
		group_id = ""
		entrance = 
		lane_group_id = ""
		lane_id = ""
		group_name = ""
		application_id = ""
		application_name = ""
		namespace_id = ""
		namespace_name = ""
		create_time = 
		update_time = 
		cluster_type = ""

  }
  program_id_list = 
  }

`
