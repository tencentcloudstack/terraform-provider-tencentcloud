package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostDeployRecordDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostDeployRecordDetailDataSource,
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_deploy_record_detail.describe_host_deploy_record_detail"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_deploy_record_detail.describe_host_deploy_record_detail", "deploy_record_id", "35364"),
				),
			},
		},
	})
}

const testAccSslDescribeHostDeployRecordDetailDataSource = `

data "tencentcloud_ssl_describe_host_deploy_record_detail" "describe_host_deploy_record_detail" {
  deploy_record_id = "35364"
}

`
