package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDescribeLiveStreamMonitorListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDescribeLiveStreamMonitorListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_live_describe_live_stream_monitor_list.describe_live_stream_monitor_list")),
			},
		},
	})
}

const testAccLiveDescribeLiveStreamMonitorListDataSource = `

data "tencentcloud_live_describe_live_stream_monitor_list" "describe_live_stream_monitor_list" {
  index = 
    }

`
