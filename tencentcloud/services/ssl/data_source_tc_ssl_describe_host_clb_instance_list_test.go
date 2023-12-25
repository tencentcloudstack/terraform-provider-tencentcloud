package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostClbInstanceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostClbInstanceListDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_clb_instance_list.describe_host_clb_instance_list"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_clb_instance_list.describe_host_clb_instance_list", "certificate_id", "9Bpk7XOu"),
				),
			},
		},
	})
}

const testAccSslDescribeHostClbInstanceListDataSource = `

data "tencentcloud_ssl_describe_host_clb_instance_list" "describe_host_clb_instance_list" {
  certificate_id = "9Bpk7XOu"
}

`
