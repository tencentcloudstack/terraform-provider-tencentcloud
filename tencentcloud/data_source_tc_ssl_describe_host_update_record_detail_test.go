package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostUpdateRecordDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostUpdateRecordDetailDataSource,
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_update_record_detail.describe_host_update_record_detail"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_update_record_detail.describe_host_update_record_detail", "deploy_record_id", "1666"),
				),
			},
		},
	})
}

const testAccSslDescribeHostUpdateRecordDetailDataSource = `

data "tencentcloud_ssl_describe_host_update_record_detail" "describe_host_update_record_detail" {
  deploy_record_id = "1666"
}

`
