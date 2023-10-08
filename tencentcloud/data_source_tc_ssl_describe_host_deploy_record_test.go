package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostDeployRecordDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostDeployRecordDataSource,
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_deploy_record.describe_host_deploy_record"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_deploy_record.describe_host_deploy_record", "certificate_id", "9D3qRt7r"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_deploy_record.describe_host_deploy_record", "resource_type", "ddos"),
				),
			},
		},
	})
}

const testAccSslDescribeHostDeployRecordDataSource = `

data "tencentcloud_ssl_describe_host_deploy_record" "describe_host_deploy_record" {
  certificate_id = "9D3qRt7r"
  resource_type = "ddos"
  }

`
