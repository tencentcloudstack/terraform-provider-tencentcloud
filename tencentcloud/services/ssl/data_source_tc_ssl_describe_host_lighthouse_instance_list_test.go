package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostLighthouseInstanceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostLighthouseInstanceListDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_lighthouse_instance_list.describe_host_lighthouse_instance_list"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_lighthouse_instance_list.describe_host_lighthouse_instance_list", "certificate_id", "9D3mK31W"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_lighthouse_instance_list.describe_host_lighthouse_instance_list", "resource_type", "lighthouse"),
				),
			},
		},
	})
}

const testAccSslDescribeHostLighthouseInstanceListDataSource = `

data "tencentcloud_ssl_describe_host_lighthouse_instance_list" "describe_host_lighthouse_instance_list" {
  certificate_id = "9D3mK31W"
  resource_type = "lighthouse"
}
`
