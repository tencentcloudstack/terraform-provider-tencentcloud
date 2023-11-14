package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDescribeDeliverLogDownListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDescribeDeliverLogDownListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_live_describe_deliver_log_down_list.describe_deliver_log_down_list")),
			},
		},
	})
}

const testAccLiveDescribeDeliverLogDownListDataSource = `

data "tencentcloud_live_describe_deliver_log_down_list" "describe_deliver_log_down_list" {
    }

`
