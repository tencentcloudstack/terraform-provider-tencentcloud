package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostDeployRecordDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostDeployRecordDetailDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_deploy_record_detail.describe_host_deploy_record_detail"),
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
