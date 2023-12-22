package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostUpdateRecordDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostUpdateRecordDetailDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_update_record_detail.describe_host_update_record_detail"),
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
