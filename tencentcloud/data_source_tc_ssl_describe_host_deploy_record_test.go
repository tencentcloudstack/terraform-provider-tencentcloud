package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostDeployRecordDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostDeployRecordDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_deploy_record.describe_host_deploy_record")),
			},
		},
	})
}

const testAccSslDescribeHostDeployRecordDataSource = `

data "tencentcloud_ssl_describe_host_deploy_record" "describe_host_deploy_record" {
  certificate_id = "8u8DII0l"
  resource_type = "ddos"
  }

`
