package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDescribeTimeShiftRecordDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDescribeTimeShiftRecordDetailDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_live_describe_time_shift_record_detail.describe_time_shift_record_detail")),
			},
		},
	})
}

const testAccLiveDescribeTimeShiftRecordDetailDataSource = `

data "tencentcloud_live_describe_time_shift_record_detail" "describe_time_shift_record_detail" {
  domain = ""
  app_name = ""
  stream_name = ""
  start_time = 
  end_time = 
  domain_group = ""
  trans_code_id = 
  }

`
