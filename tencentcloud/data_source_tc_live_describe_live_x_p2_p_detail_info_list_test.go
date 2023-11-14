package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDescribeLiveXP2PDetailInfoListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDescribeLiveXP2PDetailInfoListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_live_describe_live_x_p2_p_detail_info_list.describe_live_x_p2_p_detail_info_list")),
			},
		},
	})
}

const testAccLiveDescribeLiveXP2PDetailInfoListDataSource = `

data "tencentcloud_live_describe_live_x_p2_p_detail_info_list" "describe_live_x_p2_p_detail_info_list" {
  query_time = ""
  type = 
  stream_names = 
  dimension = 
  }

`
