package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslDescribeHostLiveInstanceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostLiveInstanceListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_live_instance_list.describe_host_live_instance_list")),
			},
		},
	})
}

const testAccSslDescribeHostLiveInstanceListDataSource = `

data "tencentcloud_ssl_describe_host_live_instance_list" "describe_host_live_instance_list" {
  certificate_id = ""
  resource_type = ""
  is_cache = 
  filters {
		filter_key = ""
		filter_value = ""

  }
  old_certificate_id = ""
  }

`
