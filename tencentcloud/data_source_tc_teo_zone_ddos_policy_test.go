package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

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

const testAccDataSourceTeoZoneDDoSPolicy = `

data "tencentcloud_teo_zone_ddos_policy" "zone_ddos_policy" {
  zone_id = ""
}

`
