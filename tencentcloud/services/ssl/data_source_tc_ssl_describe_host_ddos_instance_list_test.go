package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostDdosInstanceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostDdosInstanceListDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_ddos_instance_list.describe_host_ddos_instance_list"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_ddos_instance_list.describe_host_ddos_instance_list", "certificate_id", "8mCN3eKd"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_ddos_instance_list.describe_host_ddos_instance_list", "resource_type", "ddos"),
				),
			},
		},
	})
}

const testAccSslDescribeHostDdosInstanceListDataSource = `

data "tencentcloud_ssl_describe_host_ddos_instance_list" "describe_host_ddos_instance_list" {
  certificate_id = "8mCN3eKd"
  resource_type = "ddos"
}
`
