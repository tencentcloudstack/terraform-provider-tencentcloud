package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPrivateDnsRecordsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsRecordsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_private_dns_records.private_dns_record")),
			},
		},
	})
}

const testAccPrivateDnsRecordsDataSource = `

data "tencentcloud_private_dns_records" "private_dns_record" {
  zone_id = "zone-6t11lof0"
  filters {
	name = "Value"
	values = ["8.8.8.8"]
  }
}
`
