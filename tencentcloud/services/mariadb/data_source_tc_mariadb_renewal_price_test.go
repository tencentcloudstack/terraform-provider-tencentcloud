package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbRenewalPriceDataSource_basic -v
func TestAccTencentCloudMariadbRenewalPriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbRenewalPriceDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_renewal_price.renewal_price"),
				),
			},
		},
	})
}

const testAccMariadbRenewalPriceDataSource = `
data "tencentcloud_mariadb_renewal_price" "renewal_price" {
  instance_id = "tdsql-9vqvls95"
  period      = 2
}
`
