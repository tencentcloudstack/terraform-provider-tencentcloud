package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoAcceleration_domainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoAcceleration_domain,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_acceleration_domain.acceleration_domain", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_acceleration_domain.acceleration_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoAcceleration_domain = `

resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
  zone_id = ""
  domain_name = ""
  origin_info {
		origin_type = ""
		origin = ""
		backup_origin = ""
		private_access = ""
		private_parameters {
			name = ""
			value = ""
		}

  }
}

`
