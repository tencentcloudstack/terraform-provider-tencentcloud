package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostWafInstanceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostWafInstanceListDataSource,
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_waf_instance_list.describe_host_waf_instance_list"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_waf_instance_list.describe_host_waf_instance_list", "certificate_id", "8hUkH3xC"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_waf_instance_list.describe_host_waf_instance_list", "resource_type", "waf"),
				),
			},
		},
	})
}

const testAccSslDescribeHostWafInstanceListDataSource = `

data "tencentcloud_ssl_describe_host_waf_instance_list" "describe_host_waf_instance_list" {
  certificate_id = "8hUkH3xC"
  resource_type = "waf"
 }

`
