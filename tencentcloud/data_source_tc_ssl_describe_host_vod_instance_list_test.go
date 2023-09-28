package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostVodInstanceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostVodInstanceListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_vod_instance_list.describe_host_vod_instance_list")),
			},
		},
	})
}

const testAccSslDescribeHostVodInstanceListDataSource = `

data "tencentcloud_ssl_describe_host_vod_instance_list" "describe_host_vod_instance_list" {
  certificate_id = "8u8DII0l"
  resource_type = "vod"
  }

`
