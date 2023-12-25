package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeCompaniesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeCompaniesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_companies.describe_companies"),
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
