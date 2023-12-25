package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeManagersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeManagersDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_managers.describe_managers"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_managers.describe_managers", "company_id", "11772"),
				),
			},
		},
	})
}

const testAccSslDescribeManagersDataSource = `

data "tencentcloud_ssl_describe_managers" "describe_managers" {
  company_id = "11772"
  }

`
