package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostCdnInstanceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostCdnInstanceListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_cdn_instance_list.describe_host_cdn_instance_list")),
			},
		},
	})
}

const testAccSslDescribeHostCdnInstanceListDataSource = `

data "tencentcloud_ssl_describe_host_cdn_instance_list" "describe_host_cdn_instance_list" {
  certificate_id = "8u8DII0l"
  resource_type = "cdn"
}

`
