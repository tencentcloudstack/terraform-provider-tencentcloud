package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeCompaniesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeCompaniesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_companies.describe_companies"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_companies.describe_companies", "company_id", "122"),
				),
			},
		},
	})
}

const testAccSslDescribeCompaniesDataSource = `

data "tencentcloud_ssl_describe_companies" "describe_companies" {
  company_id = 122
}

`
