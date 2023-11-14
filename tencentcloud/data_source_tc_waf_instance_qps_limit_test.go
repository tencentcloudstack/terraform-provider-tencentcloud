package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWafInstanceQpsLimitDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafInstanceQpsLimitDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_instance_qps_limit.instance_qps_limit")),
			},
		},
	})
}

const testAccWafInstanceQpsLimitDataSource = `

data "tencentcloud_waf_instance_qps_limit" "instance_qps_limit" {
  instance_id = ""
  type = ""
  }

`
