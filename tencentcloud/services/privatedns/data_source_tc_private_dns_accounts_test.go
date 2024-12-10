package privatedns_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudPrivateDnsAccountsDataSource_basic -v
func TestAccTencentCloudPrivateDnsAccountsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccPrivateDnsAccountsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_private_dns_accounts.example"),
			),
		}},
	})
}

const testAccPrivateDnsAccountsDataSource = `
data "tencentcloud_private_dns_accounts" "example" {}
`
