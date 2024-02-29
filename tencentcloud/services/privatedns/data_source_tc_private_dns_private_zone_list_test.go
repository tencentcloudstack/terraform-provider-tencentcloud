package privatedns_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudPrivateDnsPrivateZoneListDataSource_basic -v
func TestAccTencentCloudPrivateDnsPrivateZoneListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatednsPrivateZoneListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_private_dns_private_zone_list.example"),
				),
			},
		},
	})
}

const testAccPrivatednsPrivateZoneListDataSource = `
data "tencentcloud_private_dns_private_zone_list" "example" {}
`
