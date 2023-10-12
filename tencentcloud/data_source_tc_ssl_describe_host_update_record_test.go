package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostUpdateRecordDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostUpdateRecordDataSource,
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_update_record.describe_host_update_record"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_update_record.describe_host_update_record", "old_certificate_id", "9D3qRt7r"),
				),
			},
		},
	})
}

const testAccSslDescribeHostUpdateRecordDataSource = `

data "tencentcloud_ssl_describe_host_update_record" "describe_host_update_record" {
  old_certificate_id = "9D3qRt7r"
  }

`
