package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoZoneDDoSPolicyDataSource -v
func TestAccTencentCloudTeoZoneDDoSPolicyDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeoZoneDDoSPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_teo_zone_ddos_policy.zone_ddos_policy"),
				),
			},
		},
	})
}

const testAccDataSourceTeoZoneDDoSPolicyVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}
`

const testAccDataSourceTeoZoneDDoSPolicy = testAccDataSourceTeoZoneDDoSPolicyVar + `

data "tencentcloud_teo_zone_ddos_policy" "zone_ddos_policy" {
  zone_id = var.zone_id
}

`
