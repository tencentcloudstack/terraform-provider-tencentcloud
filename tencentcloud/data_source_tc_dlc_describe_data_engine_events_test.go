package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeDataEngineEventsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeDataEngineEventsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_data_engine_events.describe_data_engine_events")),
			},
		},
	})
}

const testAccDlcDescribeDataEngineEventsDataSource = `

data "tencentcloud_dlc_describe_data_engine_events" "describe_data_engine_events" {
  data_engine_name = "iac-keep-config"
    }

`
