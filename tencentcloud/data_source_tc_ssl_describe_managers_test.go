package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeManagersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeManagersDataSource,
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_managers.describe_managers"),
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
