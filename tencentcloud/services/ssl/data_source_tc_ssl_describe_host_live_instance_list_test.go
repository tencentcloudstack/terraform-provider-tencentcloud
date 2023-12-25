package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostLiveInstanceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostLiveInstanceListDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_live_instance_list.describe_host_live_instance_list"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_live_instance_list.describe_host_live_instance_list", "certificate_id", "9D3qRt7r"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_live_instance_list.describe_host_live_instance_list", "resource_type", "live"),
				),
			},
		},
	})
}

const testAccSslDescribeHostLiveInstanceListDataSource = `

data "tencentcloud_ssl_describe_host_live_instance_list" "describe_host_live_instance_list" {
  certificate_id = "9D3qRt7r"
  resource_type = "live"
}
`
