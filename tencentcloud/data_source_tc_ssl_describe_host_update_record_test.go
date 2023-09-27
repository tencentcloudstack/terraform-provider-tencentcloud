package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslDescribeHostUpdateRecordDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostUpdateRecordDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_update_record.describe_host_update_record")),
			},
		},
	})
}

const testAccSslDescribeHostUpdateRecordDataSource = `

data "tencentcloud_ssl_describe_host_update_record" "describe_host_update_record" {
  old_certificate_id = "8u8DII0l"
  }

`
